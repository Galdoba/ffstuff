package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/gpath"
)

func listPresets() ([]string, error) {
	fullList, err := os.ReadDir(gpath.StdPath("", []string{".config", programName, "presets"}...))
	if err != nil {
		return nil, err
	}
	list := []string{}
	for _, found := range fullList {
		if strings.HasPrefix(found.Name(), "preset_") && strings.HasSuffix(found.Name(), ".json") {
			list = append(list, found.Name())
		}
	}
	return list, nil
}

func loadPreset(path string) []columnData {

	return nil
}

func createDefaultPreset() error {
	path := presetDir
	fmt.Println(path)
	return nil
}
