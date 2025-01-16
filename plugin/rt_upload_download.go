package plugin

import (
	"fmt"
	"os"
	"time"
)

var UploadCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false},
	{"--module=", "PLUGIN_MODULE", false, false},
	{"--project=", "PLUGIN_PROJECT", false, false},
	{"--server-id=", "PLUGIN_SERVER_ID", false, false},
	{"--spec=", "PLUGIN_SPEC", false, false},
	{"--spec=", "PLUGIN_SPEC_PATH", false, false},
}

func GetUploadCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string
	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs("tmpServerId", args.Username, args.Password,
		args.URL, args.AccessToken, args.APIKey)
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

	uploadCommandArgs := []string{"rt", "upload", args.Source, args.Target}
	err = PopulateArgs(&uploadCommandArgs, &args, UploadCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, uploadCommandArgs)

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

var DownloadCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false},
	{"--module=", "PLUGIN_MODULE", false, false},
	{"--project=", "PLUGIN_PROJECT", false, false},
	{"--spec=", "PLUGIN_SPEC", false, false},
	{"--spec=", "PLUGIN_SPEC_PATH", false, false},
}

func GetDownloadCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs("tmpServerId",
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
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

	downloadCommandArgs := []string{"rt", "download", args.Target, args.Source}
	err = PopulateArgs(&downloadCommandArgs, &args, DownloadCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, downloadCommandArgs)

	return cmdList, nil
}
