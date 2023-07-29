package config

import (
	"backup/constants"
	"backup/types"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

func BuildBackupConfig() (*types.BackupConfig, error) {
	consts := constants.GetOSConstants()

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

	return &types.BackupConfig{
		KeyFile:        config.Key.Path + config.Key.FileName,
		Bucket:         config.S3.Bucket,
		Prefix:         config.S3.Prefix,
		Backup:         config.Backup,
		DecryptPath:    config.DecryptPath,
		Profile:        config.Profile,
		LastModifiedDt: 0,
	}, nil
}

func parseJSONConfig(consts *constants.BackupConstants) (*types.ConfigFile, error) {
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

	keyObject := types.KeyObject{
		Path: dirname + "/.ssh/",
	}

	config := types.ConfigFile{
		Key:     keyObject,
		Profile: "default",
	}

	json.Unmarshal(jsonData, &config)

	return &config, nil
}

func parseLastModifiedFile(consts *constants.BackupConstants) int {
	file, err := os.Open(consts.ConfigLocation + "LastRun.conf")
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
