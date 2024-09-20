package stdpath

import (
	"fmt"
	"os"
	"path/filepath"
)

var app string = "<NO_APP_NAME>"

func SetAppName(name string) {
	app = name
}

func sep() string {
	return string(filepath.Separator)
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	return home + sep()
}

func LogDir() string {
	return fmt.Sprintf("%v.log%vgaldoba%v%v%v", homeDir(), sep(), sep(), app, sep())
}

func LogFile(keys ...string) string {
	key := "default"
	for _, k := range keys {
		key = k
	}
	return fmt.Sprintf("%v%v.log", LogDir(), key)
}

func ConfigDir() string {
	return fmt.Sprintf("%v.config%vgaldoba%v%v%v", homeDir(), sep(), sep(), app, sep())
}

func ConfigFile(keys ...string) string {
	key := "default"
	for _, k := range keys {
		key = k
	}
	return fmt.Sprintf("%v%v.config", ConfigDir(), key)
}

func ProgramDir(layers ...string) string {
	suffix := ""
	for _, layer := range layers {
		suffix += layer + sep()
	}
	path := fmt.Sprintf("%vPrograms%vgaldoba%v%v%v%v", homeDir(), sep(), sep(), app, sep(), suffix)
	return path
}
