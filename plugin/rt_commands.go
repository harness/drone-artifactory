package plugin

import (
	"fmt"
	"reflect"
	"sync"
)

var MavenRunCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false, nil, nil},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false, nil, nil},
	{"--detailed-summary=", "PLUGIN_DETAILED_SUMMARY", false, false, nil, nil},
	{"--format=", "PLUGIN_FORMAT", false, false, nil, nil},
	{"--insecure-tls=", "PLUGIN_INSECURE_TLS", false, false, nil, nil},
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

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs := GetConfigAddConfigCommandArgs(args.Username, args.Password, args.URL)

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

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, mvnConfigCommandArgs)
	cmdList = append(cmdList, mvnRunCommandArgs)

	return cmdList, nil
}

/*
Options:
  --deploy-ivy-desc          [Default: true] Set to false if you do not wish to deploy Ivy descriptors.
  --deploy-maven-desc        [Default: true] Set to false if you do not wish to deploy Maven descriptors.
  --global                   [Default: false] Set to true if you'd like the configuration to be global (for all projects). Specific projects can override the global configuration.
  --ivy-artifacts-pattern    [Default: '[organization]/[module]/[revision]/[artifact]-[revision](-[classifier]).[ext]' Set the deployed Ivy artifacts pattern.
  --ivy-desc-pattern         [Default: '[organization]/[module]/ivy-[revision].xml' Set the deployed Ivy descriptor pattern.
  --repo-deploy              [Optional] Repository for artifacts deployment.
  --repo-resolve             [Optional] Repository for dependencies resolution.
  --server-id-deploy         [Optional] Artifactory server ID for deployment. The server should be configured using the 'jfrog c add' command.
  --server-id-resolve        [Optional] Artifactory server ID for resolution. The server should be configured using the 'jfrog c add' command.
  --use-wrapper              [Default: false] Set to true if you wish to use the wrapper.
  --uses-plugin              [Default: false] Set to true if the Gradle Artifactory Plugin is already applied in the build script.
*/

var GradleConfigJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--deploy-ivy-desc=", "PLUGIN_DEPLOY_IVY_DESC", false, false, nil, nil},
	{"--deploy-maven-desc=", "PLUGIN_DEPLOY_MAVEN_DESC", false, false, nil, nil},
	{"--global=", "PLUGIN_GLOBAL", false, false, nil, nil},
	{"--ivy-artifacts-pattern=", "PLUGIN_IVY_ARTIFACTS_PATTERN", false, false, nil, nil},
	{"--ivy-desc-pattern=", "PLUGIN_IVY_DESC_PATTERN", false, false, nil, nil},
	{"--repo-deploy=", "PLUGIN_REPO_DEPLOY", false, false, nil, nil},
	{"--repo-resolve=", "PLUGIN_REPO_RESOLVE", false, false, nil, nil},
	{"--server-id-deploy=", "PLUGIN_SERVER_ID_DEPLOY", false, false, nil, nil},
	{"--server-id-resolve=", "PLUGIN_SERVER_ID_RESOLVE", false, false, nil, nil},
	{"--use-wrapper=", "PLUGIN_USE_WRAPPER", false, false, nil, nil},
	{"--uses-plugin=", "PLUGIN_USES_PLUGIN", false, false, nil, nil},
}

/*
Options:

	--build-name          [Optional] Providing this option will collect and record build info for this build name. Build number option is mandatory when this option is provided.
	--build-number        [Optional] Providing this option will collect and record build info for this build number. Build name option is mandatory when this option is provided.
	--detailed-summary    [Default: false] Set to true to include a list of the affected files in the command summary.
	--format              [Default: table] Defines the output format of the command. Acceptable values are: table, json, simple-json and sarif. Note: the json format doesn't include information about scans that are included as part of the Advanced Security package.
	--project             [Optional] JFrog Artifactory project key.
	--scan                [Default: false] Set if you'd like all files to be scanned by Xray on the local file system prior to the upload, and skip the upload if any of the files are found vulnerable.
	--threads             [Default: 3] Number of threads for uploading build artifacts.
*/
var GradleRunJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false, nil, nil},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false, nil, nil},
	{"--detailed-summary=", "PLUGIN_DETAILED_SUMMARY", false, false, nil, nil},
	{"--format=", "PLUGIN_FORMAT", false, false, nil, nil},
	{"--project=", "PLUGIN_PROJECT", false, false, nil, nil},
	{"--scan=", "PLUGIN_SCAN", false, false, nil, nil},
	{"--threads=", "PLUGIN_THREADS", false, false, nil, nil},
}

func GetGradleCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs := GetConfigAddConfigCommandArgs(args.Username, args.Password, args.URL)

	gradleConfigCommandArgs := []string{GradleConfig}
	err := PopulateArgs(&gradleConfigCommandArgs, &args, GradleConfigJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}

	gradleTaskCommandArgs := []string{GradleCmd, args.GradleTasks}
	err = PopulateArgs(&gradleTaskCommandArgs, &args, GradleRunJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, gradleConfigCommandArgs)
	cmdList = append(cmdList, gradleTaskCommandArgs)

	return cmdList, nil
}

type JsonTagToExeFlagMapStringItem struct {
	FlagName         string
	PluginArgJsonTag string
	IsMandatory      bool
	StopOnError      bool
	ValidationFunc   func() (bool, error)
	TransformFunc    func() (string, error)
}

var UploadCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--access-token=", "PLUGIN_ACCESS_TOKEN", false, false, nil, nil},
	{"--ant=", "PLUGIN_ANT", false, false, nil, nil},
	{"--archive=", "PLUGIN_ARCHIVE", false, false, nil, nil},
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false, nil, nil},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false, nil, nil},
	{"--chunk-size=", "PLUGIN_CHUNK_SIZE", false, false, nil, nil},
	{"--client-cert-key-path=", "PLUGIN_CLIENT_CERT_KEY_PATH", false, false, nil, nil},
	{"--client-cert-path=", "PLUGIN_CLIENT_CERT_PATH", false, false, nil, nil},
	{"--deb=", "PLUGIN_DEB", false, false, nil, nil},
	{"--detailed-summary=", "PLUGIN_DETAILED_SUMMARY", false, false, nil, nil},
	{"--dry-run=", "PLUGIN_DRY_RUN", false, false, nil, nil},
	{"--exclusions=", "PLUGIN_EXCLUSIONS", false, false, nil, nil},
	{"--explode=", "PLUGIN_EXPLODE", false, false, nil, nil},
	{"--fail-no-op=", "PLUGIN_FAIL_NO_OP", false, false, nil, nil},
	{"--include-dirs=", "PLUGIN_INCLUDE_DIRS", false, false, nil, nil},
	{"--insecure-tls=", "PLUGIN_INSECURE_TLS", false, false, nil, nil},
	{"--min-split=", "PLUGIN_MIN_SPLIT", false, false, nil, nil},
	{"--module=", "PLUGIN_MODULE", false, false, nil, nil},
	{"--project=", "PLUGIN_PROJECT", false, false, nil, nil},
	{"--quiet=", "PLUGIN_QUIET", false, false, nil, nil},
	{"--recursive=", "PLUGIN_RECURSIVE", false, false, nil, nil},
	{"--regexp=", "PLUGIN_REGEXP", false, false, nil, nil},
	{"--retry-wait-time=", "PLUGIN_RETRY_WAIT_TIME", false, false, nil, nil},
	{"--server-id=", "PLUGIN_SERVER_ID", false, false, nil, nil},
	{"--split-count=", "PLUGIN_SPLIT_COUNT", false, false, nil, nil},
	{"--ssh-key-path=", "PLUGIN_SSH_KEY_PATH", false, false, nil, nil},
	{"--ssh-passphrase=", "PLUGIN_SSH_PASSPHRASE", false, false, nil, nil},
	{"--symlinks=", "PLUGIN_SYMLINKS", false, false, nil, nil},
	{"--sync-deletes=", "PLUGIN_SYNC_DELETES", false, false, nil, nil},
}

func GetUploadCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs := GetConfigAddConfigCommandArgs(args.Username, args.Password, args.URL)

	uploadCommandArgs := []string{"rt", "upload", args.Source, args.Target}
	err := PopulateArgs(&uploadCommandArgs, &args, UploadCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, uploadCommandArgs)

	return cmdList, nil
}

var DownloadCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--access-token=", "PLUGIN_ACCESS_TOKEN", false, false, nil, nil},
	{"--archive-entries=", "PLUGIN_ARCHIVE_ENTRIES", false, false, nil, nil},
	{"--build=", "PLUGIN_BUILD", false, false, nil, nil},
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false, nil, nil},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false, nil, nil},
	{"--bundle=", "PLUGIN_BUNDLE", false, false, nil, nil},
	{"--bypass-archive-inspection=", "PLUGIN_BYPASS_ARCHIVE_INSPECTION", false, false, nil, nil},
	{"--client-cert-key-path=", "PLUGIN_CLIENT_CERT_KEY_PATH", false, false, nil, nil},
	{"--client-cert-path=", "PLUGIN_CLIENT_CERT_PATH", false, false, nil, nil},
	{"--detailed-summary=", "PLUGIN_DETAILED_SUMMARY", false, false, nil, nil},
	{"--dry-run=", "PLUGIN_DRY_RUN", false, false, nil, nil},
	{"--exclude-artifacts=", "PLUGIN_EXCLUDE_ARTIFACTS", false, false, nil, nil},
	{"--exclude-props=", "PLUGIN_EXCLUDE_PROPS", false, false, nil, nil},
	{"--exclusions=", "PLUGIN_EXCLUSIONS", false, false, nil, nil},
	{"--explode=", "PLUGIN_EXPLODE", false, false, nil, nil},
	{"--fail-no-op=", "PLUGIN_FAIL_NO_OP", false, false, nil, nil},
	{"--flat=", "PLUGIN_FLAT", false, false, nil, nil},
	{"--gpg-key=", "PLUGIN_GPG_KEY", false, false, nil, nil},
	{"--include-deps=", "PLUGIN_INCLUDE_DEPS", false, false, nil, nil},
	{"--include-dirs=", "PLUGIN_INCLUDE_DIRS", false, false, nil, nil},
	{"--insecure-tls=", "PLUGIN_INSECURE_TLS", false, false, nil, nil},
	{"--limit=", "PLUGIN_LIMIT", false, false, nil, nil},
	{"--min-split=", "PLUGIN_MIN_SPLIT", false, false, nil, nil},
	{"--module=", "PLUGIN_MODULE", false, false, nil, nil},
	{"--offset=", "PLUGIN_OFFSET", false, false, nil, nil},
	{"--project=", "PLUGIN_PROJECT", false, false, nil, nil},
	{"--props=", "PLUGIN_PROPS", false, false, nil, nil},
	{"--quiet=", "PLUGIN_QUIET", false, false, nil, nil},
	{"--recursive=", "PLUGIN_RECURSIVE", false, false, nil, nil},
	{"--retries=", "PLUGIN_RETRIES", false, false, nil, nil},
	{"--retry-wait-time=", "PLUGIN_RETRY_WAIT_TIME", false, false, nil, nil},
	{"--server-id=", "PLUGIN_SERVER_ID", false, false, nil, nil},
	{"--skip-checksum=", "PLUGIN_SKIP_CHECKSUM", false, false, nil, nil},
	{"--sort-by=", "PLUGIN_SORT_BY", false, false, nil, nil},
	{"--sort-order=", "PLUGIN_SORT_ORDER", false, false, nil, nil},
	{"--spec=", "PLUGIN_SPEC", false, false, nil, nil},
	{"--spec-vars=", "PLUGIN_SPEC_VARS", false, false, nil, nil},
	{"--split-count=", "PLUGIN_SPLIT_COUNT", false, false, nil, nil},
	{"--ssh-key-path=", "PLUGIN_SSH_KEY_PATH", false, false, nil, nil},
	{"--ssh-passphrase=", "PLUGIN_SSH_PASSPHRASE", false, false, nil, nil},
	{"--sync-deletes=", "PLUGIN_SYNC_DELETES", false, false, nil, nil},
	{"--threads=", "PLUGIN_THREADS", false, false, nil, nil},
	{"--validate-symlinks=", "PLUGIN_VALIDATE_SYMLINKS", false, false, nil, nil},
}

func GetDownloadCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs := GetConfigAddConfigCommandArgs(args.Username, args.Password, args.URL)

	downloadCommandArgs := []string{"rt", "download", args.Target, args.Source}
	err := PopulateArgs(&downloadCommandArgs, &args, DownloadCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, downloadCommandArgs)

	return cmdList, nil
}

var CleanupCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--access-token=", "PLUGIN_ACCESS, TOKEN", false, false, nil, nil},
	{"--archive-entries=", "PLUGIN_ARCHIVE_ENTRIES", false, false, nil, nil},
	{"--build=", "PLUGIN_BUILD", false, false, nil, nil},
	{"--client-cert-key-path=", "PLUGIN_CLIENT_CERT_KEY_PATH", false, false, nil, nil},
	{"--client-cert-path=", "PLUGIN_CLIENT_CERT_PATH", false, false, nil, nil},
	{"--dry-run=", "PLUGIN_DRY_RUN", false, false, nil, nil},
	{"--exclude-artifacts=", "PLUGIN_EXCLUDE_ARTIFACTS", false, false, nil, nil},
	{"--exclude-props=", "PLUGIN_EXCLUDE_PROPS", false, false, nil, nil},
	{"--exclusions=", "PLUGIN_EXCLUSIONS", false, false, nil, nil},
	{"--fail-no-op=", "PLUGIN_FAIL_NO_OP", false, false, nil, nil},
	{"--include-deps=", "PLUGIN_INCLUDE_DEPS", false, false, nil, nil},
	{"--insecure-tls=", "PLUGIN_INSECURE_TLS", false, false, nil, nil},
	{"--limit=", "PLUGIN_LIMIT", false, false, nil, nil},
	{"--offset=", "PLUGIN_OFFSET", false, false, nil, nil},
	{"--project=", "PLUGIN_PROJECT", false, false, nil, nil},
	{"--props=", "PLUGIN_PROPS", false, false, nil, nil},
	{"--quiet=", "PLUGIN_QUIET", false, false, nil, nil},
	{"--recursive=", "PLUGIN_RECURSIVE", false, false, nil, nil},
	{"--retries=", "PLUGIN_RETRIES", false, false, nil, nil},
	{"--retry-wait-time=", "PLUGIN_RETRY_WAIT_TIME", false, false, nil, nil},
	{"--server-id=", "PLUGIN_SERVER_ID", false, false, nil, nil},
	{"--sort-by=", "PLUGIN_SORT_BY", false, false, nil, nil},
	{"--sort-order=", "PLUGIN_SORT_ORDER", false, false, nil, nil},
	{"--spec=", "PLUGIN_SPEC", false, false, nil, nil},
	{"--spec-vars=", "PLUGIN_SPEC_VARS", false, false, nil, nil},
	{"--ssh-key-path=", "PLUGIN_SSH_KEY_PATH", false, false, nil, nil},
	{"--ssh-passphrase=", "PLUGIN_SSH_PASSPHRASE", false, false, nil, nil},
	{"--threads=", "PLUGIN_THREADS", false, false, nil, nil},
}

func GetCleanupCommandArgs(args Args) ([][]string, error) {
	var cmdList [][]string
	jfrogConfigAddConfigCommandArgs := GetConfigAddConfigCommandArgs(args.Username, args.Password, args.URL)
	cleanupCommandArgs := []string{"rt", "del", args.Target}
	err := PopulateArgs(&cleanupCommandArgs, &args, CleanupCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}
	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, cleanupCommandArgs)
	return cmdList, nil
}

func GetBuildInfoCommandArgs(args Args) ([][]string, error) {
	var cmdList [][]string
	jfrogConfigAddConfigCommandArgs := GetConfigAddConfigCommandArgs(args.Username, args.Password, args.URL)
	buildCollectEnvCommandArgs := []string{"rt", "build-collect-env", args.BuildName, args.BuildNumber}
	buildInfoCommandArgs := []string{"rt", "build-publish", args.BuildName, args.BuildNumber}
	err := PopulateArgs(&buildInfoCommandArgs, &args, nil)
	if err != nil {
		return cmdList, err
	}
	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, buildCollectEnvCommandArgs)
	cmdList = append(cmdList, buildInfoCommandArgs)
	return cmdList, nil
}

var PromoteCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--access-token=", "PLUGIN_ACCESS, TOKEN", false, false, nil, nil},
	{"--comment=", "PLUGIN_COMMENT", false, false, nil, nil},
	{"--copy=", "PLUGIN_COPY", false, false, nil, nil},
	{"--dry-run=", "PLUGIN_DRY_RUN", false, false, nil, nil},
	{"--fail-fast=", "PLUGIN_FAIL_FAST", false, false, nil, nil},
	{"--include-dependencies=", "PLUGIN_INCLUDE_DEPENDENCIES", false, false, nil, nil},
	{"--insecure-tls=", "PLUGIN_INSECURE_TLS", false, false, nil, nil},
	{"--project=", "PLUGIN_PROJECT", false, false, nil, nil},
	{"--props=", "PLUGIN_PROPS", false, false, nil, nil},
	{"--server-id=", "PLUGIN_SERVER_ID", false, false, nil, nil},
	{"--source-repo=", "PLUGIN_SOURCE_REPO", false, false, nil, nil},
	{"--ssh-key-path=", "PLUGIN_SSH_KEY_PATH", false, false, nil, nil},
	{"--ssh-passphrase=", "PLUGIN_SSH_PASSPHRASE", false, false, nil, nil},
	{"--status=", "PLUGIN_STATUS", false, false, nil, nil},
}

func GetPromoteCommandArgs(args Args) ([][]string, error) {
	var cmdList [][]string
	jfrogConfigAddConfigCommandArgs := GetConfigAddConfigCommandArgs(args.Username, args.Password, args.URL)
	promoteCommandArgs := []string{"rt", "build-promote", args.BuildName, args.BuildNumber, args.Target}
	err := PopulateArgs(&promoteCommandArgs, &args, PromoteCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}
	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, promoteCommandArgs)
	return cmdList, nil
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

func GetConfigAddConfigCommandArgs(userName, password, url string) []string {
	if len(userName) == 0 || len(password) == 0 || len(url) == 0 {
		return []string{}
	}
	srvConfigStr := "tmpSrvConfig"
	return []string{"config", "add", srvConfigStr, "--url=" + url,
		"--user=" + userName, "--password=" + password, "--interactive=false"}
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
	MvnCmd       = "mvn"
	MvnConfig    = "mvn-config"
	GradleCmd    = "gradle"
	GradleConfig = "gradle-config"
	UploadCmd    = "upload"
	DownloadCmd  = "download"
	CleanUpCmd   = "cleanup"
	BuildInfoCmd = "build-info"
	PromoteCmd   = "promote"
)
