package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Configuration struct {
	IN_DIR          string `json:"IN"`
	BUFFER_DIR      string `json:"BUFFER"`
	IN_PROGRESS_DIR string `json:"IN_PROGRESS"`
	DONE_DIR        string `json:"DONE"`
	OUT_DIR         string `json:"OUT"`
}

var sep string = string(filepath.Separator)

func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("can't get userhome")
	}
	return fmt.Sprintf("%v%v.config%vaue%vdefault.config", home, sep, sep, sep)
}

func Default() *Configuration {
	cfg := Configuration{}
	cfg.BUFFER_DIR = `\\192.168.31.4\buffer\IN\`
	cfg.IN_DIR = cfg.BUFFER_DIR + `@AMEDIA_IN\`
	cfg.IN_PROGRESS_DIR = cfg.BUFFER_DIR + `_IN_PROGRESS\`
	cfg.DONE_DIR = cfg.BUFFER_DIR + `_DONE\`
	cfg.OUT_DIR = `\\nas\ROOT\EDIT\_amedia\_autogen\`
	return &cfg
}

func (cfg *Configuration) Save() error {
	bt, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config failed: %v", err)
	}
	dir := filepath.Dir(ConfigPath())
	os.MkdirAll(dir, 0777)
	f, err := os.Create(ConfigPath())
	if err != nil {
		return fmt.Errorf("create file failed: %v", err)
	}
	if _, err := f.Write(bt); err != nil {
		return fmt.Errorf("write config failed: %v", err)
	}
	f.Close()
	return nil
}

func Load() (*Configuration, error) {
	cfg := &Configuration{}
	bt, err := os.ReadFile(ConfigPath())
	if err != nil {
		return nil, fmt.Errorf("read config failed: %v", err)
	}
	if err := json.Unmarshal(bt, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %v", err)
	}
	return cfg, nil
}
