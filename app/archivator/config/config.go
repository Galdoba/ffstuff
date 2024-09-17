package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	ARCHIVE_ROOT_PATH        = "Direct"
	PROCESSING_MODE_BASH     = "Bash"
	Undefined                = "UNDEFINED"
	Asset_File_Log           = "Asset_File_Log"
	Asset_File_Serial_data   = "Asset_File_Serial_data"
	Asset_File_Movie_data    = "Asset_File_Movie_data"
	Asset_File_Stats_Global  = "Asset_File_Stats_Global"
	Asset_File_Stats_Session = "Asset_File_Stats_Session"
)

type Configuration struct {
	Version           string `yaml:"version"`
	ARCHIVE_ROOT_PATH string `yaml:"Directory: Archive Root"`
	LOG               string `yaml:"File     : Log         "`
	DB                string `yaml:"File     : Index File  "`
}

var sep string = string(filepath.Separator)

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("can't get userhome")
	}
	return fmt.Sprintf("%v%v.config%varchivator%v", home, sep, sep, sep)
}

func Load(key ...string) (*Configuration, error) {
	confKey := "default"
	for _, k := range key {
		confKey = k
	}
	path := configDir() + confKey + ".config"
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
	f, err := os.Create(configDir() + confKey + ".config")
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
	cfg.ARCHIVE_ROOT_PATH = Undefined
	cfg.LOG = Undefined
	cfg.DB = Undefined
	return &cfg
}
