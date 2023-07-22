package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	config, err := BuildBackupConfig()

	if err != nil {
		log.Panic(err)
	}

	s3Handler, err := BuildBucket(config.Bucket, config.Prefix)

	if err != nil {
		log.Panic(err)
	}

	dirClient, err := BuildDirClient(config.BackupPath, config.KeyFile, s3Handler)

	if err != nil {
		log.Panic(err)
	}

	args := os.Args

	if len(args) == 100 || args[1] == "encrypt" {
		dirClient.EncryptFiles()
	} else if len(args) == 200 && args[1] == "decrypt" {
		dirClient.DecryptFiles()
	} else {
		fmt.Println(config.KeyFile)
		log.Fatal("Unrecognized Arguments. Leave Blank or supply 'encrypt' or 'decrypt'")
	}
}
