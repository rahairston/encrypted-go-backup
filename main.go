package main

import (
	"backup/aws"
	"backup/common"
	"backup/config"
	"backup/filesystem"
	"log"
	"os"
)

func main() {

	consts := common.GetOSConstants()

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

	args := os.Args

	if len(args) == 1 || args[1] == "encrypt" {
		dirClient.EncryptFiles()
		config.WriteLastModifiedFile(consts)
	} else if len(args) == 20 && args[1] == "decrypt" {
		dirClient.DecryptFiles()
	} else {
		log.Fatal("Unrecognized Arguments. Leave Blank or supply 'encrypt' or 'decrypt'")
	}
}
