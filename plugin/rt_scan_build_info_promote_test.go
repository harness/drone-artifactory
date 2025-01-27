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
		"config add tmpServeId --url=https://artifactory.test.io/artifactory/ " +
			"--user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
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
		"config add tmpServeId --url=https://artifactory.test.io/artifactory/ " +
			"--access-token $PLUGIN_ACCESS_TOKEN --interactive=false",
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
		"config add tmpServeId --url=https://artifactory.test.io/artifactory/ " +
			"--user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
		"rt build-publish t2 v1.0",
	}

	for i, cmd := range cmdList {
		cmdStr := strings.Join(cmd, " ")
		if !strings.Contains(cmdStr, wantCmds[i]) {
			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
		}
	}
}

//
//func TestGetDownloadCommandUserAccessToken(t *testing.T) {
//	args := Args{
//		AccessToken: RtAccessToken,
//		Command:     "download",
//		BuildName:   RtBuildName,
//		BuildNumber: RtBuildNumber,
//		URL:         RtUrlTestStr,
//		Module:      RtModule,
//		Project:     RtProject,
//		SpecPath:    "spec.json",
//	}
//	cmdList, err := GetDownloadCommandArgs(args)
//	if err != nil {
//		t.Errorf("Unexpected error: %v", err)
//	}
//
//	wantCmds := []string{
//		"config add tmpServerId --url=https://artifactory.test.io/artifactory/ " +
//			"--access-token $PLUGIN_ACCESS_TOKEN --interactive=false",
//		"rt download   --build-name=t2 --build-number=v1.0 --module=backend_module " +
//			"--project=backend_project --spec=spec.json",
//	}
//
//	for i, cmd := range cmdList {
//		cmdStr := strings.Join(cmd, " ")
//		if !strings.Contains(cmdStr, wantCmds[i]) {
//			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
//		}
//	}
//}
//
//func TestGetCleanupCommandUserPassword(t *testing.T) {
//	args := Args{
//		Username:    "ab",
//		Password:    "cd",
//		Command:     "cleanup",
//		BuildName:   RtBuildName,
//		BuildNumber: RtBuildNumber,
//		URL:         RtUrlTestStr,
//	}
//	cmdList, err := GetCleanupCommandArgs(args)
//	if err != nil {
//		t.Errorf("Unexpected error: %v", err)
//	}
//
//	wantCmds := []string{
//		"config add tmpServerId --url=https://artifactory.test.io/artifactory/ " +
//			"--user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
//		"rt build-clean t2 v1.0",
//	}
//
//	for i, cmd := range cmdList {
//		cmdStr := strings.Join(cmd, " ")
//		if !strings.Contains(cmdStr, wantCmds[i]) {
//			t.Errorf("Expected: |%s|, Got: |%s|", wantCmds[i], cmdStr)
//		}
//	}
//}
