package plugin

import (
	"fmt"
	"os"
	"time"
)

// rtBuildInfo

// rtPublishBuildInfo

var UploadCmdJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--access-token=", "PLUGIN_ACCESS_TOKEN", false, false},
	{"--ant=", "PLUGIN_ANT", false, false},
	{"--archive=", "PLUGIN_ARCHIVE", false, false},
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false},
	{"--chunk-size=", "PLUGIN_CHUNK_SIZE", false, false},
	{"--client-cert-key-path=", "PLUGIN_CLIENT_CERT_KEY_PATH", false, false},
	{"--client-cert-path=", "PLUGIN_CLIENT_CERT_PATH", false, false},
	{"--deb=", "PLUGIN_DEB", false, false},
	{"--detailed-summary=", "PLUGIN_DETAILED_SUMMARY", false, false},
	{"--dry-run=", "PLUGIN_DRY_RUN", false, false},
	{"--exclusions=", "PLUGIN_EXCLUSIONS", false, false},
	{"--explode=", "PLUGIN_EXPLODE", false, false},
	{"--fail-no-op=", "PLUGIN_FAIL_NO_OP", false, false},
	{"--include-dirs=", "PLUGIN_INCLUDE_DIRS", false, false},
	{"--insecure-tls=", "PLUGIN_INSECURE_TLS", false, false},
	{"--min-split=", "PLUGIN_MIN_SPLIT", false, false},
	{"--module=", "PLUGIN_MODULE", false, false},
	{"--project=", "PLUGIN_PROJECT", false, false},
	{"--quiet=", "PLUGIN_QUIET", false, false},
	{"--recursive=", "PLUGIN_RECURSIVE", false, false},
	{"--regexp=", "PLUGIN_REGEXP", false, false},
	{"--retry-wait-time=", "PLUGIN_RETRY_WAIT_TIME", false, false},
	{"--server-id=", "PLUGIN_SERVER_ID", false, false},
	{"--spec=", "PLUGIN_SPEC", false, false},
	{"--spec=", "PLUGIN_SPEC_PATH", false, false},
	{"--spec-vars=", "PLUGIN_SPEC_VARS", false, false},
	{"--split-count=", "PLUGIN_SPLIT_COUNT", false, false},
	{"--ssh-key-path=", "PLUGIN_SSH_KEY_PATH", false, false},
	{"--ssh-passphrase=", "PLUGIN_SSH_PASSPHRASE", false, false},
	{"--symlinks=", "PLUGIN_SYMLINKS", false, false},
	{"--sync-deletes=", "PLUGIN_SYNC_DELETES", false, false},
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

	for _, cmd := range cmdList {
		fmt.Println(cmd)
	}

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
	{"--access-token=", "PLUGIN_ACCESS_TOKEN", false, false},
	{"--archive-entries=", "PLUGIN_ARCHIVE_ENTRIES", false, false},
	{"--build=", "PLUGIN_BUILD", false, false},
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false},
	{"--bundle=", "PLUGIN_BUNDLE", false, false},
	{"--bypass-archive-inspection=", "PLUGIN_BYPASS_ARCHIVE_INSPECTION", false, false},
	{"--client-cert-key-path=", "PLUGIN_CLIENT_CERT_KEY_PATH", false, false},
	{"--client-cert-path=", "PLUGIN_CLIENT_CERT_PATH", false, false},
	{"--detailed-summary=", "PLUGIN_DETAILED_SUMMARY", false, false},
	{"--dry-run=", "PLUGIN_DRY_RUN", false, false},
	{"--exclude-artifacts=", "PLUGIN_EXCLUDE_ARTIFACTS", false, false},
	{"--exclude-props=", "PLUGIN_EXCLUDE_PROPS", false, false},
	{"--exclusions=", "PLUGIN_EXCLUSIONS", false, false},
	{"--explode=", "PLUGIN_EXPLODE", false, false},
	{"--fail-no-op=", "PLUGIN_FAIL_NO_OP", false, false},
	{"--flat=", "PLUGIN_FLAT", false, false},
	{"--gpg-key=", "PLUGIN_GPG_KEY", false, false},
	{"--include-deps=", "PLUGIN_INCLUDE_DEPS", false, false},
	{"--include-dirs=", "PLUGIN_INCLUDE_DIRS", false, false},
	{"--insecure-tls=", "PLUGIN_INSECURE_TLS", false, false},
	{"--limit=", "PLUGIN_LIMIT", false, false},
	{"--min-split=", "PLUGIN_MIN_SPLIT", false, false},
	{"--module=", "PLUGIN_MODULE", false, false},
	{"--offset=", "PLUGIN_OFFSET", false, false},
	{"--project=", "PLUGIN_PROJECT", false, false},
	{"--props=", "PLUGIN_PROPS", false, false},
	{"--quiet=", "PLUGIN_QUIET", false, false},
	{"--recursive=", "PLUGIN_RECURSIVE", false, false},
	{"--retries=", "PLUGIN_RETRIES", false, false},
	{"--retry-wait-time=", "PLUGIN_RETRY_WAIT_TIME", false, false},
	{"--server-id=", "PLUGIN_SERVER_ID", false, false},
	{"--skip-checksum=", "PLUGIN_SKIP_CHECKSUM", false, false},
	{"--sort-by=", "PLUGIN_SORT_BY", false, false},
	{"--sort-order=", "PLUGIN_SORT_ORDER", false, false},
	{"--spec=", "PLUGIN_SPEC", false, false},
	{"--spec-vars=", "PLUGIN_SPEC_VARS", false, false},
	{"--split-count=", "PLUGIN_SPLIT_COUNT", false, false},
	{"--ssh-key-path=", "PLUGIN_SSH_KEY_PATH", false, false},
	{"--ssh-passphrase=", "PLUGIN_SSH_PASSPHRASE", false, false},
	{"--sync-deletes=", "PLUGIN_SYNC_DELETES", false, false},
	{"--threads=", "PLUGIN_THREADS", false, false},
	{"--validate-symlinks=", "PLUGIN_VALIDATE_SYMLINKS", false, false},
}

func GetDownloadCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs("tmpServerId",
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		return cmdList, err
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
