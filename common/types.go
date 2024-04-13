package common

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type FileSystem interface {
	GetFileNames(path string, exclusions ExcludeObject, lastModifiedDt int64) []string
	ValidatePath(path string) string
	ValidatePaths(path string, folders []string) []string
	ReadFile(fileName string) ([]byte, error)
	Close()
}

type FileConfigType string

const (
	Local FileConfigType = "local"
	Smb                  = "smb"
)

type TierException struct {
	Tier    types.StorageClass `json:"tier"`
	Matches []string           `json:"matches"`
}

type S3TierObject struct {
	Default types.StorageClass `json:"default"`
	Files   []TierException    `json:"files"`
	Folders []TierException    `json:"folders"`
}

type S3Object struct {
	Bucket string       `json:"bucket"`
	Prefix string       `json:"prefix"`
	Tier   S3TierObject `json:"tier"`
}

type KeyObject struct {
	FileName string `json:"fileName"`
	Path     string `json:"path"`
}

type Authentication struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SmbConfig struct {
	Host           string         `json:"host"`
	Port           string         `json:"port"`
	Authentication Authentication `json:"authentication"`
	MountPoint     string         `json:"mountPoint"`
}

type ConnectionObject struct {
	Type      FileConfigType `json:"type"`
	SmbConfig SmbConfig      `json:"smbConfig"`
}

type ExcludeObject struct {
	Files   []string `json:"files"`
	Folders []string `json:"folders"`
}

type BackupObject struct {
	BasePath   string           `json:"basePath"`
	Folders    []string         `json:"folders"`
	Connection ConnectionObject `json:"connection"`
	Exclusions ExcludeObject    `json:"exclude"`
}

type ConfigFile struct {
	S3          S3Object     `json:"s3"`
	Key         KeyObject    `json:"key"`
	Backup      BackupObject `json:"backup"`
	DecryptPath string       `json:"decryptPath"`
	Profile     string       `json:"profile"`
}

type BackupConfig struct {
	KeyFile        string
	S3Config       S3Object
	Backup         BackupObject
	DecryptPath    string
	Profile        string
	LastModifiedDt int64
}
