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

const (
	MvnCmd       = "mvn"
	MvnConfig    = "mvn-config"
	BuildPublish = "build-publish"
)

type RtCommand struct {
	Command   string `envconfig:"PLUGIN_COMMAND"`
	BuildTool string `envconfig:"PLUGIN_BUILD_TOOL"`
	RtMvnCommand
}

type RtMvnCommand struct {
	ResolveReleaseRepo  string `envconfig:"PLUGIN_RESOLVE_RELEASE_REPO"`
	ResolveSnapshotRepo string `envconfig:"PLUGIN_RESOLVE_SNAPSHOT_REPO"`
	DeployReleaseRepo   string `envconfig:"PLUGIN_DEPLOY_RELEASE_REPO"`
	DeploySnapshotRepo  string `envconfig:"PLUGIN_DEPLOY_SNAPSHOT_REPO"`
	DeployRepo          string `envconfig:"PLUGIN_DEPLOY_REPO"`

	MvnGoals     string `envconfig:"PLUGIN_GOALS"`
	MvnPomFile   string `envconfig:"PLUGIN_POM_FILE"`
	ProjectKey   string `envconfig:"PLUGIN_PROJECT_KEY"`
	OptionalArgs string `envconfig:"PLUGIN_OPTIONAL_ARGS"`

	DeployerId string `envconfig:"PLUGIN_DEPLOYER_ID"`
	ResolverId string `envconfig:"PLUGIN_RESOLVER_ID"`
}

func HandleRtCommands(ctx context.Context, args Args) error {
	commandsList, err := GetRtCommandsList(ctx, args)
	for _, cmd := range commandsList {
		execArgs := []string{getJfrogBin()}
		execArgs = append(execArgs, cmd...)
		err := ExecCommand(args, execArgs)
		if err != nil {
			log.Println("Error Unable to run err = ", err)
			return err
		}
	}

	return err
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
