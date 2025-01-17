package plugin

import (
	"strings"
	"testing"
)

const (
	RtModule  = "backend_module"
	RtProject = "backend_project"
)

func TestGetUploadCommandUserPassword(t *testing.T) {
	args := Args{
		Username:    "ab",
		Password:    "cd",
		Command:     "upload",
		BuildName:   RtBuildName,
		BuildNumber: RtBuildNumber,
		URL:         RtUrlTestStr,
		Module:      RtModule,
		Project:     RtProject,
		SpecPath:    "spec.json",
	}
	cmdList, err := GetUploadCommandArgs(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"config add tmpServerId --url=https://artifactory.test.io/artifactory/ " +
			"--user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
		"rt upload   --build-name=t2 --build-number=v1.0 --module=backend_module " +
			"--project=backend_project --spec=",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		if !strings.Contains(cmdStr, wantCmds[i]) {
			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
		}
	}
}

func TestGetUploadCommandUserAccessToken(t *testing.T) {
	args := Args{
		AccessToken: RtAccessToken,
		Command:     "upload",
		BuildName:   RtBuildName,
		BuildNumber: RtBuildNumber,
		URL:         RtUrlTestStr,
		Module:      RtModule,
		Project:     RtProject,
		SpecPath:    "spec.json",
	}
	cmdList, err := GetUploadCommandArgs(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	wantCmds := []string{
		"config add tmpServerId --url=https://artifactory.test.io/artifactory/ " +
			"--access-token $PLUGIN_ACCESS_TOKEN --interactive=false",
		"rt upload   --build-name=t2 --build-number=v1.0 --module=backend_module --project=backend_project --spec=",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		if !strings.Contains(cmdStr, wantCmds[i]) {
			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
		}
	}
}

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
		"config add tmpServerId --url=https://artifactory.test.io/artifactory/ " +
			"--user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
		"rt download   --build-name=t2 --build-number=v1.0 --module=backend_module " +
			"--project=backend_project --spec=spec.json",
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
		"config add tmpServerId --url=https://artifactory.test.io/artifactory/ " +
			"--access-token $PLUGIN_ACCESS_TOKEN --interactive=false",
		"rt download   --build-name=t2 --build-number=v1.0 --module=backend_module " +
			"--project=backend_project --spec=spec.json",
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
		"config add tmpServerId --url=https://artifactory.test.io/artifactory/ " +
			"--user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
		"rt build-clean t2 v1.0",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		if !strings.Contains(cmdStr, wantCmds[i]) {
			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
		}
	}
}
