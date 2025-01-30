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

	authParams, err := setAuthParams([]string{}, Args{Username: args.Username,
		Password: args.Password, AccessToken: args.AccessToken, APIKey: args.APIKey})
	if err != nil {
		return cmdList, err
	}

	scanCommandArgs := []string{
		"build-scan", args.BuildName, args.BuildNumber}
	scanCommandArgs = append(scanCommandArgs, "--url="+args.URL)
	scanCommandArgs = append(scanCommandArgs, authParams...)
	cmdList = append(cmdList, scanCommandArgs)

	return cmdList, nil
}

func GetCreateBuildInfoCommandArgs(args Args) ([][]string, error) {
	var cmdList [][]string

	buildCollectEnvCommandArgs := []string{"rt", "build-collect-env", args.BuildName, args.BuildNumber}
	cmdList = append(cmdList, buildCollectEnvCommandArgs)
	return cmdList, nil
}

func GetBuildInfoPublishCommandArgs(args Args) ([][]string, error) {
	var cmdList [][]string

	tmpServerId := tmpServerId
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

func GetPromoteCommandArgs(args Args) ([][]string, error) {
	var cmdList [][]string

	promoteCommandArgs := []string{"rt", "build-promote"}
	if args.Copy != "" {
		promoteCommandArgs = append(promoteCommandArgs, "--copy="+args.Copy)
	}
	promoteCommandArgs = append(promoteCommandArgs, "--url="+args.URL)
	promoteCommandArgs = append(promoteCommandArgs, args.BuildName, args.BuildNumber, args.Target)
	authParams, err := setAuthParams([]string{}, Args{Username: args.Username, Password: args.Password, AccessToken: args.AccessToken, APIKey: args.APIKey})
	if err != nil {
		return cmdList, err
	}
	promoteCommandArgs = append(promoteCommandArgs, authParams...)
	cmdList = append(cmdList, promoteCommandArgs)
	return cmdList, nil
}
