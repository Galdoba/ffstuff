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

/*
оповещения (на перспективу):
tgChannel - в какой чат плевать текст
system    - выскакивающее уведомление на локальной машине
*/
type Configuration struct {
	Version             string         `yaml:"version"`
	DEFAULT_DESTINATION string         `yaml:"Directory  : Default Destination "`
	TASK_DIR            string         `yaml:"Directory  : Queue Storage       "`
	SEARCH_ROOTS        []string       `yaml:"Directories: Search Markers In   "`
	LOG                 string         `yaml:"File       : Log                 "`
	LOG_LEVEL           string         `yaml:"Minimum Log Level"`
	LOG_BY_SESSION      bool           `yaml:"Log By Session"`
	TRIGGER_BY_SCHEDULE bool           `yaml:"Schedule Trigger"`
	SCHEDULE            string         `yaml:"Schedule"`
	TRIGGER_BY_TIMEOUT  bool           `yaml:"Timeout Trigger"`
	TIMEOUT             int            `yaml:"Timeout (Seconds)"`
	PRIORITY_MAP        map[string]int `yaml:"Processing Priority"`
	GRAB_BY_SIZE        bool           `yaml:"Process Small Files First (Ignore Priority)"`
	COPY_PREFIX         string         `yaml:"New Copy Prefix Mask"`
	COPY_SUFFIX         string         `yaml:"New Copy Suffix Mask"`
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
	//cfg.SOURCE_ROOT_PATH = ""
	cfg.DEFAULT_DESTINATION = ""
	cfg.LOG = ""
	return &cfg
}
