package plugin

import "fmt"

func GetMavenCommandArgs(userName, password, url,
	repoResolveReleases, repoResolveSnapshots, repoDeployReleases, repoDeploySnapshots,
	pomFile, goals, buildName, buildNumber string, numThreads int,
	insecureTls string, projectKey string,
	otherOpts string) ([][]string, error) {

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

	mvnGoalCommandArgs := []string{"mvn", goals}

	if len(pomFile) > 0 {
		mvnGoalCommandArgs = append(mvnGoalCommandArgs, "-file="+pomFile)
	}
	if len(buildName) > 0 {
		mvnGoalCommandArgs = append(mvnGoalCommandArgs, "--build-name="+buildName)
	}
	if len(buildNumber) > 0 {
		mvnGoalCommandArgs = append(mvnGoalCommandArgs, "--build-number="+buildNumber)
	}
	if numThreads > 0 {
		mvnGoalCommandArgs = append(mvnGoalCommandArgs, fmt.Sprintf("--threads=%d", numThreads))
	}
	if len(insecureTls) > 0 {
		mvnGoalCommandArgs = append(mvnGoalCommandArgs, "--insecure-tls="+insecureTls)
	}
	if len(projectKey) > 0 {
		mvnGoalCommandArgs = append(mvnGoalCommandArgs, "--project="+projectKey)
	}
	mvnGoalCommandArgs = append(mvnGoalCommandArgs, otherOpts)

	mvnCmdList = append(mvnCmdList, jfrogConfigAddConfigCommandArgs)
	mvnCmdList = append(mvnCmdList, mvnConfigCommandArgs)
	mvnCmdList = append(mvnCmdList, mvnGoalCommandArgs)

	return mvnCmdList, nil

}

const (
	MvnCmd = "mvn"
)
