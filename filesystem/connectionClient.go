package filesystem

import (
	"backup/types"
)

func Connect(connectionConfig types.ConnectionObject) (types.FileSystem, error) {

	switch connectionConfig.Type {
	case types.Smb:
		return SmbConnect(connectionConfig.SmbConfig)
	case types.Local:
		return &LocalClient{}, nil
	default:
		return nil, nil
	}
}
