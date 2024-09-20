package config

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"gopkg.in/yaml.v3"
)

const (
	DESTINATION_ROOT_PATH = "Direct"
	SOURCE_ROOT_PATH      = "Bash"
)

type Configuration struct {
	Version               string `yaml:"version"`
	SOURCE_ROOT_PATH      string `yaml:"Directory: Source Root"`
	DESTINATION_ROOT_PATH string `yaml:"Directory: Destination Root"`
	LOG                   string `yaml:"File     : Log         "`
}

func init() {
	stdpath.SetAppName("archivator")
}

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
		return fmt.Errorf("config marshaling failed: %v", err)
	}
	path := stdpath.ConfigFile(confKey)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%v config file creation failed: %v", confKey, err)
	}
	_, err = f.Write(bt)
	if err != nil {
		return fmt.Errorf("%v config file writing failed: %v", confKey, err)
	}
	return nil
}

func NewConfig(version string) *Configuration {
	cfg := Configuration{}
	cfg.Version = version
	cfg.SOURCE_ROOT_PATH = ""
	cfg.DESTINATION_ROOT_PATH = ""
	cfg.LOG = ""
	return &cfg
}
