// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"fmt"
	"strings"
	"testing"
)

const (
	RtUrlTestStr          = "https://artifactory.test.io/artifactory/"
	RtAccessToken         = "qXsj28"
	RtMvnBuildTool        = "mvn"
	RtBuildName           = "t2"
	RtBuildNumber         = "v1.0"
	RtDeployerId          = "deploy_gen_maven_01"
	RtTestRelRepo         = "mvn_repo_deploy_releases_01"
	RtTestSnapshotRepo    = "mvn_repo_deploy_snapshots_01"
	RtRslvId              = "resolve_gen_maven_01"
	RtResolveRelRepo      = "mvn_repo_resolve_releases_01"
	RtResolveSnapshotRepo = "mvn_repo_resolve_snapshots_01"
)

func TestSetAuthParams(t *testing.T) {
	tests := []struct {
		cmdArgs []string
		args    Args
		output  []string
		err     error
	}{
		// Test case 1
		{
			cmdArgs: []string{"executable", "arg1", "arg2"},
			args:    Args{Username: "john", Password: "password123", APIKey: "", AccessToken: ""},
			output:  []string{"executable", "arg1", "arg2", "--user $PLUGIN_USERNAME", "--password $PLUGIN_PASSWORD"},
			err:     nil,
		},
		// Test case 2
		{
			cmdArgs: []string{"./app", "--flag"},
			args:    Args{Username: "", Password: "", APIKey: "secretkey", AccessToken: ""},
			output:  []string{"./app", "--flag", "--apikey $PLUGIN_API_KEY"},
			err:     nil,
		},
		// Test case 3
		{
			cmdArgs: []string{"script.sh", "-option"},
			args:    Args{Username: "", Password: "", APIKey: "", AccessToken: "token123"},
			output:  []string{"script.sh", "-option", "--access-token $PLUGIN_ACCESS_TOKEN"},
			err:     nil,
		},
		// Test case 4
		{
			cmdArgs: []string{"command", "arg1"},
			args:    Args{Username: "", Password: "", APIKey: "", AccessToken: ""},
			output:  nil,
			err:     fmt.Errorf("either username/password, api key or access token needs to be set"),
		},
		// Test case 5
		{
			cmdArgs: []string{"app", "-flag"},
			args:    Args{Username: "user", Password: "", APIKey: "apikey123", AccessToken: ""},
			output:  []string{"app", "-flag", "--apikey $PLUGIN_API_KEY"},
			err:     nil,
		},
	}

	for _, tc := range tests {
		result, err := setAuthParams(tc.cmdArgs, tc.args)
		if err != nil {
			if tc.err == nil {
				t.Errorf("Unexpected error: %v", err)
			} else if err.Error() != tc.err.Error() {
				t.Errorf("Expected error: %v, Got: %v", tc.err, err)
			}
		} else {
			if len(result) != len(tc.output) {
				t.Errorf("Expected output length: %d, Got: %d", len(tc.output), len(result))
			}
			for j := range result {
				if result[j] != tc.output[j] {
					t.Errorf("Mismatch at index %d. Expected: %s, Got: %s", j, tc.output[j], result[j])
				}
			}
		}
	}
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		inputURL string
		expected string
		err      error
	}{
		{
			inputURL: "https://artifactory.maryout.com/artifactory/test44",
			expected: "https://artifactory.maryout.com/artifactory/",
			err:      nil,
		},
		{
			inputURL: "https://artifactory.maryout.com/artifactory/test/newdir/",
			expected: "https://artifactory.maryout.com/artifactory/",
			err:      nil,
		},
		{
			inputURL: "https://opautomates.jfrog.io/artifactory/test55/",
			expected: "https://opautomates.jfrog.io/artifactory/",
			err:      nil,
		},
		{
			inputURL: "https://opautomates.jfrog.io/artifactory",
			expected: "https://opautomates.jfrog.io/artifactory/",
			err:      nil,
		},
		{
			inputURL: "https://example.com/notartifactory",
			expected: "",
			err:      fmt.Errorf("url does not contain '/artifactory': https://example.com/notartifactory"),
		},
		{
			inputURL: "invalid-url",
			expected: "",
			err:      fmt.Errorf("invalid URL: invalid-url"),
		},
	}

	for _, tc := range tests {
		result, err := sanitizeURL(tc.inputURL)
		if err != nil {
			if tc.err == nil {
				t.Errorf("Unexpected error: %v", err)
			} else if err.Error() != tc.err.Error() {
				t.Errorf("Expected error: %v, Got: %v", tc.err, err)
			}
		} else {
			if result != tc.expected {
				t.Errorf("For URL %s, Expected: %s, Got: %s", tc.inputURL, tc.expected, result)
			}
		}
	}
}

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
