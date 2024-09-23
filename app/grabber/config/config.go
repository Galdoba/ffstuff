package config

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/grabber/internal/validation"
	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"gopkg.in/yaml.v3"
)

const (
	DESTINATION_ROOT_PATH = "Direct"
	SOURCE_ROOT_PATH      = "Bash"
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
	Version             string         `yaml:"version"`
	DEFAULT_DESTINATION string         `yaml:"Default Destination Directory"`
	TASK_DIR            string         `yaml:"Queue Storage Directory,omitempty"`
	SEARCH_ROOTS        []string       `yaml:"Directories: Search Markers In,omitempty"`
	LOG                 string         `yaml:"Log File"`
	LOG_LEVEL           string         `yaml:"Minimum Log Level"`
	LOG_BY_SESSION      bool           `yaml:"Log By Session,omitempty"`
	TRIGGER_BY_SCHEDULE bool           `yaml:"Schedule Trigger,omitempty"`
	SCHEDULE            string         `yaml:"Schedule,omitempty"`
	TRIGGER_BY_TIMEOUT  bool           `yaml:"Timeout Trigger,omitempty"`
	TIMEOUT             int            `yaml:"Timeout (Seconds),omitempty"`
	PRIORITY_MAP        map[string]int `yaml:"Processing Priority"`
	GRAB_BY_SIZE        bool           `yaml:"Process Small Files First (Ignore Priority)"`
	COPY_PREFIX         string         `yaml:"New Copy Prefix Mask,omitempty"`
	COPY_SUFFIX         string         `yaml:"New Copy Suffix Mask,omitempty"`
	COPY_HANDLING       string         `yaml:"Existing Copy Handling"`
}

func Load(key ...string) (*Configuration, error) {
	confKey := "default"
	for _, k := range key {
		confKey = k
	}
	path := stdpath.ConfigFile(confKey)
	bt, err := os.ReadFile(path)
	if err != nil {
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
	cfg.DEFAULT_DESTINATION = ""
	cfg.LOG = ""
	cfg.LOG_LEVEL = "DEBUG"
	cfg.COPY_SUFFIX = "copy_[C]"
	cfg.COPY_HANDLING = "SKIP"
	cfg.SCHEDULE = "0 6 * * 1,2,3,4,5"
	cfg.PRIORITY_MAP = make(map[string]int)
	cfg.PRIORITY_MAP[".ready"] = 100
	cfg.PRIORITY_MAP[".srt"] = 100
	cfg.PRIORITY_MAP[".aac"] = 10
	cfg.PRIORITY_MAP[".ac3"] = 10
	cfg.PRIORITY_MAP[".m4a"] = 10
	cfg.PRIORITY_MAP[".mp4"] = 5
	cfg.PRIORITY_MAP[".mov"] = 5
	cfg.PRIORITY_MAP[".mpg"] = 5
	cfg.PRIORITY_MAP["_proxy"] = 5
	cfg.PRIORITY_MAP["_4k"] = 1
	cfg.PRIORITY_MAP["_hd"] = 3
	cfg.PRIORITY_MAP["_sd"] = 5
	cfg.PRIORITY_MAP["_audio"] = 20

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
	switch cfg.SCHEDULE {
	case "":
		errors = append(errors, fmt.Errorf("grabber shedule trigger is not set"))
	default:
		//errors = append(errors, fmt.Errorf("grabber shedule trigger: %v", testShedule(cfg.SCHEDULE)))
	}
	switch cfg.COPY_HANDLING {
	case "SKIP", "OVERWRITE", "RENAME":
	case "":
		errors = append(errors, fmt.Errorf("grabber existient copy handling is not set (expect SKIP, OVERWRITE or RENAME)"))
	default:
		errors = append(errors, fmt.Errorf("grabber existient copy handling is invalid (expect SKIP, OVERWRITE or RENAME)"))
	}
	if cfg.COPY_PREFIX+cfg.COPY_SUFFIX == "" {
		errors = append(errors, fmt.Errorf("grabber New Copy Mask is not set"))
	}
	if cfg.PRIORITY_MAP == nil {
		errors = append(errors, fmt.Errorf("grabber Priority Map is not set"))
	}
	if cfg.TIMEOUT < 0 {
		errors = append(errors, fmt.Errorf("grabber Timeout Trigger is invalid (expect int >= 0)"))
	}
	return errors
}

func testShedule(shed string) error {
	if shed == "* * * * *" {
		return nil
	}
	return nil
}
