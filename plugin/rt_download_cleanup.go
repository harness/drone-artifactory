package plugin

import (
	"fmt"
	"os"
	"time"
)

var DownloadCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false},
	{"--module=", "PLUGIN_MODULE", false, false},
	{"--project=", "PLUGIN_PROJECT", false, false},
	{"--url=", "PLUGIN_URL", false, false},
	{"--spec=", "PLUGIN_SPEC", false, false},
	{"--spec=", "PLUGIN_SPEC_PATH", false, false},
}

func GetDownloadCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string
	downloadCommandArgs := []string{"rt", "download"}

	authParams, err := setAuthParams([]string{}, Args{Username: args.Username,
		Password: args.Password, AccessToken: args.AccessToken, APIKey: args.APIKey})
	if err != nil {
		return cmdList, err
	}

	if args.Spec != "" {
		fileName := getTimestampedFileName()
		err = writeToFile(fileName, args.Spec)
		if err != nil {
			return cmdList, err
		}
		args.Spec = ""
		args.SpecPath = fileName
	}

	downloadCommandArgs = append(downloadCommandArgs, authParams...)
	downloadCommandArgs = append(downloadCommandArgs, args.Target, args.Source)
	downloadCommandArgs = append(downloadCommandArgs)

	err = PopulateArgs(&downloadCommandArgs, &args, DownloadCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}

	cmdList = append(cmdList, downloadCommandArgs)
	return cmdList, nil
}

func GetCleanupCommandArgs(args Args) ([][]string, error) {
	var cmdList [][]string
	cleanupCommandArgs := []string{"rt", "build-clean", args.BuildName, args.BuildNumber}
	cmdList = append(cmdList, cleanupCommandArgs)
	return cmdList, nil
}

func writeToFile(filePath, content string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

func getTimestampedFileName() string {
	timestamp := time.Now().Format("20060102_150405.000")
	fileName := fmt.Sprintf("%s_spec.json", timestamp)
	return fileName
}
