package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var sep = string(filepath.Separator)
var loadAtemptNum int

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("user home: %v", err.Error()))
	}
	configPath := fmt.Sprintf("%v%v%v%v%v%v%v", home, sep, ".config", sep, "prodamed", sep, "default.prodamedconf")
	return configPath
}

type Config struct {
	path                  string
	isDefaultConfig       bool
	AlternativeConfigPath string  `json:"Alternative Config Path,omitempty"`
	Option                Options `json:"Options"`
}

func defaultConfig() *Config {
	cfg := Config{}
	cfg.isDefaultConfig = true
	cfg.path = ""
	cfg.Option = defaultOptions()
	return &cfg
}

func reportConfigErrors(cfg *Config) int {
	errors := 0
	for k, v := range cfg.Option.PATH {
		switch v {
		case "":
			errors++
			fmt.Printf("path variable '%v' was not setup\n", k)
		default:
			f, err := os.Stat(v)
			if err != nil {
				fmt.Printf("can't open '%v': %v\n", k, v)
				errors++
				continue
			}
			if !f.IsDir() {
				fmt.Printf("%v: '%v' is not a directory\n", k, v)
				errors++
			}
		}
	}
	if cfg.Option.CycleSeconds < 0 {
		fmt.Printf("cycle value is negative\n")
		errors++
	}
	return errors
}

// Save - Saves config to json file.
func (cfg *Config) Save() error {
	bt, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("config marshling failed: %v", err)
	}
	configDir := filepath.Dir(cfg.path)
	if err := os.MkdirAll(configDir, 0777); err != nil {
		return fmt.Errorf("directory creation failed: %v", err)
	}
	f, err := os.Create(cfg.path)
	if err != nil {
		return fmt.Errorf("file creation failed: %v", err)
	}
	defer f.Close()
	if _, err := f.Write(bt); err != nil {
		return fmt.Errorf("config writing failed: %v", err)
	}
	return nil
}

// Load - loads config from filepath provided.
// If filepath = "", default path will be used.
func Load(path string) (*Config, error) {
	loadAtemptNum++
	if loadAtemptNum > 20 {
		return nil, fmt.Errorf("load atempts limit reached")
	}
	if path == "" {
		path = defaultConfigPath()
	}
	bt, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("file reading failed: %v", err)
	}
	cfg := &Config{}
	if err := json.Unmarshal(bt, cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling failed: %v", err)
	}
	if cfg.AlternativeConfigPath != "" {
		cfgAlt, err := Load(cfg.AlternativeConfigPath)
		if err == nil {
			cfg = cfgAlt
		}
	}
	cfg.path = path
	loadAtemptNum = 0
	return cfg, nil
}
