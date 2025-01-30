package plugin

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

const (
	MvnCmd       = "mvn"
	MvnConfig    = "mvn-config"
	BuildPublish = "build-publish"
	Deploy       = "deploy"
	Publish      = "publish"
	GradleConfig = "gradle-config"
	GradleCmd    = "gradle"
	tmpServerId  = "tmpServerId"
)

func HandleRtCommands(args Args) error {

	commandsList, err := GetRtCommandsList(args)
	if err != nil {
		log.Println("Error Unable to get rt commands list err = ", err)
		return err
	}

	err = WriteKnownGoodServerCertsForTls(args)
	if err != nil {
		log.Println("Error Unable to write TLS certs err = ", err)
		return err
	}

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

func WriteKnownGoodServerCertsForTls(args Args) error {

	insecure := parseBoolOrDefault(false, args.Insecure)
	if insecure {
		return nil
	}

	// create pem file
	if args.PEMFileContents != "" {
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
		log.Printf("Creating pem file at %q\n", path)
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
			log.Printf("Successfully created pem file at %q\n", path)
		}
	}
	return nil
}

func GetRtCommandsList(args Args) ([][]string, error) {
	log.Println("Handling rt command handleRtCommand")
	commandsList := [][]string{}
	var err error

	log.Println("Checking GetRtCommandsList args.Command ", args.Command)

	if args.BuildTool == MvnCmd && (args.Command == "" || args.Command == "build") {
		log.Println("mvn build start")
		commandsList, err = GetMavenBuildCommandArgs(args)
	}

	if args.BuildTool == MvnCmd && args.Command == "publish" {
		commandsList, err = GetMavenPublishCommand(args)
	}

	if args.BuildTool == GradleCmd && (args.Command == "" || args.Command == "build") {
		log.Println("Gradle build start")
		commandsList, err = GetGradleCommandArgs(args)
	}

	if args.BuildTool == GradleCmd && args.Command == "publish" {
		log.Println("Gradle publish start")
		commandsList, err = GetGradlePublishCommand(args)
	}

	if args.Command == "download" {
		log.Println("download start")
		commandsList, err = GetDownloadCommandArgs(args)
	}

	if args.Command == "cleanup" {
		log.Println("cleanup start")
		commandsList, err = GetCleanupCommandArgs(args)
	}

	if args.Command == "scan" {
		log.Println("scan start")
		commandsList, err = GetScanCommandArgs(args)
	}

	if args.Command == "create-build-info" {
		log.Println("create-build-info start")
		commandsList, err = GetCreateBuildInfoCommandArgs(args)
	}

	if args.Command == "publish-build-info" {
		log.Println("publish-build-info start")
		commandsList, err = GetBuildInfoPublishCommandArgs(args)
	}

	if args.Command == "promote" {
		log.Println("promote start")
		commandsList, err = GetPromoteCommandArgs(args)
	}

	return commandsList, err
}

func GetShellForOs(osName string) (string, string) {
	if osName == "windows" {
		return os.Getenv("POWERSHELL_EXE"), "-Command"
	}

	return "sh", "-c"
}

func ExecCommand(args Args, cmdArgs []string) error {

	cmdStr := strings.Join(cmdArgs[:], " ")

	shell, shArg := GetShellForOs(runtime.GOOS)

	log.Println()
	log.Printf("%s %s %s", shell, shArg, cmdStr)
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
			log.Println("Error publishing build info: ", err)
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
		log.Println("setAuthParams error: ", err)
		return []string{""}, err
	}

	cfgCommand := []string{"config", "add", srvConfigStr, "--url=" + url}
	cfgCommand = append(cfgCommand, authParams...)
	cfgCommand = append(cfgCommand, "--interactive=false")
	return cfgCommand, nil
}
