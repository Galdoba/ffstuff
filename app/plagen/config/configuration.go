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
	AppName       string            `yaml:"App Name       ,omitempty"`
	Location      string            `yaml:"Location       ,omitempty"`
	VideoFormats  []string          `yaml:"Video Formats  ,omitempty"`
	VideoPaths    map[string]string `yaml:"Video Paths    ,omitempty"`
	AudioPaths    map[string]string `yaml:"Audio Paths    ,omitempty"`
	Subtitle      string            `yaml:"Subtitle Path  ,omitempty"`
	Destination   string            `yaml:"Destination    ,omitempty"`
	LanguagePairs map[string]string `yaml:"Language Pairs ,omitempty"`
	header        string
}

func NewConfig(program string) (*Config, error) {
	cfg := Config{}
	cfg.AppName = program
	loc := stdPath(program)
	cfg.Location = loc
	cfg.VideoPaths = make(map[string]string)
	cfg.AudioPaths = make(map[string]string)
	cfg.LanguagePairs = make(map[string]string)
	if err := os.MkdirAll(gpath.StdConfigDir(program), 0777); err != nil {
		return nil, err
	}
	err := cfg.setDefaultValues()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) setDefaultValues() error {
	cfg.VideoFormats = []string{"4K", "HD", "SD169", "SD43"}
	cfg.VideoPaths["alcohol"] = `\\192.168.31.4\edit\_exchange\_age\originals\alcohol_4k_7.15.mov`
	cfg.VideoPaths["smoking"] = `\\192.168.31.4\edit\_exchange\_age\originals\smoking_4k_7.15.mov`
	cfg.AudioPaths["silence20"] = `\\192.168.31.4\edit\_exchange\_age\originals\silence_20.ac3`
	cfg.AudioPaths["silence51"] = `\\192.168.31.4\edit\_exchange\_age\originals\silence_51.ac3`
	cfg.Subtitle = `\\192.168.31.4\edit\_exchange\_age\originals\silence.srt`
	cfg.Destination = `\\192.168.31.4\edit\_exchange\_age\generated\`
	cfg.LanguagePairs["rus"] = `r`
	cfg.LanguagePairs["eng"] = `e`
	cfg.LanguagePairs["ita"] = `ita`
	cfg.LanguagePairs["heb"] = `heb`
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
