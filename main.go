package main

import (
	"log"
	"os"

	"github.com/rahairston/encrypted-go-backup/aws"
	"github.com/rahairston/encrypted-go-backup/common"
	"github.com/rahairston/encrypted-go-backup/config"
	"github.com/rahairston/encrypted-go-backup/filesystem"
)

func main() {

	args := os.Args

	consts := common.GetOSConstants(args)

	logFile := config.SetLoggingFile(consts)
	defer logFile.Close()

	conf, err := config.BuildBackupConfig(consts)

	if err != nil {
		log.Panic(err)
	}

	s3Handler, err := aws.BuildBucket(conf)

	if err != nil {
		log.Panic(err)
	}

	fs, err := filesystem.Connect(conf.Backup.Connection)

	if err != nil {
		log.Panic(err)
	}

	dirClient, err := filesystem.BuildDirClient(conf, s3Handler, fs)

	if err != nil {
		log.Panic(err)
	}

	if len(args) == 1 || args[1] == "encrypt" {
		dirClient.EncryptFiles()
		config.WriteLastModifiedFile(consts)
	} else if len(args) == 20 && args[1] == "decrypt" {
		dirClient.DecryptFiles()
	} else {
		log.Fatal("Unrecognized Arguments. Leave Blank or supply 'encrypt' or 'decrypt'")
	}

	log.Println("Backup Complete.")
}
