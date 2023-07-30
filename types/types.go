package types

type FileSystem interface {
	GetFileNames(path string, exclusions ExcludeObject) []string
	ValidatePath(path string) string
	ReadFile(fileName string) ([]byte, error)
	Close()
}

type FileConfigType string

const (
	Local FileConfigType = "local"
	Smb                  = "smb"
)

type S3Object struct {
	Bucket string `json:"bucket"`
	Prefix string `json:"prefix"`
	Tier   string `json:"tier"`
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
	Path       string           `json:"path"`
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
	Bucket         string
	Prefix         string
	Backup         BackupObject
	DecryptPath    string
	Profile        string
	LastModifiedDt int
}
