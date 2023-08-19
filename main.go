package main

import (
	"log"
	"os"

	"github.com/rahairston/EncryptedGoBackup/aws"
	"github.com/rahairston/EncryptedGoBackup/common"
	"github.com/rahairston/EncryptedGoBackup/config"
	"github.com/rahairston/EncryptedGoBackup/filesystem"
)

func main() {

	consts := common.GetOSConstants()

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
