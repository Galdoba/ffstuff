package main

import (
	"os"
	"path/filepath"
)

func defaultStorageDir() string {
	home, _ := os.UserHomeDir()
	sep := string(filepath.Separator)
	return home + sep + ".ffmpeg" + sep + "data" + sep + programName + sep
}

func defaultLogFile() string {
	home, _ := os.UserHomeDir()
	sep := string(filepath.Separator)
	return home + sep + ".ffmpeg" + sep + "logs" + sep + programName + ".log"
}
