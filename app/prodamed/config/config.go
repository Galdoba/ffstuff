package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type OptFunc func(*Options)

func defaultOptions() Options {
	return Options{
		PATH:         make(map[string]string),
		CycleSeconds: 30,
	}
}

var sep = string(filepath.Separator)

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("user home: %v", err.Error()))
	}
	configPath := fmt.Sprintf("%v%v", home, sep, ".galdoba", sep)
	return home + sep
}

type Config struct {
	path   string
	Option Options
}

type Options struct {
	PATH         map[string]string `json:"Paths"`
	CycleSeconds int               `json:"Repeat Cycle (Seconds)"`
}

func WithPath(key, path string) OptFunc {
	return func(opt *Options) {
		opt.PATH[key] = path
	}
}

func Load(path string) (*Config, error) {
	if path == "" {
		path = defaultConfigPath()
	}

}
