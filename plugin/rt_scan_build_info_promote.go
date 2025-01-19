package plugin

import (
	"errors"
	"log"
)

func GetScanCommandArgs(args Args) ([][]string, error) {
	var cmdList [][]string

	if args.BuildName == "" || args.BuildNumber == "" {
		return cmdList, errors.New("Valid BuildName and BuildNumber are required")
	}

	tmpServerId := "tmpServeId"
	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs(tmpServerId,
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		log.Println("GetConfigAddConfigCommandArgs error: ", err)
		return cmdList, err
	}

	scanCommandArgs := []string{
		"build-scan", args.BuildName, args.BuildNumber}
	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, scanCommandArgs)
	return cmdList, nil
}

func GetCreateBuildInfoCommandArgs(args Args) ([][]string, error) {
	var cmdList [][]string

	tmpServerId := "tmpServeId"
	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs(tmpServerId,
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		log.Println("GetConfigAddConfigCommandArgs error: ", err)
		return cmdList, err
	}
	buildCollectEnvCommandArgs := []string{"rt", "build-collect-env", args.BuildName, args.BuildNumber}
	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, buildCollectEnvCommandArgs)
	return cmdList, nil
}

func GetBuildInfoPublishCommandArgs(args Args) ([][]string, error) {
	var cmdList [][]string

	tmpServerId := "tmpServeId"
	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs(tmpServerId,
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		log.Println("GetConfigAddConfigCommandArgs error: ", err)
		return cmdList, err
	}
	buildInfoCommandArgs := []string{"rt", "build-publish", args.BuildName, args.BuildNumber}
	err = PopulateArgs(&buildInfoCommandArgs, &args, nil)
	if err != nil {
		return cmdList, err
	}
	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, buildInfoCommandArgs)
	return cmdList, nil
}
