package plugin

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetMavenBuildCommandUserPassword(t *testing.T) {
	args := Args{
		Username:            "ab",
		Password:            "cd",
		BuildTool:           RtMvnBuildTool,
		MvnPomFile:          "pom.xml",
		MvnGoals:            "clean install",
		BuildName:           RtBuildName,
		BuildNumber:         RtBuildNumber,
		URL:                 RtUrlTestStr,
		ResolverId:          RtRslvId,
		ResolveReleaseRepo:  RtResolveRelRepo,
		ResolveSnapshotRepo: RtResolveSnapshotRepo,
		DeployerId:          RtDeployerId,
	}
	cmdList, err := GetMavenBuildCommandArgs(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"config add resolve_gen_maven_01 --url=https://artifactory.test.io/artifactory/ " +
			"--user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
		"mvn-config --repo-resolve-releases=mvn_repo_resolve_releases_01 " +
			"--repo-resolve-snapshots=mvn_repo_resolve_snapshots_01 --server-id-resolve=resolve_gen_maven_01",
		"mvn clean install --build-name=t2 --build-number=v1.0 -f pom.xml",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		fmt.Println(cmdStr)

		ret := strings.Compare(cmdStr, wantCmds[i])
		if ret != 0 {
			t.Errorf("Expected: %s, Got: %s", wantCmds[i], cmdStr)
		}
	}

}

func TestGetMavenBuildCommandAccessToken(t *testing.T) {
	args := Args{
		AccessToken:         RtAccessToken,
		BuildTool:           RtMvnBuildTool,
		MvnPomFile:          "pom.xml",
		MvnGoals:            "clean install",
		BuildName:           RtBuildName,
		BuildNumber:         RtBuildNumber,
		URL:                 RtUrlTestStr,
		ResolverId:          RtRslvId,
		ResolveReleaseRepo:  RtResolveRelRepo,
		ResolveSnapshotRepo: RtResolveSnapshotRepo,
		DeployerId:          RtDeployerId,
	}
	cmdList, err := GetMavenBuildCommandArgs(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"config add resolve_gen_maven_01 --url=https://artifactory.test.io/artifactory/ " +
			"--access-token $PLUGIN_ACCESS_TOKEN --interactive=false",
		"mvn-config --repo-resolve-releases=mvn_repo_resolve_releases_01 " +
			"--repo-resolve-snapshots=mvn_repo_resolve_snapshots_01 --server-id-resolve=resolve_gen_maven_01",
		"mvn clean install --build-name=t2 --build-number=v1.0 -f pom.xml",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		fmt.Println(cmdStr)
		ret := strings.Compare(cmdStr, wantCmds[i])
		if ret != 0 {
			t.Errorf("Expected: %s, Got: %s", wantCmds[i], cmdStr)
		}
	}
}

func TestGetMavenPublishCommandUserNamePassword(t *testing.T) {
	args := Args{
		Username:           "ab",
		Password:           "cd",
		Command:            "publish",
		BuildTool:          RtMvnBuildTool,
		BuildName:          RtBuildName,
		BuildNumber:        RtBuildNumber,
		URL:                RtUrlTestStr,
		DeployerId:         RtDeployerId,
		DeployReleaseRepo:  RtTestRelRepo,
		DeploySnapshotRepo: RtTestSnapshotRepo,
	}
	cmdList, err := GetMavenPublishCommand(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"config add deploy_gen_maven_01 --url=https://artifactory.test.io/artifactory/ --user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
		"mvn-config --repo-deploy-releases=mvn_repo_deploy_releases_01 --repo-deploy-snapshots=mvn_repo_deploy_snapshots_01",
		"mvn deploy --build-name=t2 --build-number=v1.0",
		"rt build-publish t2 v1.0 --server-id=deploy_gen_maven_01",
	}
	_ = wantCmds

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		ret := strings.Compare(cmdStr, wantCmds[i])
		if ret != 0 {
			t.Errorf("Expected: %s, Got: %s", wantCmds[i], cmdStr)
		}
	}
}

func TestGetMavenPublishCommandAccessToken(t *testing.T) {
	args := Args{
		AccessToken:        RtAccessToken,
		Command:            "publish",
		BuildTool:          RtMvnBuildTool,
		BuildName:          RtBuildName,
		BuildNumber:        RtBuildNumber,
		URL:                RtUrlTestStr,
		DeployerId:         RtDeployerId,
		DeployReleaseRepo:  RtTestRelRepo,
		DeploySnapshotRepo: RtTestSnapshotRepo,
	}
	cmdList, err := GetMavenPublishCommand(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"config add deploy_gen_maven_01 --url=https://artifactory.test.io/artifactory/ --access-token $PLUGIN_ACCESS_TOKEN --interactive=false",
		"mvn-config --repo-deploy-releases=mvn_repo_deploy_releases_01 --repo-deploy-snapshots=mvn_repo_deploy_snapshots_01",
		"mvn deploy --build-name=t2 --build-number=v1.0",
		"rt build-publish t2 v1.0 --server-id=deploy_gen_maven_01",
	}
	_ = wantCmds

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		ret := strings.Compare(cmdStr, wantCmds[i])
		if ret != 0 {
			t.Errorf("Expected: %s, Got: %s", wantCmds[i], cmdStr)
		}
	}
}
