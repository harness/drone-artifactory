// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
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

	// PLUGIN_COMMAND
	Command string `envconfig:"PLUGIN_COMMAND"`

	// Maven parameters
	MvnResolveReleases  string `envconfig:"PLUGIN_REPO_RESOLVE_RELEASES"`
	MvnResolveSnapshots string `envconfig:"PLUGIN_REPO_RESOLVE_SNAPSHOTS"`
	MvnDeployReleases   string `envconfig:"PLUGIN_REPO_DEPLOY_RELEASES"`
	MvnDeploySnapshots  string `envconfig:"PLUGIN_REPO_DEPLOY_SNAPSHOTS"`
	MvnGoals            string `envconfig:"PLUGIN_GOALS"`
	MvnPomFile          string `envconfig:"PLUGIN_POM_FILE"`
	ProjectKey          string `envconfig:"PLUGIN_PROJECT_KEY"`
	OptionalArgs        string `envconfig:"PLUGIN_OPTIONAL_ARGS"`

	DeployerId      string `envconfig:"PLUGIN_DEPLOYER_ID"`
	ResolverId      string `envconfig:"PLUGIN_RESOLVER_ID"`
	DeployArtifacts string `envconfig:"PLUGIN_DEPLOY_ARTIFACTS"`
	ExcludePatterns string `envconfig:"PLUGIN_EXCLUDE_PATTERNS"`
	IncludePatterns string `envconfig:"PLUGIN_INCLUDE_PATTERNS"`

	// Gradle parameters
	DeployIvyDesc   string `envconfig:"PLUGIN_DEPLOY_IVY_DESC"`
	DeployMavenDesc string `envconfig:"PLUGIN_DEPLOY_MAVEN_DESC"`
	Global          string `envconfig:"PLUGIN_GLOBAL"`
	IvyArtifacts    string `envconfig:"PLUGIN_IVY_ARTIFACTS_PATTERN"`
	IvyDesc         string `envconfig:"PLUGIN_IVY_DESC_PATTERN"`
	RepoDeploy      string `envconfig:"PLUGIN_REPO_DEPLOY"`
	RepoResolve     string `envconfig:"PLUGIN_REPO_RESOLVE"`

	RepoResolverOrDeployerId string `envconfig:"PLUGIN_REPO_RESOLVER_OR_DEPLOYER_ID"`

	UseWrapper  string `envconfig:"PLUGIN_USE_WRAPPER"`
	GradleTasks string `envconfig:"PLUGIN_TASKS"`
	BuildFile   string `envconfig:"PLUGIN_BUILD_FILE"`
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
	case RtMavenDeployer:
		commandsList, err = GetRtMavenDeployerCommandArgs(args)
	case RtMavenResolver:
		commandsList, err = GetRtMavenResolverCommandArgs(args)
	case RtMavenRun:
		commandsList, err = GetMavenRunCommandArgs(args)
	case RtPublishBuildInfo:
		commandsList, err = GetRtPublishBuildInfoCommandArgs(args)
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

var MavenRunCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false, nil, nil},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false, nil, nil},
	{"--detailed-summary=", "PLUGIN_DETAILED_SUMMARY", false, false, nil, nil},
	{"--format=", "PLUGIN_FORMAT", false, false, nil, nil},
	{"--insecure-tls=", "PLUGIN_INSECURE", false, false, nil, nil},
	{"--project=", "PLUGIN_PROJECT", false, false, nil, nil},
	{"--scan=", "PLUGIN_SCAN", false, false, nil, nil},
	{"--threads=", "PLUGIN_THREADS", false, false, nil, nil},
}

var MavenConfigCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--exclude-patterns=", "PLUGIN_EXCLUDE_PATTERNS", false, false, nil, nil},
	{"--global=", "PLUGIN_GLOBAL", false, false, nil, nil},
	{"--include-patterns=", "PLUGIN_INCLUDE_PATTERNS", false, false, nil, nil},
	{"--repo-deploy-releases=", "PLUGIN_REPO_DEPLOY_RELEASES", false, false, nil, nil},
	{"--repo-deploy-snapshots=", "PLUGIN_REPO_DEPLOY_SNAPSHOTS", false, false, nil, nil},
	{"--repo-resolve-releases=", "PLUGIN_REPO_RESOLVE_RELEASES", false, false, nil, nil},
	{"--repo-resolve-snapshots=", "PLUGIN_REPO_RESOLVE_SNAPSHOTS", false, false, nil, nil},
	{"--server-id-deploy=", "PLUGIN_SERVER_ID_DEPLOY", false, false, nil, nil},
	{"--server-id-resolve=", "PLUGIN_SERVER_ID_RESOLVE", false, false, nil, nil},
	{"--use-wrapper=", "PLUGIN_USE_WRAPPER", false, false, nil, nil},
}

func GetMavenCommandArgs(args Args) ([][]string, error) {

	if args.MvnGoals == "" {
		return [][]string{}, fmt.Errorf("Missing mandatory parameter", args.MvnGoals)
	}

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs := GetConfigAddConfigCommandArgs("tmpSrvConfig",
		args.Username, args.Password, args.URL)

	mvnConfigCommandArgs := []string{MvnConfig}
	err := PopulateArgs(&mvnConfigCommandArgs, &args, MavenConfigCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}

	mvnRunCommandArgs := []string{MvnCmd, args.MvnGoals}
	err = PopulateArgs(&mvnRunCommandArgs, &args, MavenRunCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}
	if len(args.MvnPomFile) > 0 {
		mvnRunCommandArgs = append(mvnRunCommandArgs, "-f "+args.MvnPomFile)
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, mvnConfigCommandArgs)
	cmdList = append(cmdList, mvnRunCommandArgs)

	return cmdList, nil
}

var RtMavenConfigCmdJsonTagToExeFlagMapString = []JsonTagToExeFlagMapStringItem{
	{"--exclude-patterns=", "PLUGIN_EXCLUDE_PATTERNS", false, false, nil, nil},
	{"--global=", "PLUGIN_GLOBAL", false, false, nil, nil},
	{"--include-patterns=", "PLUGIN_INCLUDE_PATTERNS", false, false, nil, nil},
	{"--repo-deploy-releases=", "PLUGIN_REPO_DEPLOY_RELEASES", false, false, nil, nil},
	{"--repo-deploy-snapshots=", "PLUGIN_REPO_DEPLOY_SNAPSHOTS", false, false, nil, nil},
	{"--repo-resolve-releases=", "PLUGIN_REPO_RESOLVE_RELEASES", false, false, nil, nil},
	{"--repo-resolve-snapshots=", "PLUGIN_REPO_RESOLVE_SNAPSHOTS", false, false, nil, nil},
	{"--server-id-deploy=", "PLUGIN_DEPLOYER_ID", false, false, nil, nil},
	{"--server-id-resolve=", "PLUGIN_RESOLVER_ID", false, false, nil, nil},
	{"--use-wrapper=", "PLUGIN_USE_WRAPPER", false, false, nil, nil},
}

var RtMavenRunCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false, nil, nil},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false, nil, nil},
	{"--detailed-summary=", "PLUGIN_DETAILED_SUMMARY", false, false, nil, nil},
	{"--format=", "PLUGIN_FORMAT", false, false, nil, nil},
	{"--insecure-tls=", "PLUGIN_INSECURE", false, false, nil, nil},
	{"--project=", "PLUGIN_PROJECT", false, false, nil, nil},
	{"--scan=", "PLUGIN_SCAN", false, false, nil, nil},
	{"--threads=", "PLUGIN_THREADS", false, false, nil, nil},
}

func GetMavenRunCommandArgs(args Args) ([][]string, error) {

	fmt.Println(">>>>>>>>>>>>>> chcbh GetMavenRunCommandArgs READING CONFIGS <<<<<<<<<<<<<<<<")

	if args.MvnGoals == "" {
		return [][]string{}, fmt.Errorf("Missing mandatory parameter", args.MvnGoals)
	}

	var cmdList [][]string

	mvnConfigCommandArgs := []string{MvnConfig}
	err := PopulateArgs(&mvnConfigCommandArgs, &args, RtMavenConfigCmdJsonTagToExeFlagMapString)
	if err != nil {
		return cmdList, err
	}

	err = GetResolverCmd(args, DeployerIdType, "PLUGIN_REPO_DEPLOY_RELEASES", "PLUGIN_REPO_DEPLOY_SNAPSHOTS",
		"--repo-deploy-releases=", "--repo-deploy-snapshots=", &mvnConfigCommandArgs, &cmdList)

	if err != nil {
		return cmdList, err
	}

	err = GetResolverCmd(args, ResolverIdType, "PLUGIN_REPO_RESOLVE_RELEASES", "PLUGIN_REPO_RESOLVE_SNAPSHOTS",
		"--repo-resolve-releases=", "--repo-resolve-snapshots=", &mvnConfigCommandArgs, &cmdList)

	if err != nil {
		return cmdList, err
	}

	mvnRunCommandArgs := []string{MvnCmd, args.MvnGoals}
	err = PopulateArgs(&mvnRunCommandArgs, &args, RtMavenRunCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}
	if len(args.MvnPomFile) > 0 {
		mvnRunCommandArgs = append(mvnRunCommandArgs, "-f "+args.MvnPomFile)
	}

	cmdList = append(cmdList, mvnConfigCommandArgs)
	cmdList = append(cmdList, mvnRunCommandArgs)

	for i, cmd := range cmdList {
		fmt.Println(i, " --> ", cmd)
	}

	return cmdList, nil
}

func GetResolverCmd(args Args, resolverIdType string,
	releaseRepoPluginParam, snapShotRepoPluginParam,
	releaseCliFlag, snapShotCliFlag string,
	mvnConfigCommandArgs *[]string, cmdList *[][]string) error {

	resolverId := ""
	infoFile := ""

	if resolverIdType == DeployerIdType {
		resolverId = args.DeployerId
		infoFile = getDeployerIdFileName(args.DeployerId)
	}

	if resolverIdType == ResolverIdType {
		resolverId = args.ResolverId
		infoFile = getResolverIdFileName(args.ResolverId)
	}

	if resolverId != "" {
		resolverInfo, err := os.ReadFile(infoFile)

		if err == nil {
			var infoMap map[string]interface{}
			err = json.Unmarshal(resolverInfo, &infoMap)
			if err != nil {
				log.Println("Error unmarshalling resolverInfo ", resolverInfo, err.Error())
			}

			if tmpReleaseRepo, ok := infoMap[releaseRepoPluginParam]; ok {
				*mvnConfigCommandArgs = append(*mvnConfigCommandArgs, releaseCliFlag+tmpReleaseRepo.(string))
			}
			if tmpReleaseRepo, ok := infoMap[snapShotRepoPluginParam]; ok {
				*mvnConfigCommandArgs = append(*mvnConfigCommandArgs, snapShotCliFlag+tmpReleaseRepo.(string))
			}

			userName := ""
			if tmpUserName, ok := infoMap["PLUGIN_USERNAME"]; ok {
				userName = tmpUserName.(string)
			} else {
				fmt.Println("Unable to to find  PLUGIN_USERNAME ")
				return fmt.Errorf("Unable to to find  PLUGIN_USERNAME ")
			}

			password := ""
			if tmpPassword, ok := infoMap["PLUGIN_PASSWORD"]; ok {
				password = tmpPassword.(string)
			} else {
				fmt.Println("Unable to to find  PLUGIN_PASSWORD ")
				return fmt.Errorf("Unable to to find  PLUGIN_PASSWORD ")
			}

			url := ""
			if tmpUrl, ok := infoMap["PLUGIN_URL"]; ok {
				url = tmpUrl.(string)
			} else {
				fmt.Println("Unable to to find  PLUGIN_URL ")
				return fmt.Errorf("Unable to to find  PLUGIN_URL ")
			}

			resolverConfigCommand := GetConfigAddConfigCommandArgs(resolverId, userName, password, url)
			*cmdList = append(*cmdList, resolverConfigCommand)

		}
	} else {
		fmt.Println("ResolverId is empty")
		return fmt.Errorf("ResolverId is empty")
	}
	return nil
}

var RtMavenDeployerConfigCmdJsonTagToExeFlagMap = []JsonTagToExeFlagMapStringItem{
	{"--exclude-patterns=", "PLUGIN_EXCLUDE_PATTERNS", false, false, nil, nil},
	{"--include-patterns=", "PLUGIN_INCLUDE_PATTERNS", false, false, nil, nil},
	{"--repo-deploy-releases=", "PLUGIN_REPO_DEPLOY_RELEASES", false, false, nil, nil},
	{"--repo-deploy-snapshots=", "PLUGIN_REPO_DEPLOY_SNAPSHOTS", false, false, nil, nil},
	{"--use-wrapper=", "PLUGIN_USE_WRAPPER", false, false, nil, nil},
	{"--server-id-deploy=", "PLUGIN_DEPLOYER_ID", false, false, nil, nil},
}

func GetRtMavenDeployerCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs := GetConfigAddConfigCommandArgs(args.DeployerId,
		args.Username, args.Password, args.URL)

	rtMavenDeployerMvnConfigCommandArgs := []string{MvnConfig}
	err := PopulateArgs(&rtMavenDeployerMvnConfigCommandArgs, &args, RtMavenDeployerConfigCmdJsonTagToExeFlagMap)
	if err != nil {
		return cmdList, err
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, rtMavenDeployerMvnConfigCommandArgs)

	mvnDeployerMap := map[string]interface{}{
		"PLUGIN_REPO_DEPLOY_RELEASES":  args.MvnDeployReleases,
		"PLUGIN_REPO_DEPLOY_SNAPSHOTS": args.MvnDeploySnapshots,
		"PLUGIN_USERNAME":              args.Username,
		"PLUGIN_PASSWORD":              args.Password,
		"PLUGIN_URL":                   args.URL,
	}
	// convert mvnDeployerMap to json string
	mvnDeployerJson, err := json.Marshal(mvnDeployerMap)
	if err != nil {
		return cmdList, err
	}
	// write mvnDeployerJson to file with name getDeployerIdFileName(args.DeployerId)
	deployerIdFileName := getDeployerIdFileName(args.DeployerId)
	err = os.WriteFile(deployerIdFileName, mvnDeployerJson, 0600)
	if err != nil {
		return cmdList, err
	}

	fmt.Println("=================")
	fmt.Println(jfrogConfigAddConfigCommandArgs)
	fmt.Println(rtMavenDeployerMvnConfigCommandArgs)

	return cmdList, nil
}

func getDeployerIdFileName(deployerId string) string {
	return deployerId + ".deployer_data" + ".json"
}

var RtMavenResolverConfigCmdJsonTagToExeFlagMap = []JsonTagToExeFlagMapStringItem{
	{"--repo-resolve-releases=", "PLUGIN_REPO_RESOLVE_RELEASES", false, false, nil, nil},
	{"--repo-resolve-snapshots=", "PLUGIN_REPO_RESOLVE_SNAPSHOTS", false, false, nil, nil},
	{"--server-id-resolve=", "PLUGIN_RESOLVER_ID", false, false, nil, nil},
}

func GetRtMavenResolverCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs := GetConfigAddConfigCommandArgs(args.ResolverId,
		args.Username, args.Password, args.URL)

	rtMavenDeployerCommandArgs := []string{MvnConfig}
	err := PopulateArgs(&rtMavenDeployerCommandArgs, &args, RtMavenResolverConfigCmdJsonTagToExeFlagMap)
	if err != nil {
		return cmdList, err
	}

	mvnResolverMap := map[string]interface{}{
		"PLUGIN_REPO_RESOLVE_RELEASES":  args.MvnResolveReleases,
		"PLUGIN_REPO_RESOLVE_SNAPSHOTS": args.MvnResolveSnapshots,
		"PLUGIN_USERNAME":               args.Username,
		"PLUGIN_PASSWORD":               args.Password,
		"PLUGIN_URL":                    args.URL,
	}
	// convert mvnResolverMap to json string
	mvnResolverJson, err := json.Marshal(mvnResolverMap)
	if err != nil {
		return cmdList, err
	}
	// write mvnResolverJson to file with name getResolverIdFileName(args.ResolverId)
	resolverIdFileName := getResolverIdFileName(args.ResolverId)
	err = os.WriteFile(resolverIdFileName, mvnResolverJson, 0600)
	if err != nil {
		return cmdList, err
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, rtMavenDeployerCommandArgs)

	fmt.Println(jfrogConfigAddConfigCommandArgs)
	fmt.Println(rtMavenDeployerCommandArgs)

	return cmdList, nil
}

func getResolverIdFileName(deployerId string) string {
	return deployerId + ".resolver_data" + ".json"
}

var RtBuildInfoPublishCmdJsonTagToExeFlagMap = []JsonTagToExeFlagMapStringItem{
	{"--project=", "PLUGIN_PROJECT", false, false, nil, nil},
}

func GetRtPublishBuildInfoCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string
	var jfrogConfigAddConfigCommandArgs []string

	idType := ""
	serverId := ""

	if args.DeployerId != "" {
		idType = DeployerIdType
		serverId = args.DeployerId
	}

	if args.ResolverId != "" {
		idType = ResolverIdType
		serverId = args.ResolverId
	}

	err := GetResolverCmd(args, idType, "None", "None", "", "", &jfrogConfigAddConfigCommandArgs, &cmdList)
	rtPublishBuildInfoCommandArgs := []string{"rt", BuildPublish, args.BuildName, args.BuildNumber,
		"--server-id=" + serverId}
	err = PopulateArgs(&rtPublishBuildInfoCommandArgs, &args, RtBuildInfoPublishCmdJsonTagToExeFlagMap)
	if err != nil {
		return cmdList, err
	}
	cmdList = append(cmdList, rtPublishBuildInfoCommandArgs)

	return cmdList, nil
}

func GetConfigAddConfigCommandArgs(srvConfigStr, userName, password, url string) []string {
	if len(userName) == 0 || len(password) == 0 || len(url) == 0 {
		return []string{}
	}
	if srvConfigStr == "" {
		srvConfigStr = "tmpSrvConfig"
	}
	return []string{"config", "add", srvConfigStr, "--url=" + url,
		"--user=" + userName, "--password=" + password, "--interactive=false"}
}

type JsonTagToExeFlagMapStringItem struct {
	FlagName         string
	PluginArgJsonTag string
	IsMandatory      bool
	StopOnError      bool
	ValidationFunc   func() (bool, error)
	TransformFunc    func() (string, error)
}

func PopulateArgs(tmpCommandsList *[]string, args *Args,
	jsonTagToExeFlagMapStringItemList []JsonTagToExeFlagMapStringItem) error {

	for _, jsonTagToExeFlagMapStringItem := range jsonTagToExeFlagMapStringItemList {
		flagName := jsonTagToExeFlagMapStringItem.FlagName
		pluginArgJsonTag := jsonTagToExeFlagMapStringItem.PluginArgJsonTag
		pluginArgValue, err := GetFieldAddress[*Args, string](args, pluginArgJsonTag)

		if err != nil {
			if jsonTagToExeFlagMapStringItem.IsMandatory || jsonTagToExeFlagMapStringItem.StopOnError {
				fmt.Println("GetFieldAddress error: ", err)
				return err
			}
			fmt.Println("GetFieldAddress error: ", err)
			continue
		}

		if pluginArgValue == nil {
			if jsonTagToExeFlagMapStringItem.IsMandatory || jsonTagToExeFlagMapStringItem.StopOnError {
				fmt.Println("missing mandatory field: ", pluginArgJsonTag)
				return fmt.Errorf("missing mandatory field %s", pluginArgJsonTag)
			}
			fmt.Println("missing mandatory field: ", pluginArgJsonTag)
			continue
		}

		if pluginArgValue == nil &&
			jsonTagToExeFlagMapStringItem.IsMandatory || jsonTagToExeFlagMapStringItem.StopOnError {
			fmt.Println("missing mandatory field: ", pluginArgJsonTag)
			return fmt.Errorf("missing mandatory field %s", pluginArgJsonTag)
		}
		AppendStringArg(tmpCommandsList, flagName, pluginArgValue)
	}

	return nil
}

func AppendStringArg(argsList *[]string, argName string, argValue *string) {

	if argsList == nil {
		fmt.Println("argsList is nil")
		return
	}

	if argValue == nil {
		fmt.Println("argValue is nil")
		return
	}

	if len(*argValue) > 0 {
		*argsList = append(*argsList, argName+*argValue)
	}
}

var tagFieldCache sync.Map

func precomputeTagMapping(structType reflect.Type) map[string]int {
	tagMap := make(map[string]int)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("envconfig")
		if tag != "" {
			tagMap[tag] = i
		}
	}
	return tagMap
}

func getTagMapping(structType reflect.Type) map[string]int {
	if cachedMapping, ok := tagFieldCache.Load(structType); ok {
		return cachedMapping.(map[string]int)
	}

	tagMap := precomputeTagMapping(structType)
	tagFieldCache.Store(structType, tagMap)
	return tagMap
}

func GetFieldAddress[ST, VT any](args ST, argJsonTag string) (*VT, error) {
	v := reflect.ValueOf(args)
	if v.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("args must be a pointer to a struct; got %T", args)
	}
	if v.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("args must point to a struct; got pointer to %s", v.Elem().Kind())
	}

	v = v.Elem()
	t := v.Type()

	tagMap := getTagMapping(t)

	fieldIndex, found := tagMap[argJsonTag]
	if !found {
		return nil, fmt.Errorf("field with tag '%s' not found in struct type '%s'", argJsonTag, t.Name())
	}

	fieldValue := v.Field(fieldIndex)
	if fieldValue.CanAddr() {
		if fieldValue.Type().AssignableTo(reflect.TypeOf((*VT)(nil)).Elem()) {
			return fieldValue.Addr().Interface().(*VT), nil
		}
		return nil, fmt.Errorf("field with tag '%s' in struct '%s' is not of type '%T'; actual type is '%s'",
			argJsonTag, t.Name(), new(VT), fieldValue.Type().String())
	}

	return nil, fmt.Errorf("field with tag '%s' in struct '%s' cannot be addressed", argJsonTag, t.Name())
}

const (
	MvnCmd             = "mvn"
	MvnConfig          = "mvn-config"
	BuildPublish       = "build-publish"
	RtMavenDeployer    = "rtMavenDeployer"
	RtMavenResolver    = "rtMavenResolver"
	RtMavenRun         = "rtMavenRun"
	RtPublishBuildInfo = "rtPublishBuildInfo"
	DeployerIdType     = "deployer"
	ResolverIdType     = "resolver"
)
