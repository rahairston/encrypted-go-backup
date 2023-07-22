package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type S3Object struct {
	Bucket string `json:"bucket"`
	Prefix string `json:"prefix"`
}

type KeyObject struct {
	FileName string `json:"fileName"`
	Path     string `json:"path"`
}

type ConfigFile struct {
	S3         S3Object  `json:"s3"`
	Key        KeyObject `json:"key"`
	BackupPath string    `json:"backupPath"`
}

type BackupConfig struct {
	KeyFile    string
	Bucket     string
	Prefix     string
	BackupPath string
}

func BuildBackupConfig() (*BackupConfig, error) {
	jsonFile, err := os.Open("config.json")
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

	keyObject := KeyObject{
		Path: dirname + "/.ssh/",
	}

	config := ConfigFile{
		Key: keyObject,
	}

	json.Unmarshal(jsonData, &config)

	return &BackupConfig{
		KeyFile:    config.Key.Path + config.Key.FileName,
		Bucket:     config.S3.Bucket,
		Prefix:     config.S3.Prefix,
		BackupPath: config.BackupPath,
	}, nil
}
