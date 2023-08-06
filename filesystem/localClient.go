package filesystem

import (
	"backup/common"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
		panic(errors.New("Path provided must be a Folder."))
	} else if !strings.HasSuffix(path, common.Separator) {
		return path + common.Separator
	}

	return path
}

func (lc LocalClient) ReadFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

func (lc LocalClient) Close() {
}
