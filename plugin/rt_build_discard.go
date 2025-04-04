package plugin

import (
	"github.com/sirupsen/logrus"
)

var BuildDiscardCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--async=", "PLUGIN_ASYNC", false, false},
	{"--delete-artifacts=", "PLUGIN_DELETE_ARTIFACTS", false, false},
	{"--exclude-builds=", "PLUGIN_EXCLUDE_BUILDS", false, false},
	{"--max-builds=", "PLUGIN_MAX_BUILDS", false, false},
	{"--max-days=", "PLUGIN_MAX_DAYS", false, false},
}

func GetBuildDiscardCommandArgs(args Args) ([][]string, error) {
	bdiServerId := tmpServerId + "bdi"
	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs(bdiServerId,
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		logrus.Println("GetConfigAddConfigCommandArgs error: ", err)
		return cmdList, err
	}

	buildDiscardCmd, err := GetBuildDiscardCommand(args)

	if err != nil {
		logrus.Println("Error in GetBuildDiscardCommand ", err)
		return cmdList, err
	}
	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, buildDiscardCmd)
	return cmdList, nil
}

func GetBuildDiscardCommand(args Args) ([]string, error) {
	buildDiscardCommandArgs := []string{"rt", "build-discard"}
	err := PopulateArgs(&buildDiscardCommandArgs, &args, BuildDiscardCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		logrus.Println("Error in PopulateArgs ", err)
		return buildDiscardCommandArgs, err
	}
	buildDiscardCommandArgs = append(buildDiscardCommandArgs, args.BuildName)
	return buildDiscardCommandArgs, nil
}
