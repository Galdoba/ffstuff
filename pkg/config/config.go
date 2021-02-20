package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ConfigConstructor struct {
}

/*
config.Construct(
	config.SetField(key, val string),
	config.SetField(key2, val2 string),
)

*/

type Field func(string, string)

func Construct(fields ...Field) {
	confDir, file := configPath()
	os.MkdirAll(confDir, os.ModePerm)
	dir, file := configPath()
	f, err := os.OpenFile(dir+"\\"+file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString("This is Text"); err != nil {
		panic(err)
	}
}

func configPath() (string, string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	exe, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	exe = filepath.Base(exe)
	configDir := ""
	switch runtime.GOOS {
	case "windows":
		exe = strings.TrimSuffix(exe, ".exe")
		configDir = home + "\\config\\" + exe + "\\" // + exe + ".config"

	}
	return configDir, exe + ".config"
}

//Read - reads config file for this specific program and returns [string]string map
func Read() (configMap map[string]string, err error) {
	confDir, confFile := configPath()
	f, err := os.OpenFile(confDir+"\\"+confFile, os.O_RDONLY, 0600)
	if err != nil {
		if strings.Contains(err.Error(), "The system cannot find the path specified") {
			Construct()
		}

	}
	defer f.Close()
	//TODO: проверить есть ли файл
	// return Err.Ошибки доступа и наличия файла
	//
	//configMap = make(map[string]string)
	//
	//TODO: собрать строки из файла
	// return Err.Ошибки парсинга
	//
	return configMap, nil
}

//Load - search and return config value by key
func Load(key string) (val string, err error) {
	switch runtime.GOOS {
	default:
		return "", errors.New("Unknown OS")
	case "windows":
		Construct()
		return loadFromWindows(key)
	}
}

func loadFromWindows(key string) (val string, err error) {
	fmt.Println(os.UserHomeDir())
	return "", nil
}
