package gconfig

import (
	"fmt"
	"os"
	"path/filepath"
)

func DefineConfigPath(programName string) string {
	userdir, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	sep := string(filepath.Separator)

	return fmt.Sprintf("%v%v.config%v%v%v%v.json", userdir, sep, sep, programName, sep, programName)
}

func DefineProgramDirectory(programName string) string {
	userdir, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	sep := string(filepath.Separator)
	return fmt.Sprintf("%v%v%v%v", userdir, sep, programName, sep)
}

/*
var programConfig *config

type config struct {
	ReadMe           string   `json:"ReadMe",omitempty`
	ListSeparatorArg string   `json:"ListSeparatorArg"`
	Rules            []string `json:"Rules"`
	Verbose          bool     `json:"Verbose Mode"`
	LogFile          string   `json:"LogFile"`
	LogLenght        int      `json:"LogLenght"`
}

func defaultConfig() *config {
	cfg := config{}
	cfg.ListSeparatorArg = "======"
	cfg.Rules = []string{"Rule_1"}
	cfg.Verbose = true
	cfg.LogFile = ""
	return &cfg
}

func init() {
	programName = "glisort"
	configPath = defineConfigPath()
	data, err := os.ReadFile(configPath)
	fmt.Println(configPath)

	if err != nil {
		switch {
		default:
			fmt.Println("Неизвестная ошибка при проверки наличия конфига:")
			println(err.Error())
			panic(0)
		case strings.Contains(err.Error(), "The system cannot find the file specified"), strings.Contains(err.Error(), "The system cannot find the path specified"):
			fmt.Println("Config file not found")
			err := os.MkdirAll(strings.TrimSuffix(configPath, "glisort.json"), 0777)
			if err != nil {
				panic(err.Error())
			}
			programConfig = defaultConfig()
			data, err = json.MarshalIndent(programConfig, "", "  ")
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Creating default config at %v: ", configPath)
			f, err := os.Create(configPath)
			if err != nil {
				panic(err.Error())
			}
			defer f.Close()
			_, err = f.Write(data)
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("ok")
		}
	}

	err = json.Unmarshal(data, &programConfig)
	if err != nil {
		panic(err.Error())
	}
}

*/
