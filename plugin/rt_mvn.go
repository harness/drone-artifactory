package plugin

import (
	"fmt"
	"reflect"
	"sync"
)

func GetMavenCommandArgs(userName, password, url,
	repoResolveReleases, repoResolveSnapshots, repoDeployReleases, repoDeploySnapshots,
	pomFile, goals, buildName, buildNumber string, numThreads int,
	insecureTls string, projectKey string,
	otherOpts string) ([][]string, error) {

	var mvnCmdList [][]string

	if len(userName) == 0 || len(password) == 0 {
		return mvnCmdList, fmt.Errorf("missing username or password")
	}

	if len(url) == 0 {
		return mvnCmdList, fmt.Errorf("missing url")
	}

	if len(pomFile) == 0 {
		pomFile = "pom.xml"
	}

	if len(goals) == 0 {
		goals = "clean install"
	}

	jfrogConfigAddConfigCommandArgs := GetConfigAddConfigCommandArgs(userName, password, url)

	mvnConfigCommandArgs := []string{"mvn-config", "--global", "--repo-resolve-releases=" + repoResolveReleases,
		"--repo-resolve-snapshots=" + repoResolveSnapshots,
		"--repo-deploy-releases=" + repoDeployReleases, "--repo-deploy-snapshots=" + repoDeploySnapshots}

	mvnGoalCommandArgs := []string{"mvn", goals}

	if len(pomFile) > 0 {
		mvnGoalCommandArgs = append(mvnGoalCommandArgs, "-file="+pomFile)
	}
	if len(buildName) > 0 {
		mvnGoalCommandArgs = append(mvnGoalCommandArgs, "--build-name="+buildName)
	}
	if len(buildNumber) > 0 {
		mvnGoalCommandArgs = append(mvnGoalCommandArgs, "--build-number="+buildNumber)
	}
	if numThreads > 0 {
		mvnGoalCommandArgs = append(mvnGoalCommandArgs, fmt.Sprintf("--threads=%d", numThreads))
	}
	if len(insecureTls) > 0 {
		mvnGoalCommandArgs = append(mvnGoalCommandArgs, "--insecure-tls="+insecureTls)
	}
	if len(projectKey) > 0 {
		mvnGoalCommandArgs = append(mvnGoalCommandArgs, "--project="+projectKey)
	}
	mvnGoalCommandArgs = append(mvnGoalCommandArgs, otherOpts)

	mvnCmdList = append(mvnCmdList, jfrogConfigAddConfigCommandArgs)
	mvnCmdList = append(mvnCmdList, mvnConfigCommandArgs)
	mvnCmdList = append(mvnCmdList, mvnGoalCommandArgs)

	return mvnCmdList, nil
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
	MvnCmd    = "mvn"
	GradleCmd = "gradle"
	UploadCmd = "upload"
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
