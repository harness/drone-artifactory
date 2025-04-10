package plugin

import (
	"strings"
	"testing"
)

func TestBuildDiscard(t *testing.T) {
	args := Args{
		Command:         "build-discard",
		Username:        "ab",
		Password:        "cd",
		URL:             "https://artifactory.test.io/artifactory/",
		BuildName:       "t2",
		BuildNumber:     "v1.0",
		DeleteArtifacts: "true",
		MaxBuilds:       "5",
		MaxDays:         "7",
	}

	cmdList, err := GetBuildDiscardCommandArgs(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	gotCmds := []string{}
	for _, cmd := range cmdList {
		gotCmds = append(gotCmds, strings.Join(cmd, " "))
	}

	wantCmds := []string{
		"config add tmpServerIdbdi --url=https://artifactory.test.io/artifactory/ --user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
		"rt build-discard --delete-artifacts=true --max-builds=5 --max-days=7 t2",
	}

	if len(gotCmds) != len(wantCmds) {
		t.Fatalf("Expected %d commands, but got %d", len(wantCmds), len(gotCmds))
	}

	for i := range wantCmds {
		if gotCmds[i] != wantCmds[i] {
			t.Errorf("Command mismatch at index %d:\nExpected: %q\nGot:      %q", i, wantCmds[i], gotCmds[i])
		}
	}
}

func TestGradleBuildDiscard(t *testing.T) {
	args := Args{
		BuildTool:       "gradle",
		Command:         "publish",
		Username:        "ab0",
		Password:        "cd",
		URL:             "https://artifactory.test.io/artifactory/",
		BuildName:       "t2",
		BuildNumber:     "v1.0",
		DeleteArtifacts: "true",
		MaxBuilds:       "5",
		MaxDays:         "7",
	}

	cmdList, err := GetGradlePublishCommand(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	gotCmds := []string{}
	for _, cmd := range cmdList {
		gotCmds = append(gotCmds, strings.Join(cmd, " "))
	}

	wantCmds := []string{
		"config add tmpServerId --url=https://artifactory.test.io/artifactory/ --user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
		"gradle-config --server-id-deploy= --server-id-resolve=",
		"gradle publish -Pusername=ab0 -Ppassword=cd --build-name=t2 --build-number=v1.0",
		"rt build-publish t2 v1.0 --server-id=",
		"config add tmpServerIdbdi --url=https://artifactory.test.io/artifactory/ --user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
		"rt build-discard --delete-artifacts=true --max-builds=5 --max-days=7 t2",
	}

	if len(gotCmds) != len(wantCmds) {
		t.Fatalf("Expected %d commands, but got %d", len(wantCmds), len(gotCmds))
	}

	for i := range wantCmds {
		if gotCmds[i] != wantCmds[i] {
			t.Errorf("Command mismatch at index %d:\nExpected: %q\nGot:      %q", i, wantCmds[i], gotCmds[i])
		}
	}
}

func TestMvnBuildDiscard(t *testing.T) {
	args := Args{
		BuildTool:       "mvn",
		Command:         "publish",
		Username:        "ab0",
		Password:        "cd",
		URL:             "https://artifactory.test.io/artifactory/",
		BuildName:       "t2",
		BuildNumber:     "v1.0",
		DeleteArtifacts: "true",
		MaxBuilds:       "5",
		MaxDays:         "7",
	}

	cmdList, err := GetMavenPublishCommand(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	gotCmds := []string{}
	for _, cmd := range cmdList {
		gotCmds = append(gotCmds, strings.Join(cmd, " "))
	}

	wantCmds := []string{
		"config add tmpServerId --url=https://artifactory.test.io/artifactory/ --user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
		"mvn-config",
		"mvn deploy --build-name=t2 --build-number=v1.0",
		"rt build-publish t2 v1.0 --server-id=",
		"config add tmpServerIdbdi --url=https://artifactory.test.io/artifactory/ --user $PLUGIN_USERNAME --password $PLUGIN_PASSWORD --interactive=false",
		"rt build-discard --delete-artifacts=true --max-builds=5 --max-days=7 t2",
	}

	if len(gotCmds) != len(wantCmds) {
		t.Fatalf("Expected %d commands, but got %d", len(wantCmds), len(gotCmds))
	}

	for i := range wantCmds {
		if gotCmds[i] != wantCmds[i] {
			t.Errorf("Command mismatch at index %d:\nExpected: %q\nGot:      %q", i, wantCmds[i], gotCmds[i])
		}
	}
}
