package filesystem

import (
	"backup/common"
	"errors"
	"net"
	"strings"

	"github.com/hirochachacha/go-smb2"
)

type SmbClient struct {
	s    *smb2.Session
	fs   *smb2.Share
	conn net.Conn
}

func SmbConnect(config common.SmbConfig) (*SmbClient, error) {
	var port = config.Port
	if port == "" {
		port = "445"
	}

	conn, err := net.Dial("tcp", config.Host+":"+config.Port)

	if err != nil {
		return nil, err
	}

	d := smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     config.Authentication.Username,
			Password: config.Authentication.Password,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		return nil, err
	}

	fs, err := s.Mount(config.MountPoint)

	return &SmbClient{
		s:    s,
		conn: conn,
		fs:   fs,
	}, nil
}

func (smbClient SmbClient) GetFileNames(path string, exclusions common.ExcludeObject) []string {
	var result []string
	var adjustedPath string = path
	if !strings.HasSuffix(path, "\\") { // Keep \\ since SMB is Windows file pathing
		adjustedPath = path + "\\"
	}

	files, _ := smbClient.fs.ReadDir(path)

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		} else if file.IsDir() && !common.ShouldBeExcluded(file.Name(), exclusions.Folders) {
			result = append(result, smbClient.GetFileNames(adjustedPath+file.Name(), exclusions)...)
		} else if !file.IsDir() && !common.ShouldBeExcluded(file.Name(), exclusions.Files) {
			result = append(result, adjustedPath+file.Name())
			_, err := smbClient.fs.ReadFile(adjustedPath + file.Name())

			if err != nil {
				panic(err)
			}
		}
	}

	return result
}

func (smbClient SmbClient) ValidatePath(path string) string {
	info, err := smbClient.fs.Stat(path)

	if err != nil {
		panic(err)
	} else if !info.IsDir() {
		panic(errors.New("Path provided must be a Folder."))
	} else if !strings.HasSuffix(path, "\\") { // Keep \\ since SMB is Windows file pathing
		return path + "\\"
	}

	return path
}

func (smbClient SmbClient) ReadFile(fileName string) ([]byte, error) {
	return smbClient.fs.ReadFile(fileName)
}

func (smbClient SmbClient) Close() {
	smbClient.s.Logoff()
	smbClient.conn.Close()
}
