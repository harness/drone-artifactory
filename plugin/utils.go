package plugin

import (
	"context"
	"log"
)

const (
	MvnCmd       = "mvn"
	MvnConfig    = "mvn-config"
	BuildPublish = "build-publish"
)

type RtCommand struct {
	Command   string `envconfig:"PLUGIN_COMMAND"`
	BuildTool string `envconfig:"PLUGIN_BUILD_TOOL"`
	RtMvnCommand
}

type RtMvnCommand struct {
	ResolveReleaseRepo  string `envconfig:"PLUGIN_RESOLVE_RELEASE_REPO"`
	ResolveSnapshotRepo string `envconfig:"PLUGIN_RESOLVE_SNAPSHOT_REPO"`
	DeployReleaseRepo   string `envconfig:"PLUGIN_DEPLOY_RELEASE_REPO"`
	DeploySnapshotRepo  string `envconfig:"PLUGIN_DEPLOY_SNAPSHOT_REPO"`
	DeployRepo          string `envconfig:"PLUGIN_DEPLOY_REPO"`

	MvnGoals     string `envconfig:"PLUGIN_GOALS"`
	MvnPomFile   string `envconfig:"PLUGIN_POM_FILE"`
	ProjectKey   string `envconfig:"PLUGIN_PROJECT_KEY"`
	OptionalArgs string `envconfig:"PLUGIN_OPTIONAL_ARGS"`

	DeployerId string `envconfig:"PLUGIN_DEPLOYER_ID"`
	ResolverId string `envconfig:"PLUGIN_RESOLVER_ID"`
}

func HandleRtCommands(ctx context.Context, args Args) error {
	commandsList, err := GetRtCommandsList(ctx, args)
	for _, cmd := range commandsList {
		execArgs := []string{getJfrogBin()}
		execArgs = append(execArgs, cmd...)
		err := ExecCommand(args, execArgs)
		if err != nil {
			log.Println("Error Unable to run err = ", err)
			return err
		}
	}

	return err
}
