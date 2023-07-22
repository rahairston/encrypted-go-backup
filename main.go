package main

import (
	"log"
	"os"
)

func main() {

	config, err := BuildBackupConfig()

	if err != nil {
		log.Panic(err)
	}

	s3Handler, err := BuildBucket(config)

	if err != nil {
		log.Panic(err)
	}

	dirClient, err := BuildDirClient(config.BackupPath, config.KeyFile, s3Handler)

	if err != nil {
		log.Panic(err)
	}

	args := os.Args

	if len(args) == 1 || args[1] == "encrypt" {
		dirClient.EncryptFiles()
	} else if len(args) == 2 && args[1] == "decrypt" {
		dirClient.DecryptFiles()
	} else {
		log.Fatal("Unrecognized Arguments. Leave Blank or supply 'encrypt' or 'decrypt'")
	}
}
