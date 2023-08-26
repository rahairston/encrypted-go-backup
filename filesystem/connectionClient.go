package filesystem

import (
	"github.com/rahairston/encrypted-go-backup/common"
)

func Connect(connectionConfig common.ConnectionObject) (common.FileSystem, error) {

	switch connectionConfig.Type {
	case common.Smb:
		return SmbConnect(connectionConfig.SmbConfig)
	case common.Local:
		return &LocalClient{}, nil
	default:
		return nil, nil
	}
}
