package filesystem

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/rahairston/encrypted-go-backup/common"
)

type LocalClient struct {
}

func (lc LocalClient) GetFileNames(path string, exclusions common.ExcludeObject, lastModifiedDt int64) []string {
	var result []string
	var adjustedPath string = path
	if !strings.HasSuffix(path, common.Separator) {
		adjustedPath = path + common.Separator
	}

	entries, err := os.ReadDir(adjustedPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		fileInfo, _ := e.Info()
		if strings.HasPrefix(e.Name(), ".") {
			continue
		} else if e.IsDir() && !common.ShouldBeExcluded(e.Name(), exclusions.Folders) {
			result = append(result, lc.GetFileNames(adjustedPath+e.Name(), exclusions, lastModifiedDt)...)
		} else if !e.IsDir() && !common.ShouldBeExcluded(e.Name(), exclusions.Files) && fileInfo.ModTime().Unix() > lastModifiedDt {
			result = append(result, adjustedPath+e.Name())
		}
	}

	return result
}

func (lc LocalClient) ValidatePath(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	} else if !info.IsDir() {
		panic(errors.New("path provided must be a folder"))
	} else if !strings.HasSuffix(path, common.Separator) {
		return path + common.Separator
	}

	return path
}

func (lc LocalClient) ValidatePaths(path string, folders []string) []string {
	var result []string
	info, err := os.Stat(path)

	adjustedPath := path

	if err != nil {
		panic(err)
	} else if !info.IsDir() {
		panic(errors.New("base path provided must be a folder"))
	} else if !strings.HasSuffix(path, common.Separator) {
		adjustedPath += common.Separator
	}

	if len(folders) == 0 {
		result = append(result, adjustedPath)
	}

	for _, entry := range folders {
		folder := entry
		if strings.HasPrefix(entry, common.Separator) {
			folder = strings.TrimPrefix(entry, common.Separator)
		}

		result = append(result, lc.ValidatePath(adjustedPath+folder))
	}

	return result
}

func (lc LocalClient) ReadFile(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}

func (lc LocalClient) Close() {
}
