package config

import (
	"backup/common"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

func BuildBackupConfig() (*common.BackupConfig, error) {
	consts := common.GetOSConstants()

	_, err := os.Stat(consts.ConfigLocation)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(consts.LoggingLocation)
	if err != nil {
		return nil, err
	}

	config, err := parseJSONConfig(consts)

	if err != nil {
		return nil, err
	}

	return &common.BackupConfig{
		KeyFile:        config.Key.Path + config.Key.FileName,
		Bucket:         config.S3.Bucket,
		Prefix:         config.S3.Prefix,
		Backup:         config.Backup,
		DecryptPath:    config.DecryptPath,
		Profile:        config.Profile,
		LastModifiedDt: parseLastModifiedFile(consts),
	}, nil
}

func parseJSONConfig(consts *common.BackupConstants) (*common.ConfigFile, error) {
	jsonFile, err := os.Open(consts.ConfigLocation + "config.json")
	defer jsonFile.Close()
	if err != nil {
		return nil, err
	}

	jsonData, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}

	dirname, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	keyObject := common.KeyObject{
		Path: dirname + "/.ssh/",
	}

	config := common.ConfigFile{
		Key:     keyObject,
		Profile: "default",
	}

	json.Unmarshal(jsonData, &config)

	return &config, nil
}

func parseLastModifiedFile(consts *common.BackupConstants) int {
	file, err := os.Open(consts.ConfigLocation + "last_run.conf")
	if err != nil {
		return -1
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return -1
	}

	timestamp, err := strconv.Atoi(string(data))

	if err != nil {
		return -1
	}

	return timestamp
}
