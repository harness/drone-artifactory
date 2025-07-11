package plugin

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

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

func GetMavenBuildCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs(args.ResolverId,
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		return cmdList, err
	}

	mvnConfigCommandArgs := []string{MvnConfig}
	// Add necessary parameters for Windows to prevent all interactive prompts
	if runtime.GOOS == "windows" {
		// These parameters prevent all interactive prompts
		mvnConfigCommandArgs = append(mvnConfigCommandArgs, "--global=true")
		// Add server ID for deployment/resolution
		if args.ResolverId != "" {
			mvnConfigCommandArgs = append(mvnConfigCommandArgs, "--server-id-resolve="+args.ResolverId)
			mvnConfigCommandArgs = append(mvnConfigCommandArgs, "--server-id-deploy="+args.ResolverId)
		}
		// Add repos to prevent prompts
		// Must set both release and snapshot repos to prevent errors
		if args.ResolveReleaseRepo == "" {
			mvnConfigCommandArgs = append(mvnConfigCommandArgs, "--repo-resolve-releases=libs-release")
		}
		if args.ResolveSnapshotRepo == "" {
			mvnConfigCommandArgs = append(mvnConfigCommandArgs, "--repo-resolve-snapshots=libs-snapshot")
		}
		if args.DeployReleaseRepo == "" {
			mvnConfigCommandArgs = append(mvnConfigCommandArgs, "--repo-deploy-releases=libs-release-local")
		}
		if args.DeploySnapshotRepo == "" {
			mvnConfigCommandArgs = append(mvnConfigCommandArgs, "--repo-deploy-snapshots=libs-snapshot-local")
		}
	}

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

var RtBuildInfoPublishCmdJsonTagToExeFlagMap = []JsonTagToExeFlagMapStringItem{
	{"--project=", "PLUGIN_PROJECT", false, false},
}

func GetMavenPublishCommand(args Args) ([][]string, error) {
	var cmdList [][]string
	var jfrogConfigAddConfigCommandArgs []string

	tmpServerId := args.DeployerId
	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs(tmpServerId,
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		logrus.Println("GetConfigAddConfigCommandArgs error: ", err)
		return cmdList, err
	}

	mvnConfigCommandArgs := []string{MvnConfig}
	err = PopulateArgs(&mvnConfigCommandArgs, &args, MavenConfigCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		logrus.Println("mvnConfigCommandArgs PopulateArgs error: ", err)
		return cmdList, err
	}

	rtPublishCommandArgs := []string{MvnCmd, Deploy,
		"--build-name=" + args.BuildName, "--build-number=" + args.BuildNumber}
	err = PopulateArgs(&rtPublishCommandArgs, &args, RtBuildInfoPublishCmdJsonTagToExeFlagMap)
	if err != nil {
		logrus.Println("rtPublishCommandArgs PopulateArgs error: ", err)
		return cmdList, err
	}

	rtPublishBuildInfoCommandArgs := []string{"rt", BuildPublish, args.BuildName, args.BuildNumber,
		"--server-id=" + tmpServerId}
	err = PopulateArgs(&rtPublishBuildInfoCommandArgs, &args, RtBuildInfoPublishCmdJsonTagToExeFlagMap)
	if err != nil {
		logrus.Println("PopulateArgs error: ", err)
		return cmdList, err
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, mvnConfigCommandArgs)
	cmdList = append(cmdList, rtPublishCommandArgs)
	cmdList = append(cmdList, rtPublishBuildInfoCommandArgs)

	if IsBuildDiscardArgs(args) {
		buildDiscardBuildArgsList, err := GetBuildDiscardCommandArgs(args)
		if err != nil {
			logrus.Println("GetBuildDiscardCommandArgs error: ", err)
			return cmdList, err
		}
		cmdList = append(cmdList, buildDiscardBuildArgsList...)
	}
	return cmdList, nil
}
