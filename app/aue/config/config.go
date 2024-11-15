package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/ffstuff/pkg/logman/v1"
	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"gopkg.in/yaml.v3"
)

const (
	PROCESSING_MODE_DIRECT   = "Direct"
	PROCESSING_MODE_BASH     = "Bash"
	Undefined                = "UNDEFINED"
	Asset_File_Log           = "Asset_File_Log"
	Asset_File_Serial_data   = "Asset_File_Serial_data"
	Asset_File_Movie_data    = "Asset_File_Movie_data"
	Asset_File_Stats_Global  = "Asset_File_Stats_Global"
	Asset_File_Stats_Session = "Asset_File_Stats_Session"
)

type Configuration struct {
	Version             string            `yaml:"version"`
	IN_DIR              string            `yaml:"Directory: IN         "`
	BUFFER_DIR          string            `yaml:"Directory: BUFFER     "`
	IN_PROGRESS_DIR     string            `yaml:"Directory: IN_PROGRESS"`
	DONE_DIR            string            `yaml:"Directory: DONE       "`
	OUT_DIR             string            `yaml:"Directory: OUT        "`
	NotificationDir     string            `yaml:"Directory: NOTIFY_DIR "`
	DirectProcessing    bool              `yaml:"Direct Processing     "`
	DebugMode           bool              `yaml:"Debug Mode            "`
	BashGeneration      bool              `yaml:"Generate Bash File    "`
	AssetFiles          map[string]string `yaml:"Asset Files"`
	BashPathTranslation map[string]string `yaml:"Bash Paths Translation"`
	SleepSeconds        int               `yaml:"Repeat Cycle (seconds)"`
	CONSOLE_LOG_LEVEL   string            `yaml:"Log Level: Console"`
	FILE_LOG_LEVEL      string            `yaml:"Log Level: File"`
}

var sep string = string(filepath.Separator)

// func configDir() string {
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		panic("can't get userhome")
// 	}
// 	return fmt.Sprintf("%v%v.config%vaue%v", home, sep, sep, sep)
// }

func Load(key ...string) (*Configuration, error) {
	confKey := "default"
	for _, k := range key {
		confKey = k
	}
	path := stdpath.ConfigFile(confKey)
	bt, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't read config file: %v", err)
	}
	cfg := &Configuration{}
	err = yaml.Unmarshal(bt, cfg)
	if err != nil {
		return nil, fmt.Errorf("config file unmarshaling failed: %v", err)
	}
	return cfg, nil
}

func Save(cfg *Configuration, key ...string) error {
	confKey := "default"
	for _, k := range key {
		confKey = k
	}
	bt, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Errorf("config marshaling failed: %v", err)
	}
	path := stdpath.ConfigFile(confKey)
	fmt.Println(path)
	os.MkdirAll(filepath.Dir(path), 0666)
	f, err := os.Create(path)
	if err != nil {
		fmt.Errorf("%v config file creation failed: %v", confKey, err)
	}
	_, err = f.Write(bt)
	if err != nil {
		fmt.Errorf("%v config file writing failed: %v", confKey, err)
	}
	return nil
}

func NewConfig(version string) *Configuration {
	cfg := Configuration{}
	cfg.Version = version
	cfg.BUFFER_DIR = Undefined
	cfg.IN_DIR = Undefined
	cfg.IN_PROGRESS_DIR = Undefined
	cfg.DONE_DIR = Undefined
	cfg.OUT_DIR = Undefined
	cfg.NotificationDir = Undefined
	cfg.AssetFiles = make(map[string]string)
	cfg.AssetFiles[Asset_File_Log] = stdpath.LogFile() //configDir() + "aue.log"
	cfg.AssetFiles[Asset_File_Serial_data] = cfg.IN_DIR + "amedia_tv_series.xml"
	cfg.AssetFiles[Asset_File_Movie_data] = "???"
	cfg.AssetFiles[Asset_File_Stats_Global] = "???"
	cfg.AssetFiles[Asset_File_Stats_Session] = "???"
	cfg.DirectProcessing = false
	cfg.BashGeneration = true
	cfg.BashPathTranslation = make(map[string]string)
	cfg.BashPathTranslation[Undefined+"1"] = Undefined + "1_translated"
	cfg.BashPathTranslation[Undefined+"2"] = Undefined + "2_translated"
	cfg.CONSOLE_LOG_LEVEL = logman.DEBUG
	cfg.FILE_LOG_LEVEL = logman.DEBUG
	cfg.SleepSeconds = 300

	return &cfg
}
