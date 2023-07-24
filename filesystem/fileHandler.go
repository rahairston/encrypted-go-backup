package filesystem

import (
	"backup/aws"
	"backup/encryption"
	"backup/types"
	"log"
	"strings"
)

type DirClient struct {
	path      *string
	keys      *encryption.KeyHandler
	s3Handler *aws.BucketHandler
	fs        types.FileSystem
}

func BuildDirClient(path string, keyFileName string, s3Handler *aws.BucketHandler, fs types.FileSystem) (*DirClient, error) {

	adjustedPath := fs.ValidatePath(path)

	keys, err := encryption.BuildKeyHandler(keyFileName)

	if err != nil {
		return nil, err
	}

	return &DirClient{
		path:      &adjustedPath,
		keys:      keys,
		s3Handler: s3Handler,
		fs:        fs,
	}, nil
}

func (dir DirClient) EncryptFiles() {
	fileNames := dir.fs.GetFileNames(*dir.path)

	c := make(chan string, len(fileNames))

	for _, fileName := range fileNames {
		go dir.EncryptAndUploadFile(fileName, c)
	}

	for i := 0; i < cap(c); i++ {
		log.Println(<-c)
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

	var adjustedS3Key = strings.TrimPrefix(fileName, *dir.path)

	err = dir.s3Handler.PutObject(adjustedS3Key, encrypted)
	if err != nil {
		log.Panic(err)
	}

	c <- fileName
}

func (dir DirClient) DecryptFiles() {
	fileNames := dir.fs.GetFileNames(*dir.path)

	c := make(chan string, len(fileNames))

	for _, fileName := range fileNames {
		go dir.DownloadAndDecryptFile(fileName, c)
	}

	for i := 0; i < cap(c); i++ {
		log.Println(<-c)
	}
}

func (dir DirClient) DownloadAndDecryptFile(fileName string, c chan string) {
	var adjustedS3Key = strings.TrimPrefix(fileName, *dir.path)

	data, err := dir.s3Handler.GetObject(adjustedS3Key)
	if err != nil {
		log.Panic(err)
	}

	_, err = dir.keys.Decrypt(data)

	if err != nil {
		log.Panic(err)
	}

	// TODO Write file to decrypted path
	// ioutil.WriteFile(path + fileName, decrypted, 0777)

	c <- fileName
}
