// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	harnessHTTPProxy  = "HARNESS_HTTP_PROXY"
	harnessHTTPSProxy = "HARNESS_HTTPS_PROXY"
	harnessNoProxy    = "HARNESS_NO_PROXY"
	httpProxy         = "HTTP_PROXY"
	httpsProxy        = "HTTPS_PROXY"
	noProxy           = "NO_PROXY"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	// TODO replace or remove
	Username         string `envconfig:"PLUGIN_USERNAME"`
	Password         string `envconfig:"PLUGIN_PASSWORD"`
	APIKey           string `envconfig:"PLUGIN_API_KEY"`
	AccessToken      string `envconfig:"PLUGIN_ACCESS_TOKEN"`
	URL              string `envconfig:"PLUGIN_URL"`
	Source           string `envconfig:"PLUGIN_SOURCE"`
	Target           string `envconfig:"PLUGIN_TARGET"`
	Retries          int    `envconfig:"PLUGIN_RETRIES"`
	Flat             string `envconfig:"PLUGIN_FLAT"`
	Spec             string `envconfig:"PLUGIN_SPEC"`
	Threads          int    `envconfig:"PLUGIN_THREADS"`
	SpecVars         string `envconfig:"PLUGIN_SPEC_VARS"`
	TargetProps      string `envconfig:"PLUGIN_TARGET_PROPS"`
	Insecure         string `envconfig:"PLUGIN_INSECURE"`
	PEMFileContents  string `envconfig:"PLUGIN_PEM_FILE_CONTENTS"`
	PEMFilePath      string `envconfig:"PLUGIN_PEM_FILE_PATH"`
	BuildNumber      string `envconfig:"PLUGIN_BUILD_NUMBER"`
	BuildName        string `envconfig:"PLUGIN_BUILD_NAME"`
	PublishBuildInfo bool   `envconfig:"PLUGIN_PUBLISH_BUILD_INFO"`
	EnableProxy      string `envconfig:"PLUGIN_ENABLE_PROXY"`

	// RT commands
	BuildTool string `envconfig:"PLUGIN_BUILD_TOOL"`
	Command   string `envconfig:"PLUGIN_COMMAND"`

	// Mvn commands
	ResolveReleaseRepo  string `envconfig:"PLUGIN_RESOLVE_RELEASE_REPO"`
	ResolveSnapshotRepo string `envconfig:"PLUGIN_RESOLVE_SNAPSHOT_REPO"`
	DeployReleaseRepo   string `envconfig:"PLUGIN_DEPLOY_RELEASE_REPO"`
	DeploySnapshotRepo  string `envconfig:"PLUGIN_DEPLOY_SNAPSHOT_REPO"`
	DeployRepo          string `envconfig:"PLUGIN_DEPLOY_REPO"`
	MvnGoals            string `envconfig:"PLUGIN_GOALS"`
	MvnPomFile          string `envconfig:"PLUGIN_POM_FILE"`
	DeployerId          string `envconfig:"PLUGIN_DEPLOYER_ID"`
	ResolverId          string `envconfig:"PLUGIN_RESOLVER_ID"`

	// Gradle commands
	GradleTasks string `envconfig:"PLUGIN_TASKS"`
	BuildFile   string `envconfig:"PLUGIN_BUILD_FILE"`
	RepoDeploy  string `envconfig:"PLUGIN_REPO_DEPLOY"`
	RepoResolve string `envconfig:"PLUGIN_REPO_RESOLVE"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {

	if args.BuildTool != "" {
		log.Println("Handle RT commands")
		return HandleRtCommands(args)
	}

	enableProxy := parseBoolOrDefault(false, args.EnableProxy)
	if enableProxy {
		log.Printf("setting proxy config for upload")
		setSecureConnectProxies()
	}

	// write code here
	if args.URL == "" {
		return fmt.Errorf("url needs to be set")
	}

	cmdArgs := []string{getJfrogBin(), "rt", "u", fmt.Sprintf("--url %s", args.URL)}
	if args.Retries != 0 {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--retries=%d", args.Retries))
	}

	// Set authentication params
	cmdArgs, error := setAuthParams(cmdArgs, args)
	if error != nil {
		return error
	}

	flat := parseBoolOrDefault(false, args.Flat)
	cmdArgs = append(cmdArgs, fmt.Sprintf("--flat=%s", strconv.FormatBool(flat)))

	if args.Threads > 0 {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--threads=%d", args.Threads))
	}
	// Set insecure flag
	insecure := parseBoolOrDefault(false, args.Insecure)
	if insecure {
		cmdArgs = append(cmdArgs, "--insecure-tls")
	}

	// Add --build-number and --build-name flags if provided
	if args.BuildNumber != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--build-number=%s", args.BuildNumber))
	}
	if args.BuildName != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--build-name='%s'", args.BuildName))
	}

	// create pem file
	if args.PEMFileContents != "" && !insecure {
		var path string
		// figure out path to write pem file
		if args.PEMFilePath == "" {
			if runtime.GOOS == "windows" {
				path = "C:/users/ContainerAdministrator/.jfrog/security/certs/cert.pem"
			} else {
				path = "/root/.jfrog/security/certs/cert.pem"
			}
		} else {
			path = args.PEMFilePath
		}
		fmt.Printf("Creating pem file at %q\n", path)
		// write pen contents to path
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// remove filename from path
			dir := filepath.Dir(path)
			pemFolderErr := os.MkdirAll(dir, 0700)
			if pemFolderErr != nil {
				return fmt.Errorf("error creating pem folder: %s", pemFolderErr)
			}
			// write pem contents
			pemWriteErr := os.WriteFile(path, []byte(args.PEMFileContents), 0600)
			if pemWriteErr != nil {
				return fmt.Errorf("error writing pem file: %s", pemWriteErr)
			}
			fmt.Printf("Successfully created pem file at %q\n", path)
		}
	}
	// Take in spec file or use source/target arguments
	if args.Spec != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--spec=%s", args.Spec))
		if args.SpecVars != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("--spec-vars='%s'", args.SpecVars))
		}
	} else {
		filteredTargetProps := filterTargetProps(args.TargetProps)
		if filteredTargetProps != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("--target-props='%s'", filteredTargetProps))
		}
		if args.Source == "" {
			return fmt.Errorf("source file needs to be set")
		}
		if args.Target == "" {
			return fmt.Errorf("target path needs to be set")
		}
		cmdArgs = append(cmdArgs, fmt.Sprintf("\"%s\"", args.Source), args.Target)
	}

	cmdStr := strings.Join(cmdArgs[:], " ")

	shell, shArg := getShell()

	cmd := exec.Command(shell, shArg, cmdStr)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "JFROG_CLI_OFFER_CONFIG=false")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)

	err := cmd.Run()
	if err != nil {
		return err
	}

	// Call publishBuildInfo if PLUGIN_PUBLISH_BUILD_INFO is set to true
	if args.PublishBuildInfo {
		if err := publishBuildInfo(args); err != nil {
			return err
		}
	}

	return nil
}

func publishBuildInfo(args Args) error {
	if args.BuildName == "" || args.BuildNumber == "" {
		return fmt.Errorf("both build name and build number need to be set when publishing build info")
	}

	sanitizedURL, err := sanitizeURL(args.URL)
	if err != nil {
		return err
	}

	publishCmdArgs := []string{
		getJfrogBin(),
		"rt",
		"build-publish",
		"\"" + args.BuildName + "\"",
		"\"" + args.BuildNumber + "\"",
		fmt.Sprintf("--url=%s", sanitizedURL),
	}

	if args.AccessToken != "" {
		publishCmdArgs = append(publishCmdArgs, fmt.Sprintf("--access-token=%sPLUGIN_ACCESS_TOKEN", getEnvPrefix()))
	} else if args.Username != "" && args.Password != "" {
		publishCmdArgs = append(publishCmdArgs, fmt.Sprintf("--user=%sPLUGIN_USERNAME", getEnvPrefix()))
		publishCmdArgs = append(publishCmdArgs, fmt.Sprintf("--password=%sPLUGIN_PASSWORD", getEnvPrefix()))
	} else {
		return fmt.Errorf("either access token or username/password need to be set for publishing build info")
	}

	publishCmdStr := strings.Join(publishCmdArgs, " ")
	shell, shArg := getShell()
	publishCmd := exec.Command(shell, shArg, publishCmdStr)
	publishCmd.Env = os.Environ()
	publishCmd.Env = append(publishCmd.Env, "JFROG_CLI_OFFER_CONFIG=false")
	publishCmd.Stdout = os.Stdout
	publishCmd.Stderr = os.Stderr
	trace(publishCmd)

	if err := publishCmd.Run(); err != nil {
		return fmt.Errorf("error publishing build info: %s", err)
	}

	return nil
}

// Function to filter TargetProps based on criteria
func filterTargetProps(rawProps string) string {
	keyValuePairs := strings.Split(rawProps, ",")
	validPairs := []string{}

	for _, pair := range keyValuePairs {
		keyValuePair := strings.SplitN(pair, "=", 2)
		if len(keyValuePair) != 2 {
			continue // skip if it's not a valid key-value pair
		}

		key := strings.TrimSpace(keyValuePair[0])
		value := strings.TrimSpace(keyValuePair[1])

		// Remove single or double quotes from value
		trimmedValue := strings.Trim(value, "\"'")

		// Check value is not empty, not "null", and not just whitespace
		if trimmedValue != "" && strings.ToLower(trimmedValue) != "null" {
			validPairs = append(validPairs, key+"="+value)
		}
	}

	return strings.Join(validPairs, ",")
}

// sanitizeURL trims the URL to include only up to the '/artifactory/' path.
func sanitizeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %s", inputURL)
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL: %s", inputURL)
	}
	parts := strings.Split(parsedURL.Path, "/artifactory")
	if len(parts) < 2 {
		return "", fmt.Errorf("url does not contain '/artifactory': %s", inputURL)
	}

	// Always set the path to the first part + "/artifactory/"
	parsedURL.Path = parts[0] + "/artifactory/"

	return parsedURL.String(), nil
}

// setAuthParams appends authentication parameters to cmdArgs based on the provided credentials.
func setAuthParams(cmdArgs []string, args Args) ([]string, error) {
	// Set authentication params
	envPrefix := getEnvPrefix()
	if args.Username != "" && args.Password != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--user %sPLUGIN_USERNAME", envPrefix))
		cmdArgs = append(cmdArgs, fmt.Sprintf("--password %sPLUGIN_PASSWORD", envPrefix))
	} else if args.APIKey != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--apikey %sPLUGIN_API_KEY", envPrefix))
	} else if args.AccessToken != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--access-token %sPLUGIN_ACCESS_TOKEN", envPrefix))
	} else {
		return nil, fmt.Errorf("either username/password, api key or access token needs to be set")
	}
	return cmdArgs, nil
}

func getShell() (string, string) {
	if runtime.GOOS == "windows" {
		return "powershell", "-Command"
	}

	return "sh", "-c"
}

func getJfrogBin() string {
	if runtime.GOOS == "windows" {
		return "C:/bin/jfrog.exe"
	}
	return "jfrog"
}

func getEnvPrefix() string {
	if runtime.GOOS == "windows" {
		return "$Env:"
	}
	return "$"
}

func parseBoolOrDefault(defaultValue bool, s string) (result bool) {
	var err error
	result, err = strconv.ParseBool(s)
	if err != nil {
		result = defaultValue
	}

	return
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}

func setSecureConnectProxies() {
	copyEnvVariableIfExists(harnessHTTPProxy, httpProxy)
	copyEnvVariableIfExists(harnessHTTPSProxy, httpsProxy)
	copyEnvVariableIfExists(harnessNoProxy, noProxy)
}

func copyEnvVariableIfExists(src string, dest string) {
	srcValue := os.Getenv(src)
	if srcValue == "" {
		return
	}
	err := os.Setenv(dest, srcValue)
	if err != nil {
		log.Printf("Failed to copy env variable from %s to %s with error %v", src, dest, err)
	}
}
