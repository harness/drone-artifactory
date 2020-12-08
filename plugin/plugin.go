// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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

	cmdArgs := []string{"rt", "u", fmt.Sprintf("--url=%s", args.URL)}
	if args.Retries != 0 {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--retries=%d", args.Retries))
	}

	// Set authentication params
	if args.Username != "" && args.Password != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--user=%s", args.Username))
		cmdArgs = append(cmdArgs, fmt.Sprintf("--password=%s", args.Password))
	} else if args.APIKey != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--apikey=%s", args.APIKey))
	} else if args.AccessToken != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--access-token=%s", args.AccessToken))
	} else {
		return fmt.Errorf("either username/password, api key or access token needs to be set")
	}

	if args.Source == "" {
		return fmt.Errorf("source file needs to be set")
	}
	if args.Target == "" {
		return fmt.Errorf("target path needs to be set")
	}
	cmdArgs = append(cmdArgs, args.Source, args.Target)

	cmd := exec.Command("jfrog", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
