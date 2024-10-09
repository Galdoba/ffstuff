package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
	"github.com/Galdoba/ffstuff/app/grabber/internal/validation"
	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"gopkg.in/yaml.v3"
)

const (
	DESTINATION_ROOT_PATH  = "Direct"
	SOURCE_ROOT_PATH       = "Bash"
	SORT_BY_SIZE           = "SIZE"
	SORT_BY_PRIORITY       = "PRIORITY"
	SORT_BY_NONE           = "NONE"
	COPY_HANDLER_SKIP      = "SKIP"
	COPY_HANDLER_RENAME    = "RENAME"
	COPY_HANDLER_OVERWRITE = "OVERWRITE"
)

/*
оповещения (на перспективу):
tgChannel - в какой чат плевать текст
system    - выскакивающее уведомление на локальной машине

ОПЦИОНАЛЬНО:
Затирание Маркера
Плевать в чат
проверки место на диске
проверка возраста файлов
*/
type Configuration struct {
	Version               string            `yaml:"version"`
	MARKER_FILE_EXTENTION string            `yaml:"Marker File Extention"`
	DEFAULT_DESTINATION   string            `yaml:"Default Destination Directory"`
	SEARCH_ROOTS          []string          `yaml:"Directories: Search Markers (track mode only)"`
	LOG                   string            `yaml:"Log File"`
	LOG_LEVEL             string            `yaml:"Minimum Log Level"`
	LOG_BY_SESSION        bool              `yaml:"Create Log File for Every Session,omitempty"`
	CRON_TRIGGERS         map[string]string `yaml:"Cron Triggers"`
	COPY_MARKER           string            `yaml:"New Copy Suffix Mask,omitempty"`
	COPY_MARKER_comment   string            `yaml:"Copy Marker Comment,omitempty"`
	COPY_PREFIX           bool              `yaml:"Use Prefix instead of Suffix for Renaming"`
	COPY_HANDLING         string            `yaml:"Existing Copy Handling"`
	DELETE_ORIGINAL       string            `yaml:"Delete Original Files"`
	SORT_METHOD           string            `yaml:"Default Sorting Method"`

	FILE_PRIORITY_WEIGHTS      map[string]int `yaml:"File Priority Weights"`
	DIRECTORY_PRIORITY_WEIGHTS map[string]int `yaml:"Directory Priority Weights"`
}

var ErrNoConfig = errors.New("no config found")

func Load(key ...string) (*Configuration, error) {
	confKey := "default"
	for _, k := range key {
		confKey = k
	}
	path := stdpath.ConfigFile(confKey)
	bt, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoConfig
		}
		return nil, fmt.Errorf("can't read config file: %v", err)
	}
	cfg := &Configuration{}
	err = yaml.Unmarshal(bt, cfg)
	if err != nil {
		return nil, fmt.Errorf("config file unmarshaling failed: %v", err)
	}

	return cfg, nil
}

func Save(cfg *Configuration, key ...string) error {
	confKey := "default"
	for _, k := range key {
		confKey = k
	}
	bt, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("config marshaling failed: %v", err)
	}
	path := stdpath.ConfigFile(confKey)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%v config file creation failed: %v", confKey, err)
	}
	_, err = f.Write(bt)
	if err != nil {
		return fmt.Errorf("%v config file writing failed: %v", confKey, err)
	}
	return nil
}

func NewConfig(version string) *Configuration {
	cfg := Configuration{}
	cfg.Version = version
	cfg.MARKER_FILE_EXTENTION = ".ready"
	cfg.SORT_METHOD = SORT_BY_NONE
	cfg.DEFAULT_DESTINATION = ""
	cfg.LOG = ""
	cfg.LOG_LEVEL = "DEBUG"
	cfg.COPY_HANDLING = grabberflag.VALUE_COPY_SKIP
	cfg.DELETE_ORIGINAL = grabberflag.VALUE_DELETE_MARKER
	cfg.SORT_METHOD = grabberflag.VALUE_SORT_PRIORITY
	cfg.COPY_MARKER = "_copy_([C])"
	cfg.COPY_MARKER = "valid tags: [C] - Counter; TO BE EXPECTED: [T]; [D]; [S]"
	cfg.CRON_TRIGGERS = make(map[string]string)
	cfg.CRON_TRIGGERS["* * * * *"] = "test every minute"
	cfg.FILE_PRIORITY_WEIGHTS = make(map[string]int)
	cfg.FILE_PRIORITY_WEIGHTS[".ready"] = 100
	cfg.FILE_PRIORITY_WEIGHTS[".srt"] = 100
	cfg.FILE_PRIORITY_WEIGHTS[".aac"] = 10
	cfg.FILE_PRIORITY_WEIGHTS[".ac3"] = 10
	cfg.FILE_PRIORITY_WEIGHTS[".m4a"] = 10
	cfg.FILE_PRIORITY_WEIGHTS[".mp4"] = 5
	cfg.FILE_PRIORITY_WEIGHTS[".mov"] = 5
	cfg.FILE_PRIORITY_WEIGHTS[".mpg"] = 5
	cfg.FILE_PRIORITY_WEIGHTS["_proxy"] = 5
	cfg.FILE_PRIORITY_WEIGHTS["_4k"] = 1
	cfg.FILE_PRIORITY_WEIGHTS["_hd"] = 3
	cfg.FILE_PRIORITY_WEIGHTS["_sd"] = 5
	cfg.FILE_PRIORITY_WEIGHTS["_audio"] = 20
	cfg.DIRECTORY_PRIORITY_WEIGHTS = make(map[string]int)
	cfg.DIRECTORY_PRIORITY_WEIGHTS["amedia"] = 20

	return &cfg
}

func Validate(cfg *Configuration) []error {
	errors := []error{}
	switch cfg.LOG {
	case "":
		errors = append(errors, fmt.Errorf("log filepath is not set"))
	default:
		if err := validation.FileValidation(cfg.LOG); err != nil {
			errors = append(errors, fmt.Errorf("log filepath: %v", err))
		}
	}
	logLevelExpect := "\nexpecting: TRACE, DEBUG, INFO, WARN, ERROR or FATAL"
	switch cfg.LOG_LEVEL {
	case "":
		errors = append(errors, fmt.Errorf("grabber log level is not set %v", logLevelExpect))
	default:
		errors = append(errors, fmt.Errorf("grabber log level is unknown: '%v'%v", cfg.LOG_LEVEL, logLevelExpect))
	case "TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL":
	}
	switch cfg.DEFAULT_DESTINATION {
	case "":
		errors = append(errors, fmt.Errorf("grabber default destination is not set"))
	default:
		if err := validation.DirectoryValidation(cfg.DEFAULT_DESTINATION); err != nil {
			errors = append(errors, fmt.Errorf("grabber default destination: %v", err))
		}
	}

	if cfg.COPY_MARKER == "" {
		errors = append(errors, fmt.Errorf("grabber New Copy Mask is not set"))
	}
	if cfg.FILE_PRIORITY_WEIGHTS == nil {
		errors = append(errors, fmt.Errorf("grabber File Priority Weights are not set"))
	}
	if cfg.DIRECTORY_PRIORITY_WEIGHTS == nil {
		errors = append(errors, fmt.Errorf("grabber Directory Priority Weights are not set"))
	}

	switch cfg.SORT_METHOD {
	case grabberflag.VALUE_SORT_PRIORITY, grabberflag.VALUE_SORT_SIZE, grabberflag.VALUE_SORT_NONE:
	default:
		errors = append(errors, fmt.Errorf("grabber default sorting method is invalid (expect string '%v', '%v' or '%v')", grabberflag.VALUE_SORT_PRIORITY, grabberflag.VALUE_SORT_SIZE, grabberflag.VALUE_SORT_NONE))
	}

	return errors
}

func testShedule(shed string) error {
	if shed == "* * * * *" {
		return nil
	}
	return nil
}
