package config

import (
	"backup/common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

func BuildBackupConfig(consts *common.BackupConstants) (*common.BackupConfig, error) {
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
		S3Config:       config.S3,
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

	s3Object := common.S3Object{
		Tier: "STANDARD",
	}

	config := common.ConfigFile{
		Key:     keyObject,
		S3:      s3Object,
		Profile: "default",
	}

	json.Unmarshal(jsonData, &config)

	return &config, nil
}

func parseLastModifiedFile(consts *common.BackupConstants) int {
	file, err := os.Open(consts.ConfigLocation + common.LastRunFileName)
	if err != nil {
		return -1
	}

	defer file.Close()

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

func WriteLastModifiedFile(consts *common.BackupConstants) {
	file, err := os.OpenFile(consts.ConfigLocation+common.LastRunFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	now := time.Now()

	file.WriteString(fmt.Sprint(now.Unix()))
}
