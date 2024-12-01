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

func GetGradleCommandArgs(userName, password, url, repoResolve, repoDeploy,
	gradleTasks, buildName, buildNumber string,
	numThreads int, projectKey, otherOpts string) ([][]string, error) {

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs := GetConfigAddConfigCommandArgs(userName, password, url)

	gradleConfigCommandArgs := []string{"gradle-config",
		"--repo-resolve=" + repoResolve, "--repo-deploy=" + repoDeploy}
	gradleTaskCommandArgs := []string{"gradle", gradleTasks}

	if len(buildName) > 0 {
		gradleTaskCommandArgs = append(gradleTaskCommandArgs, "--build-name="+buildName)
	}
	if len(buildNumber) > 0 {
		gradleTaskCommandArgs = append(gradleTaskCommandArgs, "--build-number="+buildNumber)
	}
	if numThreads > 0 {
		gradleTaskCommandArgs = append(gradleTaskCommandArgs, fmt.Sprintf("--threads=%d", numThreads))
	}
	if len(projectKey) > 0 {
		gradleTaskCommandArgs = append(gradleTaskCommandArgs, "--project="+projectKey)
	}
	gradleTaskCommandArgs = append(gradleTaskCommandArgs, otherOpts)

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
	srvConfigStr := "tmpSrvConfig"
	return []string{"config", "add", srvConfigStr, "--url=" + url,
		"--user=" + userName, "--password=" + password, "--interactive=false"}
}

const (
	MvnCmd      = "mvn"
	MvnConfig   = "mvn-config"
	GradleCmd   = "gradle"
	UploadCmd   = "upload"
	DownloadCmd = "download"
)

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
