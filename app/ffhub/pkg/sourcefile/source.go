package sourcefile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	tag_separator = "--"
)

func CheckProject(files []string) error {
	//check path
	for _, path := range files {
		if err := checkPath(path); err != nil {
			return err
		}
	}
	//check types

}

func checkPath(path string) error {
	phase := "check path"
	fs, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%v: %v", phase, err)
	}
	if fs.IsDir() {
		return fmt.Errorf("%v: %v is directory", phase, fs.Name())
	}
	f, err := os.OpenFile(path, os.O_RDWR, 0777)
	if err != nil {
		return fmt.Errorf("%v: %v", phase, err)
	}
	defer f.Close()
	return nil
}

type tagInfo struct {
	baseName     string
	taskType     string
	sourceName   string
	season_num   string
	episode_num  string
	episode_part string
}

func decodeTags(paths []string) map[string]tagInfo {
	for _, path := range paths {
		name := filepath.Base(path)
		parts := strings.Split(name, tag_separator)
		tInfo := tagInfo{}
		tInfo.sourceName = parts[len(parts)-1]
		for _, part := range parts {
			switch part {
			case "TRL", "FLM", "SER":
				tInfo.taskType = part

			}
		}
	}

}

func decodeEpisodeInfo(str string) (string, string, string) {
	if !strings.HasPrefix(str, "s") || !strings.Contains(str, "e") {
		return "", "", ""
	}
	prts := strings.Split(str, "e")

}
