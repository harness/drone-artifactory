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
	// Cleanup parameters
	CleanupPattern string `envconfig:"PLUGIN_CLEANUP_PATTERN"`

	// PLUGIN_COMMAND
	Command string `envconfig:"PLUGIN_COMMAND"`

	// Xray parameters
	XrayWatchName   string `envconfig:"PLUGIN_XRAY_WATCH_NAME"`
	XrayBuildName   string `envconfig:"PLUGIN_XRAY_BUILD_NAME"`
	XrayBuildNumber string `envconfig:"PLUGIN_XRAY_BUILD_NUMBER"`

	// Docker parameters
	DockerImageName string `envconfig:"PLUGIN_DOCKER_IMAGE_NAME"`
	DockerRepo      string `envconfig:"PLUGIN_DOCKER_REPO"`
	DockerUsername  string `envconfig:"PLUGIN_DOCKER_USERNAME"`
	DockerPassword  string `envconfig:"PLUGIN_DOCKER_PASSWORD"`

	// Maven parameters
	MvnResolveReleases  string `envconfig:"PLUGIN_REPO_RESOLVE_RELEASES"`
	MvnResolveSnapshots string `envconfig:"PLUGIN_REPO_RESOLVE_SNAPSHOTS"`
	MvnDeployReleases   string `envconfig:"PLUGIN_REPO_DEPLOY_RELEASES"`
	MvnDeploySnapshots  string `envconfig:"PLUGIN_REPO_DEPLOY_SNAPSHOTS"`
	MvnGoals            string `envconfig:"PLUGIN_GOALS"`
	MvnPomFile          string `envconfig:"PLUGIN_POM_FILE"`
	ProjectKey          string `envconfig:"PLUGIN_PROJECT_KEY"`
	OptionalArgs        string `envconfig:"PLUGIN_OPTIONAL_ARGS"`

	// Gradle parameters
	DeployIvyDesc   string `envconfig:"PLUGIN_DEPLOY_IVY_DESC"`
	DeployMavenDesc string `envconfig:"PLUGIN_DEPLOY_MAVEN_DESC"`
	Global          string `envconfig:"PLUGIN_GLOBAL"`
	IvyArtifacts    string `envconfig:"PLUGIN_IVY_ARTIFACTS_PATTERN"`
	IvyDesc         string `envconfig:"PLUGIN_IVY_DESC_PATTERN"`
	RepoDeploy      string `envconfig:"PLUGIN_REPO_DEPLOY"`
	RepoResolve     string `envconfig:"PLUGIN_REPO_RESOLVE"`
	ServerIdDeploy  string `envconfig:"PLUGIN_SERVER_ID_DEPLOY"`
	ServerIdResolve string `envconfig:"PLUGIN_SERVER_ID_RESOLVE"`
	UseWrapper      string `envconfig:"PLUGIN_USE_WRAPPER"`
	GradleTasks     string `envconfig:"PLUGIN_TASKS"`

	// UploadParams
	Ant               string `envconfig:"PLUGIN_ANT"`
	Archive           string `envconfig:"PLUGIN_ARCHIVE"`
	ChunkSize         string `envconfig:"PLUGIN_CHUNK_SIZE"`
	ClientCertKeyPath string `envconfig:"PLUGIN_CLIENT_CERT_KEY_PATH"`
	ClientCertPath    string `envconfig:"PLUGIN_CLIENT_CERT_PATH"`
	Deb               string `envconfig:"PLUGIN_DEB"`
	DetailedSummary   string `envconfig:"PLUGIN_DETAILED_SUMMARY"`
	DryRun            string `envconfig:"PLUGIN_DRY_RUN"`
	Exclusions        string `envconfig:"PLUGIN_EXCLUSIONS"`
	Explode           string `envconfig:"PLUGIN_EXPLODE"`
	FailNoOp          string `envconfig:"PLUGIN_FAIL_NO_OP"`
	IncludeDirs       string `envconfig:"PLUGIN_INCLUDE_DIRS"`
	MinSplit          string `envconfig:"PLUGIN_MIN_SPLIT"`
	Module            string `envconfig:"PLUGIN_MODULE"`
	Project           string `envconfig:"PLUGIN_PROJECT"`
	Quiet             string `envconfig:"PLUGIN_QUIET"`
	Recursive         string `envconfig:"PLUGIN_RECURSIVE"`
	Regexp            string `envconfig:"PLUGIN_REGEXP"`
	RetryWaitTime     string `envconfig:"PLUGIN_RETRY_WAIT_TIME"`
	ServerId          string `envconfig:"PLUGIN_SERVER_ID"`
	SplitCount        string `envconfig:"PLUGIN_SPLIT_COUNT"`
	SSHKeyPath        string `envconfig:"PLUGIN_SSH_KEY_PATH"`
	SSHPassphrase     string `envconfig:"PLUGIN_SSH_PASSPHRASE"`
	Symlinks          string `envconfig:"PLUGIN_SYMLINKS"`
	SyncDeletes       string `envconfig:"PLUGIN_SYNC_DELETES"`

	// DownloadParams
	ArchiveEntries          string `envconfig:"PLUGIN_ARCHIVE_ENTRIES"`
	Build                   string `envconfig:"PLUGIN_BUILD"`
	Bundle                  string `envconfig:"PLUGIN_BUNDLE"`
	BypassArchiveInspection string `envconfig:"PLUGIN_BYPASS_ARCHIVE_INSPECTION"`
	ExcludeArtifacts        string `envconfig:"PLUGIN_EXCLUDE_ARTIFACTS"`
	ExcludeProps            string `envconfig:"PLUGIN_EXCLUDE_PROPS"`
	GpgKey                  string `envconfig:"PLUGIN_GPG_KEY"`
	IncludeDeps             string `envconfig:"PLUGIN_INCLUDE_DEPS"`
	Limit                   string `envconfig:"PLUGIN_LIMIT"`
	Offset                  string `envconfig:"PLUGIN_OFFSET"`
	Props                   string `envconfig:"PLUGIN_PROPS"`
	SkipChecksum            string `envconfig:"PLUGIN_SKIP_CHECKSUM"`
	SortBy                  string `envconfig:"PLUGIN_SORT_BY"`
	SortOrder               string `envconfig:"PLUGIN_SORT_ORDER"`
	ValidateSymlinks        string `envconfig:"PLUGIN_VALIDATE_SYMLINKS"`
	BypassArchiveLimits     string `envconfig:"PLUGIN_BYPASS_ARCHIVE_LIMITS"`

	// xray
	ScanFilePattern string `envconfig:"PLUGIN_SCAN_FILE_PATTERN"`
	Vuln            string `envconfig:"PLUGIN_VULN"`
	ExtendedTable   string `envconfig:"PLUGIN_EXTENDED_TABLE"`
	Fail            string `envconfig:"PLUGIN_FAIL"`
	FixableOnly     string `envconfig:"PLUGIN_FIXABLE_ONLY"`
	Format          string `envconfig:"PLUGIN_FORMAT"`
	License         string `envconfig:"PLUGIN_LICENSE"`
	MinSeverity     string `envconfig:"PLUGIN_MIN_SEVERITY"`
	RepoPath        string `envconfig:"PLUGIN_REPO_PATH"`
	Watches         string `envconfig:"PLUGIN_WATCHES"`
}

func Exec(ctx context.Context, args Args) error {
	var cmdArgs []string
	var err error

	switch {
	case len(args.Command) > 0:
		_, err := handleRtCommand(ctx, args)

		return err
	default:
		cmdArgs, err = NativeJfCommandExec(ctx, args)
	}

	enableProxy := parseBoolOrDefault(false, args.EnableProxy)
	if enableProxy {
		log.Printf("setting proxy config for upload")
		setSecureConnectProxies()
	}

	cmdStr := strings.Join(cmdArgs[:], " ")

	shell, shArg := getShell()

	cmd := exec.Command(shell, shArg, cmdStr)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "JFROG_CLI_OFFER_CONFIG=false")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)

	err = cmd.Run()
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

func ExecCommand(args Args, cmdArgs []string) error {

	cmdStr := strings.Join(cmdArgs[:], " ")

	shell, shArg := getShell()

	fmt.Println()
	fmt.Printf("%s %s %s", shell, shArg, cmdStr)
	fmt.Println()

	cmd := exec.Command(shell, shArg, cmdStr)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "JFROG_CLI_OFFER_CONFIG=false")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)

	err := cmd.Run()
	if err != nil {
		fmt.Println(" Error: ", err)
		return err
	}

	if args.PublishBuildInfo {
		if err := publishBuildInfo(args); err != nil {
			return err
		}
	}

	return nil
}

func handleRtCommand(ctx context.Context, args Args) ([][]string, error) {
	commandsList := [][]string{}
	var err error

	switch args.Command {
	case MvnCmd:
		commandsList, err = GetMavenCommandArgs(args)
	case GradleCmd:
		commandsList, err = GetGradleCommandArgs(args)
	case UploadCmd:
		commandsList, err = GetUploadCommandArgs(args)
	case DownloadCmd:
		commandsList, err = GetDownloadCommandArgs(args)
	case CleanUpCmd:
		commandsList, err = GetCleanupCommandArgs(args)
	case BuildInfoCmd:
		commandsList, err = GetBuildInfoCommandArgs(args)
	case PromoteCmd:
		commandsList, err = GetPromoteCommandArgs(args)
	case ScanCommand:
		commandsList, err = GetScanCommandArgs(args)
	case PublishCommand:
		commandsList, err = GetPublishCommandArgs(args)
	}

	for _, cmd := range commandsList {
		execArgs := []string{getJfrogBin()}
		execArgs = append(execArgs, cmd...)
		err := ExecCommand(args, execArgs)
		if err != nil {
			return commandsList, err
		}
		fmt.Println()
	}
	fmt.Println()

	return commandsList, err
}

func NativeJfCommandExec(ctx context.Context, args Args) ([]string, error) {

	if args.URL == "" {
		return []string{}, fmt.Errorf("url needs to be set")
	}

	cmdArgs := []string{getJfrogBin(), "rt", "u", fmt.Sprintf("--url %s", args.URL)}
	if args.Retries != 0 {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--retries=%d", args.Retries))
	}

	// Set authentication params
	cmdArgs, error := setAuthParams(cmdArgs, args)
	if error != nil {
		return []string{}, error
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
				return []string{}, fmt.Errorf("error creating pem folder: %s", pemFolderErr)
			}
			// write pem contents
			pemWriteErr := os.WriteFile(path, []byte(args.PEMFileContents), 0600)
			if pemWriteErr != nil {
				return []string{}, fmt.Errorf("error writing pem file: %s", pemWriteErr)
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
			return []string{}, fmt.Errorf("source file needs to be set")
		}
		if args.Target == "" {
			return []string{}, fmt.Errorf("target path needs to be set")
		}
		cmdArgs = append(cmdArgs, fmt.Sprintf("\"%s\"", args.Source), args.Target)
	}
	return cmdArgs, nil
}

func setupCommonArgs(cmdArgs []string, args Args) ([]string, error) {
	// Add URL
	cmdArgs = append(cmdArgs, fmt.Sprintf("--url=%s", args.URL))

	// Add retries if set
	if args.Retries != 0 {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--retries=%d", args.Retries))
	}

	// Set authentication parameters
	var err error
	cmdArgs, err = setAuthParams(cmdArgs, args)
	if err != nil {
		return nil, err
	}

	// Handle insecure flag
	insecure := parseBoolOrDefault(false, args.Insecure)
	if insecure {
		cmdArgs = append(cmdArgs, "--insecure-tls")
	}

	// Create PEM file if necessary
	if args.PEMFileContents != "" && !insecure {
		err := createPemFile(args.PEMFileContents, args.PEMFilePath)
		if err != nil {
			return nil, err
		}
	}

	return cmdArgs, nil
}

/*
func handlePromote(cmdArgs []string, args Args) ([]string, error) {
	// Set up common arguments
	var err error
	cmdArgs, err = setupCommonArgs(cmdArgs, args)
	if err != nil {
		return cmdArgs, err
	}
	if args.SourceRepo == "" || args.TargetRepo == "" {
		log.Fatalf("source repo and target repo need to be set for promote")
	}
	cmdArgs = append(cmdArgs, fmt.Sprintf("--source-repo='%s'", args.SourceRepo))
	cmdArgs = append(cmdArgs, fmt.Sprintf("--target-repo='%s'", args.TargetRepo))
	if args.PromotionStatus != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--status='%s'", args.PromotionStatus))
	}
	if args.PromotionComment != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--comment='%s'", args.PromotionComment))
	}
	if args.IncludeDeps {
		cmdArgs = append(cmdArgs, "--include-dependencies")
	}
	return cmdArgs, nil
}
*/

func handleBuildInfo(cmdArgs []string, args Args) ([]string, error) {
	// Set up common arguments
	var err error
	cmdArgs, err = setupCommonArgs(cmdArgs, args)
	if err != nil {
		return cmdArgs, err
	}
	if args.BuildName == "" || args.BuildNumber == "" {
		log.Fatalf("build name and build number need to be set for build-info")
	}
	cmdArgs = append(cmdArgs, fmt.Sprintf("--build-name='%s'", args.BuildName))
	cmdArgs = append(cmdArgs, fmt.Sprintf("--build-number='%s'", args.BuildNumber))
	return cmdArgs, nil
}

func handleCleanup(cmdArgs []string, args Args) ([]string, error) {
	// Set up common arguments
	cmdArgs, err := setupCommonArgs(cmdArgs, args)
	if err != nil {
		return cmdArgs, err
	}

	// Handle cleanup-specific arguments
	if args.CleanupPattern != "" {
		// Create cleanup-spec.json with pattern and target
		spec := fmt.Sprintf(`{
			"files": [
				{
					"pattern": "%s",
					"delete": true
				}
			]
		}`, args.CleanupPattern)

		// Write spec to file using os.WriteFile (new in Go 1.16)
		specFilePath := "cleanup-spec.json"
		err := os.WriteFile(specFilePath, []byte(spec), 0644) // Replacing ioutil.WriteFile
		if err != nil {
			return cmdArgs, fmt.Errorf("failed to write spec file: %v", err)
		}

		// Add the spec file to the command arguments
		cmdArgs = append(cmdArgs, fmt.Sprintf("--spec=%s", specFilePath))
		cmdArgs = append(cmdArgs, "--quiet")
	}

	return cmdArgs, nil
}

func handleDocker(cmdArgs []string, args Args) ([]string, error) {
	// Ensure both Docker image name and Docker repo are provided
	if args.DockerImageName == "" || args.DockerRepo == "" {
		log.Fatalf("docker image name and docker repo need to be set for docker push")
	}

	// Set up common arguments
	cmdArgs, err := setupCommonArgs(cmdArgs, args)
	if err != nil {
		return cmdArgs, err
	}

	// Append Docker image tag and target repository
	cmdArgs = append(cmdArgs, fmt.Sprintf("%s %s", args.DockerImageName, args.DockerRepo))

	// Optionally, handle username and password for authentication if provided
	if args.DockerUsername != "" && args.DockerPassword != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--user=%s:%s", args.DockerUsername, args.DockerPassword))
	}

	return cmdArgs, nil
}

func handleXrayScan(cmdArgs []string, args Args) ([]string, error) {
	// Set up common arguments
	var err error
	cmdArgs, err = setupCommonArgs(cmdArgs, args)
	if err != nil {
		return cmdArgs, err
	}
	if args.XrayWatchName != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--watch='%s'", args.XrayWatchName))
	}
	if args.XrayBuildName != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--build-name='%s'", args.XrayBuildName))
	}
	if args.XrayBuildNumber != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--build-number='%s'", args.XrayBuildNumber))
	}
	return cmdArgs, nil
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

// createPemFile writes the PEM file to the specified path
func createPemFile(pemContents, pemFilePath string) error {
	var path string
	// Determine path to write pem file
	if pemFilePath == "" {
		if runtime.GOOS == "windows" {
			path = "C:/users/ContainerAdministrator/.jfrog/security/certs/cert.pem"
		} else {
			path = "/root/.jfrog/security/certs/cert.pem"
		}
	} else {
		path = pemFilePath
	}

	fmt.Printf("Creating pem file at %q\n", path)

	// Create folder and write PEM contents
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir := filepath.Dir(path)
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return fmt.Errorf("failed to create pem folder: %v", err)
		}
	}

	err := os.WriteFile(path, []byte(pemContents), 0600)
	if err != nil {
		return fmt.Errorf("failed to create pem file %v", err)
	}

	return nil
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
