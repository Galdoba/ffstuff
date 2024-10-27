package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"gopkg.in/yaml.v3"
)

const ()

/*
оповещения (на перспективу):
tgChannel - в какой чат плевать текст
system    - выскакивающее уведомление на локальной машине

ОПЦИОНАЛЬНО:
Затирание Маркера
Плевать в чат
проверки место на диске
проверка возраста файлов
*/
type Configuration struct {
	Version           string   `yaml:"version"`
	LOGS              []string `yaml:"Log Files"`
	CONSOLE_LOG_LEVEL string   `yaml:"Minimum Log Level: Terminal"`
	FILE_LOG_LEVEL    string   `yaml:"Minimum Log Level: File    "`
	LOG_BY_SESSION    bool     `yaml:"Create Log File for Every Session,omitempty"`
}

var ErrNoConfig = errors.New("no config found")

func Load(key ...string) (*Configuration, error) {
	confKey := "default"
	for _, k := range key {
		confKey = k
	}
	path := stdpath.ConfigFile(confKey)
	bt, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoConfig
		}
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
	os.MkdirAll(filepath.Dir(path), 0666)
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
	cfg.LOGS = append(cfg.LOGS, stdpath.LogFile())
	cfg.CONSOLE_LOG_LEVEL = "DEBUG"
	cfg.FILE_LOG_LEVEL = "DEBUG"
	return &cfg
}

func Validate(cfg *Configuration) []error {
	errors := []error{}

	return errors
}
