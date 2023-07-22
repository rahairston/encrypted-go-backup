package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type DirClient struct {
	path      *string
	keys      *KeyHandler
	s3Handler *BucketHandler
}

func BuildDirClient(path string, keyFileName string, s3Handler *BucketHandler) (*DirClient, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("Path provided must be a Folder.")
	}

	keys, err := BuildKeyHandler(keyFileName)

	if err != nil {
		return nil, err
	}

	var adjustedPath string = path
	if !strings.HasSuffix(path, "/") {
		adjustedPath = path + "/"
	}

	return &DirClient{
		path:      &adjustedPath,
		keys:      keys,
		s3Handler: s3Handler,
	}, nil
}

func GetFileNames(path string) []string {
	var result []string
	var adjustedPath string = path
	if !strings.HasSuffix(path, "/") {
		adjustedPath = path + "/"
	}

	entries, err := os.ReadDir(adjustedPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if e.IsDir() {
			result = append(result, GetFileNames(adjustedPath+e.Name())...)
		} else {
			result = append(result, adjustedPath+e.Name())
		}
	}

	return result
}

func (dir DirClient) EncryptFiles() {
	fileNames := GetFileNames(*dir.path)

	c := make(chan string, len(fileNames))

	for _, fileName := range fileNames {
		go dir.EncryptAndUploadFile(fileName, c)
	}

	for i := 0; i < cap(c); i++ {
		log.Println(<-c)
	}
}

func (dir DirClient) EncryptAndUploadFile(fileName string, c chan string) {
	plaintext, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Panic(err)
	}

	encrypted, err := dir.keys.encrypt(plaintext)

	if err != nil {
		log.Panic(fileName, err)
	}
	var adjustedS3Key = strings.TrimPrefix(fileName, *dir.path)

	err = dir.s3Handler.putObject(adjustedS3Key, encrypted)
	if err != nil {
		log.Panic(err)
	}

	c <- fileName
}

func (dir DirClient) DecryptFiles() {
	fileNames := GetFileNames(*dir.path)

	c := make(chan []byte, len(fileNames))

	for _, fileName := range fileNames {
		go dir.DownloadAndDecryptFile(fileName, c)
	}

	x := <-c

	ioutil.WriteFile("test_encryption.txt", x, 0777)

	test_decrypt, _ := dir.keys.decrypt(x)

	ioutil.WriteFile("test.pdf", test_decrypt, 0777)
}

func (dir DirClient) DownloadAndDecryptFile(fileName string, c chan []byte) {
	plaintext, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Panic(err)
	}

	encrypted, err := dir.keys.encrypt(plaintext)

	if err != nil {
		log.Panic(fileName, err)
	}
	c <- encrypted
	// var adjustedS3Key = strings.TrimPrefix(fileName, *dir.path)

	// err = dir.s3Handler.putObject(adjustedS3Key, encrypted)
	// if err != nil {
	// 	log.Panic(err)
	// }
}
