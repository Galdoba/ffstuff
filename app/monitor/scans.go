package main

import (
	"bufio"
	"os"
	"sort"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/fdf"
	"github.com/Galdoba/ffstuff/pkg/mdm/inputinfo"
	"github.com/Galdoba/ffstuff/pkg/namedata"
)

func updateStoredInfo(list []string) error {
	dataStore, err := os.OpenFile(storagePath+storageFile, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	oldData := []string{}
	newData := []string{}
	scanner := bufio.NewScanner(dataStore)
	for scanner.Scan() {
		oldData = append(oldData, scanner.Text())
	}
	dataStore.Close()
	oldText := strings.Join(oldData, "\n")
	dHash := makeHash(oldData)
	for _, path := range list {
		data := dHash[path]
		toAdd := ""
		fInfo, err := os.Stat(path)
		if err != nil {
			toAdd += "uptodate:false|error:" + err.Error() + "|"
			newData = append(newData, path+"  "+data+toAdd)
			continue
		}
		if fInfo.IsDir() {
			continue
		}

		//		pi, err := inputinfo.ParseFile(path)
		en := namedata.EditForm(path)
		if err == nil {
			if !hasKey(allKeys(data), "mProfile") {
				pi, _ := inputinfo.ParseFile(path)
				toAdd += "mProfile:" + fdf.FMP(pi) + "|"
			}
			if !hasKey(allKeys(data), "fSize") {
				toAdd += "fSize:" + fdf.Size(path) + "|"
			}
		}
		if en.EditName() != "" {
			if !hasKey(allKeys(data), "mTag") {
				toAdd += "mTag:" + en.ContentType() + "|"
				if en.ContentType() == "SER" {
					toAdd += "season:" + en.Season() + "|"
					toAdd = strings.TrimSuffix(toAdd, "season:|")
					toAdd += "episode:" + en.Episode() + "|"
					toAdd = strings.TrimSuffix(toAdd, "episode:|")
				}
			}
			if !hasKey(allKeys(data), "editName") {
				toAdd += "editName:" + en.EditName() + "|"
			}
		}
		newData = append(newData, path+"  "+data+toAdd)
	}
	sort.Strings(newData)
	newText := strings.Join(newData, "\n")
	if newText == oldText {
		return nil
	}
	os.RemoveAll(storagePath + storageFile)
	dataStore, err = os.OpenFile(storagePath+storageFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer dataStore.Close()

	_, err = dataStore.WriteString(strings.Join(newData, "\n"))
	if err != nil {
		return err
	}
	return nil

}

func hasMP(s string) bool {
	switch len(strings.Split(s, "  ")) {
	default:
	case 2:
		return true
	}
	return false
}

func allKeys(data string) []string {
	keys := []string{}
	dParts := strings.Split(data, "  ")
	for _, check := range dParts {
		pairs := strings.Split(check, "|")
		for _, pair := range pairs {
			key := strings.Split(pair, ":")
			switch len(key) {
			case 2:
				keys = append(keys, key[0])
			}
		}
	}
	return keys
}

func hasKey(allKeys []string, key string) bool {
	for _, k := range allKeys {
		if k == key {
			return true
		}
	}
	return false
}

func storedPath(data string) string {
	return strings.Split(data, "  ")[0]
}

func isLinkedToSheet(name string) bool {
	if len(strings.Split(name, "--")) == 3 {
		return true
	}
	return false
}

func makeHash(data []string) map[string]string {
	mp := make(map[string]string)
	for _, d := range data {
		dt := strings.Split(d, "  ")
		mp[dt[0]] = dt[1]
	}
	return mp
}
