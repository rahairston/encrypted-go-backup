package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/rahairston/encrypted-go-backup/common"
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

	s3TierObject := common.S3TierObject{
		Default: "STANDARD",
	}

	s3Object := common.S3Object{
		Tier: s3TierObject,
	}

	config := common.ConfigFile{
		Key:     keyObject,
		S3:      s3Object,
		Profile: "default",
	}

	json.Unmarshal(jsonData, &config)

	return &config, nil
}

func parseLastModifiedFile(consts *common.BackupConstants) int64 {
	file, err := os.Open(consts.ConfigLocation + common.LastRunFileName)
	if err != nil {
		return -1
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return -1
	}

	timestamp, err := strconv.ParseInt(string(data), 10, 64)

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

func SetLoggingFile(consts *common.BackupConstants) *os.File {
	now := time.Now().UTC()

	logFile, err := os.OpenFile(consts.LoggingLocation+now.Format("2006-01-02")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Panic(err)
	}

	log.SetOutput(logFile)

	log.Println("Starting...")

	return logFile
}
