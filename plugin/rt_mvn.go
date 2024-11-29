package plugin

import "fmt"

func GetMavenCommandArgs(userName, password, url,
	repoResolveReleases, repoResolveSnapshots, repoDeployReleases, repoDeploySnapshots,
	pomFile, goals, otherOpts string) ([][]string, error) {

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

	srvConfigStr := "tmpSrvConfig"

	jfrogConfigAddConfigCommandArgs := []string{
		"config", "add", srvConfigStr, "--url=" + url,
		"--user=" + userName, "--password=" + password, "--interactive=false"}

	mvnConfigCommandArgs := []string{"mvn-config", "--global", "--repo-resolve-releases=" + repoResolveReleases,
		"--repo-resolve-snapshots=" + repoResolveSnapshots,
		"--repo-deploy-releases=" + repoDeployReleases, "--repo-deploy-snapshots=" + repoDeploySnapshots}

	mvnGoalCommandArgs := []string{goals, otherOpts}

	mvnCmdList = append(mvnCmdList, jfrogConfigAddConfigCommandArgs)
	mvnCmdList = append(mvnCmdList, mvnConfigCommandArgs)
	mvnCmdList = append(mvnCmdList, mvnGoalCommandArgs)

	return mvnCmdList, nil

}

const (
	MvnCmd = "mvn"
)
