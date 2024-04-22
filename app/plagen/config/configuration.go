package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// config struct  
type Config struct {
	AppName    string            `yaml:"App Name       ,omitempty"`
	Location   string            `yaml:"Location       ,omitempty"`
	VideoPaths map[string]string `yaml:"Video Paths    ,omitempty"`
	AudioPaths map[string]string `yaml:"Audio Paths    ,omitempty"`
	Subtitle   string            `yaml:"Subtitle Path  ,omitempty"`
	Storage    string            `yaml:"Storage Dir    ,omitempty"`
	header     string
}

func NewConfig(program string) (*Config, error) {
	cfg := Config{}
	cfg.AppName = program
	loc := plagenConfigPath()
	cfg.Location = loc
	cfg.VideoPaths = make(map[string]string)
	cfg.AudioPaths = make(map[string]string)
	// if err := os.MkdirAll(gpath.StdConfigDir(program), 0777); err != nil {
	// 	return nil, err
	// }
	err := cfg.setDefaultValues()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) setDefaultValues() error {
	sep := string(filepath.Separator)
	cfg.VideoPaths["alcohol"] = os.Getenv("AGELOGOPATH") + `originals` + sep + `alcohol__4k_7.15.mov`
	cfg.VideoPaths["msmoking"] = os.Getenv("AGELOGOPATH") + `originals` + sep + `msmoking__4k_7.15.mov`
	cfg.AudioPaths["2"] = os.Getenv("AGELOGOPATH") + `originals` + sep + `silence_20.ac3`
	cfg.AudioPaths["6"] = os.Getenv("AGELOGOPATH") + `originals` + sep + `silence_51.ac3`
	cfg.Subtitle = os.Getenv("AGELOGOPATH") + `originals` + sep + `silence.srt`
	cfg.Storage = os.Getenv("AGELOGOPATH") + `originals` + sep + `.cache` + sep
	return nil
}

func (cfg *Config) Save() error {

	data := []byte(header())
	bt, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	dir := plagenConfigPath()

	f, err := os.OpenFile(dir, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Truncate(0)
	_, err = f.WriteString(header())
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

func Load() (*Config, error) {
	path := plagenConfigPath()
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

func plagenConfigPath() string {
	path := os.Getenv("AGELOGOPATH") + "originals" + string(filepath.Separator) + "config.yaml"
	return path
}

func header() string {
	return strings.Join([]string{
		fmt.Sprintf("#######################################################"),
		fmt.Sprintf("#  This is auto generated config for plagen app       #"),
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
