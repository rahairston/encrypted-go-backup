package filesystem

import (
	"backup/types"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/hirochachacha/go-smb2"
)

type SmbClient struct {
	s    *smb2.Session
	fs   *smb2.Share
	conn net.Conn
}

func SmbConnect(config types.SmbConfig) (*SmbClient, error) {
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

func (smbClient SmbClient) GetFileNames(path string) []string {
	var result []string
	var adjustedPath string = path
	if !strings.HasSuffix(path, "\\") {
		adjustedPath = path + "\\"
	}
	// defer smbClient.conn.Close()
	// defer smbClient.s.Logoff()

	files, _ := smbClient.fs.ReadDir(path)

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.IsDir() {
			result = append(result, smbClient.GetFileNames(adjustedPath+file.Name())...)
		} else {
			result = append(result, adjustedPath+file.Name())
			_, err := smbClient.fs.ReadFile(adjustedPath + file.Name())

			if err != nil {
				fmt.Println("erro")
				panic(err)
			}

			fmt.Println(adjustedPath + file.Name())
		}
	}

	return result
}

func (smbClient SmbClient) ValidatePath(path string) string {
	info, err := smbClient.fs.Stat(path)

	if err != nil {
		panic(err)
	}

	if !info.IsDir() {
		panic(errors.New("Path provided must be a Folder."))
	}

	if !strings.HasSuffix(path, "\\") {
		return path + "\\"
	}

	return path
}

func (smbClient SmbClient) ReadFile(fileName string) ([]byte, error) {
	return smbClient.fs.ReadFile(fileName)
}
