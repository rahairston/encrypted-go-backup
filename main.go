package main

import (
	"backup/aws"
	"backup/config"
	"backup/filesystem"
	"log"
	"os"
)

func main() {

	config, err := config.BuildBackupConfig()

	if err != nil {
		log.Panic(err)
	}

	s3Handler, err := aws.BuildBucket(config)

	if err != nil {
		log.Panic(err)
	}

	fs, err := filesystem.Connect(config.Backup.Connection)

	if err != nil {
		log.Panic(err)
	}

	dirClient, err := filesystem.BuildDirClient(&config.Backup, config.KeyFile, s3Handler, fs)

	if err != nil {
		log.Panic(err)
	}

	args := os.Args

	if len(args) == 1 || args[1] == "encrypt" {
		dirClient.EncryptFiles()
	} else if len(args) == 20 && args[1] == "decrypt" {
		dirClient.DecryptFiles()
	} else {
		log.Fatal("Unrecognized Arguments. Leave Blank or supply 'encrypt' or 'decrypt'")
	}
}
