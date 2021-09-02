// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	// TODO replace or remove
	Username    string `envconfig:"PLUGIN_USERNAME"`
	Password    string `envconfig:"PLUGIN_PASSWORD"`
	APIKey      string `envconfig:"PLUGIN_API_KEY"`
	AccessToken string `envconfig:"PLUGIN_ACCESS_TOKEN"`
	URL         string `envconfig:"PLUGIN_URL"`
	Source      string `envconfig:"PLUGIN_SOURCE"`
	Target      string `envconfig:"PLUGIN_TARGET"`
	Retries     int    `envconfig:"PLUGIN_RETRIES"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {
	// write code here
	if args.URL == "" {
		return fmt.Errorf("url needs to be set")
	}

	cmdArgs := []string{"jfrog", "rt", "u", fmt.Sprintf("--url %s", args.URL)}
	if args.Retries != 0 {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--retries=%d", args.Retries))
	}

	// Set authentication params
	if args.Username != "" && args.Password != "" {
		cmdArgs = append(cmdArgs, "--user $PLUGIN_USERNAME")
		cmdArgs = append(cmdArgs, "--password $PLUGIN_PASSWORD")
	} else if args.APIKey != "" {
		cmdArgs = append(cmdArgs, "--apikey $PLUGIN_API_KEY")
	} else if args.AccessToken != "" {
		cmdArgs = append(cmdArgs, "--access-token $PLUGIN_ACCESS_TOKEN")
	} else {
		return fmt.Errorf("either username/password, api key or access token needs to be set")
	}

	if args.Source == "" {
		return fmt.Errorf("source file needs to be set")
	}
	if args.Target == "" {
		return fmt.Errorf("target path needs to be set")
	}
	cmdArgs = append(cmdArgs, fmt.Sprintf("\"%s\"", args.Source), args.Target)
	cmdStr := strings.Join(cmdArgs[:], " ")

	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "JFROG_CLI_OFFER_CONFIG=false")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)

	err := cmd.Run()
	return err
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
