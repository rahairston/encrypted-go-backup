package filesystem

import (
	"log"
	"strings"

	"github.com/rahairston/encrypted-go-backup/aws"
	"github.com/rahairston/encrypted-go-backup/common"
	"github.com/rahairston/encrypted-go-backup/encryption"
)

type DirClient struct {
	base           *string
	paths          *[]string
	keys           *encryption.KeyHandler
	s3Handler      *aws.BucketHandler
	fs             common.FileSystem
	exclusions     common.ExcludeObject
	lastModifiedDt int64
}

func BuildDirClient(backupConfig *common.BackupConfig,
	s3Handler *aws.BucketHandler, fs common.FileSystem) (*DirClient, error) {

	path := fs.ValidatePath(backupConfig.Backup.BasePath)
	folders := backupConfig.Backup.Folders

	adjustedPaths := fs.ValidatePaths(path, folders)

	keys, err := encryption.BuildKeyHandler(backupConfig.KeyFile)

	if err != nil {
		return nil, err
	}

	return &DirClient{
		base:           &path,
		paths:          &adjustedPaths,
		keys:           keys,
		s3Handler:      s3Handler,
		fs:             fs,
		exclusions:     backupConfig.Backup.Exclusions,
		lastModifiedDt: backupConfig.LastModifiedDt,
	}, nil
}

func (dir DirClient) EncryptFiles() {
	defer dir.fs.Close()

	for _, entry := range *dir.paths {
		fileNames := dir.fs.GetFileNames(entry, dir.exclusions, dir.lastModifiedDt)

		c := make(chan string, len(fileNames))

		for _, fileName := range fileNames {
			go dir.EncryptAndUploadFile(fileName, c)
		}

		for i := 0; i < cap(c); i++ {
			log.Println(<-c)
		}
	}

}

func (dir DirClient) EncryptAndUploadFile(fileName string, c chan string) {

	plaintext, err := dir.fs.ReadFile(fileName)

	if err != nil {
		log.Panic(err)
	}

	encrypted, err := dir.keys.Encrypt(plaintext)

	if err != nil {
		log.Panic(fileName, err)
	}

	var adjustedS3Key = strings.TrimPrefix(fileName, *dir.base)

	err = dir.s3Handler.PutObject(adjustedS3Key, encrypted)
	if err != nil {
		log.Panic(err)
	}

	c <- fileName
}

// Maybe instead we List all items in S3 Bucket
// then go Download and Decrypt into location
func (dir DirClient) DecryptFiles() {
	defer dir.fs.Close()
	for _, entry := range *dir.paths {
		fileNames := dir.fs.GetFileNames(entry, dir.exclusions, dir.lastModifiedDt)

		c := make(chan string, len(fileNames))

		for _, fileName := range fileNames {
			go dir.DownloadAndDecryptFile(fileName, c)
		}

		for i := 0; i < cap(c); i++ {
			log.Println(<-c)
		}
	}
}

func (dir DirClient) DownloadAndDecryptFile(fileName string, c chan string) {
	var adjustedS3Key = strings.TrimPrefix(fileName, *dir.base)

	data, err := dir.s3Handler.GetObject(adjustedS3Key)
	if err != nil {
		log.Panic(err)
	}

	_, err = dir.keys.Decrypt(data)

	if err != nil {
		log.Panic(err)
	}

	// TODO Write file to decrypted path
	// os.WriteFile(path + fileName, decrypted, 0777)

	c <- fileName
}
