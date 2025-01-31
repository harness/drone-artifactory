package plugin

import (
	"strings"
	"testing"
)

const (
	RtModule  = "backend_module"
	RtProject = "backend_project"
)

func TestGetDownloadCommandUserPassword(t *testing.T) {
	args := Args{
		Username:    "ab",
		Password:    "cd",
		Command:     "download",
		BuildName:   RtBuildName,
		BuildNumber: RtBuildNumber,
		URL:         RtUrlTestStr,
		Module:      RtModule,
		Project:     RtProject,
		SpecPath:    "spec.json",
	}
	cmdList, err := GetDownloadCommandArgs(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"rt download --user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD   " + "--build-name=t2 --build-number=v1.0 " +
			"--module=backend_module --project=backend_project --url=https://artifactory.test.io/artifactory/ --spec=spec.json",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		if !strings.Contains(cmdStr, wantCmds[i]) {
			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
		}
	}
}

func TestGetDownloadCommandUserAccessToken(t *testing.T) {
	args := Args{
		AccessToken: RtAccessToken,
		Command:     "download",
		BuildName:   RtBuildName,
		BuildNumber: RtBuildNumber,
		URL:         RtUrlTestStr,
		Module:      RtModule,
		Project:     RtProject,
		SpecPath:    "spec.json",
	}
	cmdList, err := GetDownloadCommandArgs(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"rt download --access-token $PLUGIN_ACCESS_TOKEN   --build-name=t2 --build-number=v1.0 --module=backend_module" +
			" --project=backend_project --url=https://artifactory.test.io/artifactory/ --spec=spec.json",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		if !strings.Contains(cmdStr, wantCmds[i]) {
			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
		}
	}
}

func TestGetCleanupCommandUserPassword(t *testing.T) {
	args := Args{
		Username:    "ab",
		Password:    "cd",
		Command:     "cleanup",
		BuildName:   RtBuildName,
		BuildNumber: RtBuildNumber,
		URL:         RtUrlTestStr,
	}
	cmdList, err := GetCleanupCommandArgs(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"rt build-clean t2 v1.0",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		if !strings.Contains(cmdStr, wantCmds[i]) {
			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
		}
	}
}
