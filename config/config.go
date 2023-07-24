package config

import (
	"backup/types"
	"encoding/json"
	"io/ioutil"
	"os"
)

func BuildBackupConfig() (*types.BackupConfig, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	jsonFile, err := os.Open(cwd + "/config.json")
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

	return &types.BackupConfig{
		KeyFile:     config.Key.Path + config.Key.FileName,
		Bucket:      config.S3.Bucket,
		Prefix:      config.S3.Prefix,
		Backup:      config.Backup,
		DecryptPath: config.DecryptPath,
		Profile:     config.Profile,
	}, nil
}
