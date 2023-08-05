package common

import (
	"os"
	"runtime"
)

type BackupConstants struct {
	LoggingLocation string
	ConfigLocation  string
}

const (
	Separator       string = string(os.PathSeparator)
	LastRunFileName string = "last_run.conf"
)

func GetOSConstants() *BackupConstants {
	opsys := runtime.GOOS
	switch opsys {
	case "windows":
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return &BackupConstants{
			LoggingLocation: home + "\\AppData\\Local\\backitup\\",
			ConfigLocation:  home + "\\AppData\\Local\\backitup\\",
		}
	default:
		return &BackupConstants{
			LoggingLocation: "/var/log/backitup/",
			ConfigLocation:  "/etc/backitup/",
		}
	}
}
