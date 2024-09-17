package stdpath

import "fmt"

const (
	Dir    = "dir:"
	File   = "file:"
	LOG    = "log"
	STATS  = "stats"
	ASSETS = "assets"
	DATA   = "data"
	TMP    = "tmp"

//log - логирование (конструируется при рантайме)
//config - конфиг для работы программы (конструируется при компеляции, используется при рантайме)
//assets - файлы необходимые для работы программы (конструируются при компеляции, используется при рантайме)
//data - файлы необходимые для работы программы (используется при рантайме)
//tmp - временные файлы (конструируются при компеляции, используется при рантайме, удаляются при старте)
)

var app = "stdpath"

type path struct {
	isDir    bool
	appName  string
	tag      string
	layers   []string
	filename string
	fileext  string
	err      error
}

type pathData struct {
	isDir    bool
	appName  string
	tag      string
	layers   []string
	filename string
	fileext  string
}

type PathOption func(*pathData)

func defaultFile() pathData {
	return pathData{
		isDir:    false,
		appName:  app,
		tag:      "",
		layers:   []string{},
		filename: "",
		fileext:  "",
	}
}

func NewFile(opts ...PathOption) *path {
	p := path{}
	p.isDir = false
	settings := defaultData()
	for _, enrich := range opts {
		enrich(&settings)
	}
	p.isDir = settings.isDir
	p.appName = settings.appName
	p.tag = settings.tag
	p.layers = settings.layers
	p.filename = settings.filename
	if p.isDir {
		p.err = fmt.Errorf("path is directory (expected file)")
	}
	if p.filename == "" {
		p.err = fmt.Errorf("file has no name")
	}
	if p.fileext == "" {
		p.err = fmt.Errorf("file has no extention")
	}
	return &p
}

func NewDir(opts ...PathOption) *path {
	p := path{}
	p.isDir = true
	settings := defaultData()
	for _, enrich := range opts {
		enrich(&settings)
	}
	p.isDir = settings.isDir
	p.appName = settings.appName
	p.tag = settings.tag
	p.layers = settings.layers
	p.filename = settings.filename
	if p.isDir {
		p.err = fmt.Errorf("path is directory (expected file)")
	}
	if p.filename == "" {
		p.err = fmt.Errorf("file has no name")
	}
	if p.fileext == "" {
		p.err = fmt.Errorf("file has no extention")
	}
	return &p
}

/*
stdpath.SetAppName(app.Name)
stdpath.Create(
	stdpath.FileLog(),
	stdpath.FileConfig(),
	stdpath.DirAssets("career"),
	stdpath.DirAssets("characteristics"),
	stdpath.DirAssets("events"),
	stdpath.DirAssets("skills"),
	stdpath.DirData("presets"),
)
*/

func SetAppName(appName string) {
	app = appName
}

// func Create(paths ...path) error {
// 	if app == "stdpath" {
// 		return fmt.Errorf("app name was not set: use stdpath.SetAppName(string)")
// 	}
// 	for _, path := range paths {

// 	}
// 	dir := path("")
// 	file := path("")
// 	switch tag {
// 	case Dir + LOG:
// 		dir = newLogDir(app)
// 	}
// 	err := os.MkdirAll(dir, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	if file != "" && notExist(file) {
// 		f, err := os.Create(file)
// 		defer f.Close()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// var sep = string(filepath.Separator)

// func homeDir() string {
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return home + sep
// }

// func newLogDir(app string) path {
// 	return path(fmt.Sprintf("%v.log%v%v%v", homeDir(), sep, app, sep))
// }

// func newLogFile(app string, keys ...string) path {
// 	k := "default"
// 	for _, key := range keys {
// 		k = key
// 	}
// 	return path(fmt.Sprintf("%v.log%v%v%v%v.log", homeDir(), sep, app, sep, k))
// }

// func newConfigDir(app string) path {
// 	return path(fmt.Sprintf("%v.config%v%v%v", homeDir(), sep, app, sep))
// }

// func newConfigFile(app string, keys ...string) path {
// 	k := "default"
// 	for _, key := range keys {
// 		k = key
// 	}
// 	return path(fmt.Sprintf("%v.config%v%v%v%v.config", homeDir(), sep, app, sep, k))
// }

// func newProgramDir(app string) path {
// 	return path(fmt.Sprintf("%v%v%v%v%v%v%v", homeDir(), "Programs", sep, "galdoba", sep, app, sep))
// }

// func newAssetDir(app string, layers ...string) path {
// 	dir := ""
// 	for _, layer := range layers {
// 		dir += layer + sep
// 	}
// 	prefix := fmt.Sprintf("%v%v%v", newProgramDir(app), "assets", sep)
// 	return path(fmt.Sprintf("%v%v", prefix, dir))
// }

// func newDataDir(app string, layers ...string) path {
// 	dir := ""
// 	for _, layer := range layers {
// 		dir += layer + sep
// 	}
// 	prefix := fmt.Sprintf("%v%v%v", newProgramDir(app), "data", sep)
// 	return path(fmt.Sprintf("%v%v", prefix, dir))
// }

// func newTempDir(app string) path {
// 	return path(fmt.Sprintf("%v%v%v", newProgramDir(app), "tmp", sep))
// }

// func notExist(path string) bool {
// 	if _, err := os.Stat("/path/to/whatever"); err == nil {
// 		return false
// 	} else if errors.Is(err, os.ErrNotExist) {
// 		return true
// 	} else {
// 		panic(err.Error())
// 	}
// }
