package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func renameAmediaFiles(dir string) error {
	amediaFiles, err := listAmediaFiles(dir)
	if err != nil {
		return err
	}
	for _, fileName := range amediaFiles {

		if err := normalizeName(dir + fileName); err != nil {
			return fmt.Errorf("rename failed: %v", err)
		}
	}
	return nil
}

func normalizeName(path string) error {
	name := filepath.Base(path)
	dir := filepath.Dir(path)
	sep := string(filepath.Separator)
	if strings.Contains(name, "--") {
		return nil
	}
	re := regexp.MustCompile(`(_s[0-9]{1,}e[0-9]{1,}_PRT[0-9]{1,})`)
	amediaTags := re.FindString(name)
	base := headOfString(name, amediaTags)
	tags := strings.Split(amediaTags, "_")
	newName := base + "--" + "SER" + "--" + tags[1] + "--" + name
	if err := os.Rename(path, dir+sep+newName); err != nil {
		return err
	}

	return nil
}

func headOfString(str, separator string) string {
	parts := strings.Split(str, separator)
	return parts[0]
}

func listAmediaFiles(dir string) ([]string, error) {
	fileList := []string{}

	fi, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("scan: %v", err)
	}
	for _, f := range fi {
		if f.IsDir() {
			continue
		}
		if isAmediaFile(f.Name()) {
			fileList = append(fileList, f.Name())
		}
	}
	return fileList, nil
}

func isAmediaFile(name string) bool {
	re := regexp.MustCompile(`(_s[0-9]{1,}e[0-9]{1,}_PRT[0-9]{1,})`)
	amediaTags := re.FindString(name)
	return amediaTags != ""
}
