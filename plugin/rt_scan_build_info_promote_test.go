package plugin

import (
	"strings"
	"testing"
)

func TestGetScanCommandUserPassword(t *testing.T) {
	args := Args{
		Username:    "ab",
		Password:    "cd",
		Command:     "scan",
		BuildName:   RtBuildName,
		BuildNumber: RtBuildNumber,
		URL:         RtUrlTestStr,
	}
	cmdList, err := GetScanCommandArgs(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"build-scan t2 v1.0",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		if !strings.Contains(cmdStr, wantCmds[i]) {
			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
		}
	}
}

func TestGetBuildInfoCollectCommandUserAccessToken(t *testing.T) {
	args := Args{
		AccessToken: RtAccessToken,
		Command:     "create-build-info",
		BuildName:   RtBuildName,
		BuildNumber: RtBuildNumber,
		URL:         RtUrlTestStr,
	}
	cmdList, err := GetCreateBuildInfoCommandArgs(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"rt build-collect-env t2 v1.0",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		if !strings.Contains(cmdStr, wantCmds[i]) {
			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
		}
	}
}

func TestGetBuildInfoPublishCommandUserPassword(t *testing.T) {
	args := Args{
		Username:    "ab",
		Password:    "cd",
		Command:     "publish-build-info",
		BuildName:   RtBuildName,
		BuildNumber: RtBuildNumber,
		URL:         RtUrlTestStr,
		Module:      RtModule,
		Project:     RtProject,
		SpecPath:    "spec.json",
	}
	cmdList, err := GetBuildInfoPublishCommandArgs(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"config add tmpServerId --url=https://artifactory.test.io/artifactory/ --user $PLUGIN_USERNAME " +
			"--password $PLUGIN_PASSWORD --interactive=false",
		"rt build-publish t2 v1.0",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		if !strings.Contains(cmdStr, wantCmds[i]) {
			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
		}
	}
}

func TestPromoteBuildCommandUserPassword(t *testing.T) {
	args := Args{
		Username:    "ab",
		Password:    "cd",
		Command:     "promote",
		BuildName:   RtBuildName,
		BuildNumber: RtBuildNumber,
		URL:         RtUrlTestStr,
		Target:      "promoted-repo",
		Copy:        "true",
	}
	cmdList, err := GetPromoteCommandArgs(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"rt build-promote --copy=true --url=https://artifactory.test.io/artifactory/ t2 v1.0 promoted-repo " +
			"--user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		if !strings.Contains(cmdStr, wantCmds[i]) {
			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
		}
	}
}
