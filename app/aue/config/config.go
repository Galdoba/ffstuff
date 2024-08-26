package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	PROCESSING_MODE_DIRECT = "Direct"
	PROCESSING_MODE_BASH   = "Bash"
)

type Configuration struct {
	IN_DIR              string            `json:"Directory: IN         "`
	BUFFER_DIR          string            `json:"Directory: BUFFER     "`
	IN_PROGRESS_DIR     string            `json:"Directory: IN_PROGRESS"`
	DONE_DIR            string            `json:"Directory: DONE       "`
	OUT_DIR             string            `json:"Directory: OUT        "`
	DirectProcessing    bool              `json:"Direct Processing"`
	BashGeneration      bool              `json:"Generate Bash File"`
	BashPathTranslation map[string]string `json:"Bash Paths Translation"`
	SleepSeconds        int               `json:"Repeat Cycle (seconds)"`
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
	cfg.BUFFER_DIR = `//192.168.31.4/buffer/IN/`
	cfg.IN_DIR = cfg.BUFFER_DIR + `@AMEDIA_IN/`
	cfg.IN_PROGRESS_DIR = cfg.BUFFER_DIR + `_IN_PROGRESS/`
	cfg.DONE_DIR = cfg.BUFFER_DIR + `_DONE/`
	cfg.OUT_DIR = `//nas/ROOT/EDIT/_amedia/_autogen/`
	cfg.BashGeneration = true
	cfg.BashPathTranslation = make(map[string]string)
	cfg.BashPathTranslation[`//192.168.31.4/buffer/IN/`] = "/home/pemaltynov/IN/"
	cfg.BashPathTranslation[cfg.BUFFER_DIR+`@AMEDIA_IN/`] = "/home/pemaltynov/IN/@AMEDIA_IN/"
	cfg.BashPathTranslation[cfg.BUFFER_DIR+`_IN_PROGRESS/`] = "/home/pemaltynov/IN/_IN_PROGRESS/"
	cfg.BashPathTranslation[cfg.BUFFER_DIR+`_DONE/`] = "/home/pemaltynov/IN/_DONE/"
	cfg.BashPathTranslation[cfg.OUT_DIR] = "/mnt/pemaltynov/ROOT/EDIT/_amedia/_autogen/"
	cfg.SleepSeconds = 300
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
