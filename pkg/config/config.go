package config

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

//Load - search and return config value by key
func Load(key string) (val string, err error) {
	switch runtime.GOOS {
	default:
		return "", errors.New("Unknown OS")
	case "windows":
		return loadFromWindows(key)
	}
}

func loadFromWindows(key string) (val string, err error) {
	fmt.Println(os.UserHomeDir())
	return "", nil
}
