package plugin

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"sync"
)

func ExecCommand(args Args, cmdArgs []string) error {

	cmdStr := strings.Join(cmdArgs[:], " ")

	shell, shArg := getShell()

	log.Println()
	fmt.Printf("%s %s %s", shell, shArg, cmdStr)
	log.Println()

	cmd := exec.Command(shell, shArg, cmdStr)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "JFROG_CLI_OFFER_CONFIG=false")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)

	err := cmd.Run()
	if err != nil {
		log.Println(" Error: ", err)
		return err
	}

	if args.PublishBuildInfo {
		if err := publishBuildInfo(args); err != nil {
			fmt.Println("Error publishing build info: ", err)
			return err
		}
	}

	return nil
}

func GetRtCommandsList(ctx context.Context, args Args) ([][]string, error) {
	log.Println("Handling rt command handleRtCommand")
	commandsList := [][]string{}
	var err error

	if args.BuildTool == MvnCmd && (args.Command == "" || args.Command == "build") {
		log.Println("first handlers")
		commandsList, err = GetMavenBuildCommandArgs(args)
	}

	if args.BuildTool == MvnCmd && args.Command == "publish" {
		commandsList, err = GetMavenPublishCommand(args)
	}

	return commandsList, err
}

func GetMavenBuildCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs(args.ResolverId,
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		return cmdList, err
	}

	mvnConfigCommandArgs := []string{MvnConfig}
	err = PopulateArgs(&mvnConfigCommandArgs, &args, MavenConfigCmdJsonTagToExeFlagMapStringItemList)
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

func GetMavenPublishCommand(args Args) ([][]string, error) {

	var cmdList [][]string
	var jfrogConfigAddConfigCommandArgs []string

	tmpServerId := "tmpSrvConfig"
	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs(tmpServerId,
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		log.Println("GetConfigAddConfigCommandArgs error: ", err)
		return cmdList, err
	}

	mvnConfigCommandArgs := []string{MvnConfig}
	err = PopulateArgs(&mvnConfigCommandArgs, &args, MavenConfigCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		log.Println("PopulateArgs error: ", err)
		return cmdList, err
	}

	rtPublishBuildInfoCommandArgs := []string{"rt", BuildPublish, args.BuildName, args.BuildNumber,
		"--server-id=" + tmpServerId}
	err = PopulateArgs(&rtPublishBuildInfoCommandArgs, &args, RtBuildInfoPublishCmdJsonTagToExeFlagMap)
	if err != nil {
		log.Println("PopulateArgs error: ", err)
		return cmdList, err
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, mvnConfigCommandArgs)
	cmdList = append(cmdList, rtPublishBuildInfoCommandArgs)

	return cmdList, nil
}

var RtBuildInfoPublishCmdJsonTagToExeFlagMap = []JsonTagToExeFlagMapStringItem{
	{"--project=", "PLUGIN_PROJECT", false, false},
}

func GetConfigAddConfigCommandArgs(srvConfigStr, userName, password, url,
	accessToken, apiKey string) ([]string, error) {

	if srvConfigStr == "" {
		srvConfigStr = "tmpSrvConfig"
	}

	authParams, err := setAuthParams([]string{}, Args{Username: userName,
		Password: password, AccessToken: accessToken, APIKey: apiKey})
	if err != nil {
		fmt.Println("setAuthParams error: ", err)
		return []string{""}, err
	}

	cfgCommand := []string{"config", "add", srvConfigStr, "--url=" + url}
	cfgCommand = append(cfgCommand, authParams...)
	cfgCommand = append(cfgCommand, "--interactive=false")
	return cfgCommand, nil
}

var MavenRunCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false},
	{"--detailed-summary=", "PLUGIN_DETAILED_SUMMARY", false, false},
	{"--format=", "PLUGIN_FORMAT", false, false},
	{"--insecure-tls=", "PLUGIN_INSECURE", false, false},
	{"--project=", "PLUGIN_PROJECT", false, false},
	{"--scan=", "PLUGIN_SCAN", false, false},
	{"--threads=", "PLUGIN_THREADS", false, false},
}

var MavenConfigCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--exclude-patterns=", "PLUGIN_EXCLUDE_PATTERNS", false, false},
	{"--global=", "PLUGIN_GLOBAL", false, false},
	{"--include-patterns=", "PLUGIN_INCLUDE_PATTERNS", false, false},
	{"--repo-deploy-releases=", "PLUGIN_DEPLOY_RELEASE_REPO", false, false},
	{"--repo-deploy-snapshots=", "PLUGIN_DEPLOY_SNAPSHOT_REPO", false, false},
	{"--repo-resolve-releases=", "PLUGIN_RESOLVE_RELEASE_REPO", false, false},
	{"--repo-resolve-snapshots=", "PLUGIN_RESOLVE_SNAPSHOT_REPO", false, false},
	{"--server-id-deploy=", "PLUGIN_SERVER_ID_DEPLOY", false, false},
	{"--server-id-resolve=", "PLUGIN_RESOLVER_ID", false, false},
	{"--use-wrapper=", "PLUGIN_USE_WRAPPER", false, false},
}

type JsonTagToExeFlagMapStringItem struct {
	FlagName         string
	PluginArgJsonTag string
	IsMandatory      bool
	StopOnError      bool
}

func PopulateArgs(tmpCommandsList *[]string, args *Args,
	jsonTagToExeFlagMapStringItemList []JsonTagToExeFlagMapStringItem) error {

	for _, jsonTagToExeFlagMapStringItem := range jsonTagToExeFlagMapStringItemList {
		flagName := jsonTagToExeFlagMapStringItem.FlagName
		pluginArgJsonTag := jsonTagToExeFlagMapStringItem.PluginArgJsonTag
		pluginArgValue, err := GetFieldAddress[*Args, string](args, pluginArgJsonTag)

		if err != nil {
			if jsonTagToExeFlagMapStringItem.IsMandatory || jsonTagToExeFlagMapStringItem.StopOnError {
				log.Println("GetFieldAddress error: ", err)
				return err
			}
			//log.Println("GetFieldAddress error: ", err)
			continue
		}

		if pluginArgValue == nil {
			if jsonTagToExeFlagMapStringItem.IsMandatory || jsonTagToExeFlagMapStringItem.StopOnError {
				log.Println("missing mandatory field: ", pluginArgJsonTag)
				return fmt.Errorf("missing mandatory field %s", pluginArgJsonTag)
			}
			log.Println("missing mandatory field: ", pluginArgJsonTag)
			continue
		}

		if pluginArgValue == nil &&
			jsonTagToExeFlagMapStringItem.IsMandatory || jsonTagToExeFlagMapStringItem.StopOnError {
			log.Println("missing mandatory field: ", pluginArgJsonTag)
			return fmt.Errorf("missing mandatory field %s", pluginArgJsonTag)
		}
		AppendStringArg(tmpCommandsList, flagName, pluginArgValue)
	}

	return nil
}

func AppendStringArg(argsList *[]string, argName string, argValue *string) {

	if argsList == nil {
		log.Println("argsList is nil")
		return
	}

	if argValue == nil {
		log.Println("argValue is nil")
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
