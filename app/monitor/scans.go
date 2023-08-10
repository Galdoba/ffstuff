package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/fdf"
	"github.com/Galdoba/ffstuff/pkg/mdm/inputinfo"
)

func updateMediaProfile() error {
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

	dataLen := len(oldData)
	for i, data := range oldData {
		fmt.Printf("updating data: %v/%v                                \r", i+1, dataLen)
		switch hasMP(data) {
		case false:
			pi, err := inputinfo.ParseFile(data)
			if err == nil {
				data += "  " + fdf.FMP(pi)
			}
		case true:
		}
		newData = append(newData, data)
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
	return strings.Contains(s, "  ")
}
