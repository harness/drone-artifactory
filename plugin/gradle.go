package plugin

import (
	"fmt"
	"log"
)

var GradleConfigJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--deploy-ivy-desc=", "PLUGIN_DEPLOY_IVY_DESC", false, false},
	{"--deploy-maven-desc=", "PLUGIN_DEPLOY_MAVEN_DESC", false, false},
	{"--global=", "PLUGIN_GLOBAL", false, false},
	{"--ivy-artifacts-pattern=", "PLUGIN_IVY_ARTIFACTS_PATTERN", false, false},
	{"--ivy-desc-pattern=", "PLUGIN_IVY_DESC_PATTERN", false, false},
	{"--repo-deploy=", "PLUGIN_REPO_DEPLOY", false, false},
	{"--repo-resolve=", "PLUGIN_REPO_RESOLVE", false, false},
	{"--server-id-deploy=", "PLUGIN_SERVER_ID_DEPLOY", false, false},
	{"--server-id-resolve=", "PLUGIN_SERVER_ID_RESOLVE", false, false},
	{"--use-wrapper=", "PLUGIN_USE_WRAPPER", false, false},
	{"--uses-plugin=", "PLUGIN_USES_PLUGIN", false, false},
}

var GradleRunJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false},
	{"--detailed-summary=", "PLUGIN_DETAILED_SUMMARY", false, false},
	{"--format=", "PLUGIN_FORMAT", false, false},
	{"--project=", "PLUGIN_PROJECT", false, false},
	{"--scan=", "PLUGIN_SCAN", false, false},
	{"--threads=", "PLUGIN_THREADS", false, false},
}

func GetGradleCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs(args.ResolverId,
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		return cmdList, err
	}

	gradleConfigCommandArgs := []string{GradleConfig}
	err = PopulateArgs(&gradleConfigCommandArgs, &args, GradleConfigJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}

	gradleTaskCommandArgs := []string{GradleCmd, args.GradleTasks}
	err = PopulateArgs(&gradleTaskCommandArgs, &args, GradleRunJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}

	if len(args.BuildFile) > 0 {
		gradleTaskCommandArgs = append(gradleTaskCommandArgs, "-b "+args.BuildFile)
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, gradleConfigCommandArgs)
	cmdList = append(cmdList, gradleTaskCommandArgs)

	return cmdList, nil
}

var GradleConfigCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--exclude-patterns=", "PLUGIN_EXCLUDE_PATTERNS", false, false},
	{"--global=", "PLUGIN_GLOBAL", false, false},
	{"--include-patterns=", "PLUGIN_INCLUDE_PATTERNS", false, false},
	{"--repo-deploy-releases=", "PLUGIN_DEPLOY_RELEASE_REPO", false, false},
	{"--repo-deploy-snapshots=", "PLUGIN_DEPLOY_SNAPSHOT_REPO", false, false},
	{"--repo-resolve-releases=", "PLUGIN_RESOLVE_RELEASE_REPO", false, false},
	{"--repo-resolve-snapshots=", "PLUGIN_RESOLVE_SNAPSHOT_REPO", false, false},
	{"--repo-deploy=", "PLUGIN_REPO_DEPLOY", false, false},
	{"--repo-resolve=", "PLUGIN_REPO_RESOLVE", false, false},
	{"--server-id-deploy=", "PLUGIN_SERVER_ID_DEPLOY", false, false},
	{"--server-id-resolve=", "PLUGIN_RESOLVER_ID", false, false},
	{"--use-wrapper=", "PLUGIN_USE_WRAPPER", false, false},
}

func GetGradlePublishCommand(args Args) ([][]string, error) {

	var cmdList [][]string
	var jfrogConfigAddConfigCommandArgs []string

	tmpServerId := args.DeployerId
	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs(tmpServerId,
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		log.Println("GetConfigAddConfigCommandArgs error: ", err)
		return cmdList, err
	}

	gradleConfigCommandArgs := []string{GradleConfig}
	err = PopulateArgs(&gradleConfigCommandArgs, &args, GradleConfigCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		log.Println("PopulateArgs error: ", err)
		return cmdList, err
	}
	gradleConfigCommandArgs = append(gradleConfigCommandArgs, "--server-id-deploy="+tmpServerId)
	gradleConfigCommandArgs = append(gradleConfigCommandArgs, "--server-id-resolve="+tmpServerId)

	rtPublishCommandArgs := []string{"gradle", Publish}
	switch {
	case args.Username != "":
		rtPublishCommandArgs = append(rtPublishCommandArgs, "-Pusername="+args.Username)
		rtPublishCommandArgs = append(rtPublishCommandArgs, "-Ppassword="+args.Password)
	case args.AccessToken != "":
		errMsg := "AccessToken is not supported for Gradle" +
			" try username: <username> , password: <access_token> instead"
		log.Println(errMsg)
		return cmdList, fmt.Errorf(errMsg)
	}
	rtPublishCommandArgs = append(rtPublishCommandArgs, "--build-name="+args.BuildName)
	rtPublishCommandArgs = append(rtPublishCommandArgs, "--build-number="+args.BuildNumber)

	rtPublishBuildInfoCommandArgs := []string{"rt", BuildPublish, args.BuildName, args.BuildNumber,
		"--server-id=" + tmpServerId}
	err = PopulateArgs(&rtPublishBuildInfoCommandArgs, &args, RtBuildInfoPublishCmdJsonTagToExeFlagMap)
	if err != nil {
		log.Println("PopulateArgs error: ", err)
		return cmdList, err
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, gradleConfigCommandArgs)
	cmdList = append(cmdList, rtPublishCommandArgs)
	cmdList = append(cmdList, rtPublishBuildInfoCommandArgs)

	return cmdList, nil
}
