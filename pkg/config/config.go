package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Galdoba/utils"
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
	fmt.Println("Start Read")
	confDir, confFile := configPath()
	f, err := os.OpenFile(confDir+"\\"+confFile, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println(err.Error())
		errStr := err.Error()
		if strings.Contains(errStr, "cannot find the file") {
			fmt.Println("Do construct")
			Construct()
		}

	}
	defer f.Close()
	keyVal := make(map[string]string)
	for _, ln := range utils.LinesFromTXT(f.Name()) {
		kv := strings.Split(ln, ":=")
		if len(kv) == 2 {
			keyVal[kv[0]] = kv[1]
		}
	}

	//TODO: проверить есть ли файл
	// return Err.Ошибки доступа и наличия файла
	//
	//configMap = make(map[string]string)
	//
	//TODO: собрать строки из файла
	// return Err.Ошибки парсинга
	//
	fmt.Println("End Read")
	fmt.Println(keyVal)
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

func SetField(key, val string) {
	confDir, file := configPath()
	f, err := os.OpenFile(confDir+"\\"+file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	key = strings.ToUpper(key)
	key = strings.ReplaceAll(key, " ", "_")
	key = strings.TrimSuffix(key, "_")
	val = strings.ToUpper(val)
	val = strings.ReplaceAll(val, " ", "_")
	val = strings.TrimSuffix(val, "_")
	lines := utils.LinesFromTXT(confDir + "\\" + file)
	for n, line := range lines {
		if strings.Contains(line, key+":=") {
			utils.EditLineInFile(confDir+"\\"+file, n, key+":="+val)
			return
		}
	}
	if _, err = f.WriteString(key + ":=" + val); err != nil {
		panic(err)
	}
}
