package v2

import (
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	ProgramVersion       string            `yaml:"version"`
	ProgramDirectories   map[string]string `yaml:"Program Directories"`
	DirectProcessing     bool              `yaml:"Direct Processing Mode"`
	GenerateBash         bool              `yaml:"Bash Script Generation"`
	BashPathsTranslation map[string]string `yaml:"Bash Paths Translation"`
}

func programDirKeys() []string {
	return []string{
		"IN              ",
		"OUT             ",
		"Pre-Processing  ",
		"Processing      ",
		"Post-Processing ",
		"Stats           ",
	}
}

func NewConfig(version string) *Configuration {
	cfg := &Configuration{}
	cfg.ProgramVersion = version
	cfg.ProgramDirectories = make(map[string]string)
	for _, dir := range programDirKeys() {
		cfg.ProgramDirectories[dir] = ""
	}
	cfg.DirectProcessing = false
	cfg.GenerateBash = true
	cfg.BashPathsTranslation = make(map[string]string)
	for _, dir := range programDirKeys() {
		cfg.BashPathsTranslation[dir] = ""
	}
	return cfg
}

func example() *Configuration {
	return &Configuration{
		ProgramVersion:       "aaa",
		DirectProcessing:     true,
		GenerateBash:         false,
		BashPathsTranslation: map[string]string{"aaa": "AAA", "bbbb": "BBBBB"},
	}
}

func marshal(cfg *Configuration) ([]byte, error) {
	bt, err := yaml.Marshal(cfg)
	return bt, err
}

/*
Input Root Directory:
	'default' : '//192.168.31.55/path/dir/'
   	'   bash' : 'mnt/pemaltynov/path/dir/'
*/
