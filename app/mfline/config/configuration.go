package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/devtools/gpath"
	"gopkg.in/yaml.v3"
)

// config struct  
type Config struct {
	AppName       string
	Location      string
	header        string
	StorageDir    string   `yaml:"Scan Storage Directory,omitempty"`
	WriteLogs     bool     `yaml:"Write Logs,omitempty"`
	LogFile       string   `yaml:"Log File,omitempty"`
	OldScan       float64  `yaml:"Old Scan Age (hours),omitempty"`
	AutoDeleteOld bool     `yaml:"Delete Old Scans,omitempty"`
	AutoScan      bool     `yaml:"Scan All Files in Tracked Directories,omitempty"`
	RescanIfErr   bool     `yaml:"Repeat Scans if Error is met,omitempty"`
	TrackDirs     []string `yaml:"Track Directories,omitempty"`
}

func NewConfig(program string) (*Config, error) {
	cfg := Config{}
	cfg.AppName = program
	loc := stdPath(program)
	cfg.Location = loc
	if err := os.MkdirAll(gpath.StdConfigDir(program), 0777); err != nil {
		return nil, err
	}
	err := cfg.setDefaultValues()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(cfg.StorageDir, 0777); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) setDefaultValues() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	sep := string(filepath.Separator)
	//cfg.Header = header(cfg.AppName)
	cfg.StorageDir = home + sep + ".ffstuff" + sep + "data" + sep + cfg.AppName + sep
	cfg.LogFile = home + sep + ".ffstuff" + sep + "logs" + sep + cfg.AppName + ".log"
	cfg.WriteLogs = false
	cfg.OldScan = 72.0
	cfg.AutoDeleteOld = false
	cfg.AutoScan = false
	cfg.RescanIfErr = false
	cfg.TrackDirs = []string{}
	return nil
}

func (cfg *Config) Save() error {

	data := []byte(header(cfg.AppName))
	bt, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	dir := filepath.Dir(cfg.Location)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(cfg.Location, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Truncate(0)
	_, err = f.WriteString(header(cfg.AppName))
	// if cfg.header != "" {
	// 	_, err := f.WriteString(cfg.Header)
	if err != nil {
		return fmt.Errorf("can't save file: write header: %v", err.Error())
	}
	// }
	data = append(data, bt...)
	_, err = f.Write(bt)
	return err
}

func Load(program string) (*Config, error) {
	path := stdPath(program)
	bt, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't load config: %v", err.Error())
	}
	cfg := &Config{}
	err = yaml.Unmarshal(bt, cfg)
	if err != nil {
		return nil, fmt.Errorf("can't load config: %v", err.Error())
	}
	return cfg, nil
}

func stdPath(program string) string {
	path := gpath.StdConfigDir(program) + "config.yaml"
	return path
}

func header(program string) string {
	return strings.Join([]string{
		fmt.Sprintf("#######################################################"),
		fmt.Sprintf("#  This is auto generated config file.                #"),
		fmt.Sprintf("#  Check formatting rules before manual edit.         #"),
		fmt.Sprintf("#  https://docs.fileformat.com/programming/yaml/      #"),
		fmt.Sprintf("#######################################################\n"),
	}, "\n")
}

func (cfg *Config) String() string {
	data, err := os.ReadFile(cfg.Location)
	if err != nil {
		return fmt.Sprintf("%v", err.Error())
	}
	return string(data)
}
