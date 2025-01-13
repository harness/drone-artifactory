package plugin

import (
	"strings"
	"testing"
)

func TestGetGradleBuildCommandArgs(t *testing.T) {
	tests := []struct {
		args   Args
		output []string
		err    error
	}{
		{
			args: Args{
				BuildTool:   "gradle",
				Username:    "user",
				Password:    "pass",
				URL:         RtUrlTestStr,
				RepoResolve: RtResolveRelRepo,
				RepoDeploy:  RtTestRelRepo,
				GradleTasks: "clean build",
				BuildName:   RtBuildName,
				BuildNumber: RtBuildNumber,
			},
			output: []string{
				"config add tmpSrvConfig --url=" + RtUrlTestStr +
					" --user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
				"gradle-config --repo-deploy=" + RtTestRelRepo + " --repo-resolve=" + RtResolveRelRepo,
				"gradle clean build --build-name=" + RtBuildName + " --build-number=" + RtBuildNumber,
			},
			err: nil,
		},
	}

	for _, tc := range tests {
		result, err := GetGradleCommandArgs(tc.args)
		if err != nil {
			if tc.err == nil {
				t.Errorf("Unexpected error: %v", err)
			} else if err.Error() != tc.err.Error() {
				t.Errorf("Expected error: %v, Got: %v", tc.err, err)
			}
		} else {
			for i, cmd := range result {
				cmdStr := strings.Join(cmd, " ")
				outputStr := tc.output[i]
				if cmdStr != outputStr {
					t.Errorf("Mismatch at index %d. Expected: %s, Got: %s", i, outputStr, cmdStr)
				}
			}
		}
	}
}

func TestGetGradlePublishCommandArgs(t *testing.T) {
	tests := []struct {
		args   Args
		output []string
		err    error
	}{
		{
			args: Args{
				BuildTool:   "gradle",
				Command:     "publish",
				URL:         RtUrlTestStr,
				Username:    "user",
				Password:    "pass",
				RepoResolve: RtResolveRelRepo,
				RepoDeploy:  RtTestRelRepo,
				BuildName:   RtBuildName,
				BuildNumber: RtBuildNumber,
				DeployerId:  RtDeployerId,
			},
			output: []string{
				"config add " + RtDeployerId + " --url=" + RtUrlTestStr +
					" --user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
				"gradle-config --repo-deploy=" + RtTestRelRepo + " --repo-resolve=" +
					RtResolveRelRepo + " --server-id-deploy=" + RtDeployerId + " --server-id-resolve=" + RtDeployerId,
				"gradle publish -Pusername=user -Ppassword=pass --build-name=" +
					RtBuildName + " --build-number=" + RtBuildNumber,
				"rt build-publish " + RtBuildName + " " + RtBuildNumber + " --server-id=" + RtDeployerId,
			},
			err: nil,
		},
	}

	for _, tc := range tests {
		result, err := GetGradlePublishCommand(tc.args)
		if err != nil {
			if tc.err == nil {
				t.Errorf("Unexpected error: %v", err)
			} else if err.Error() != tc.err.Error() {
				t.Errorf("Expected error: %v, Got: %v", tc.err, err)
			}
		} else {
			for i, cmd := range result {
				cmdStr := strings.Join(cmd, " ")
				outputStr := tc.output[i]
				if cmdStr != outputStr {
					t.Errorf("Mismatch at index %d. Expected: %s, Got: %s", i, outputStr, cmdStr)
				}
			}
		}
	}
}
