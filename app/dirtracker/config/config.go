package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/Galdoba/devtools/configmanager"
)

const (
	appName = "dirtracker"
)

type configuration struct {
	app        string
	path       string
	isCustom   bool
	Root       string   `yaml:"Search Root"`            //
	WhLst      []string `yaml:"White List"`             //
	BlLst      []string `yaml:"Black List"`             //
	BlLstEnbl  bool     `yaml:"Black List Enabled"`     //
	WhLstEnbl  bool     `yaml:"White List Enabled"`     //
	UpdtCycl   int      `yaml:"Update Cycle"`           //seconds [default : 15]
	MaxThreads int      `yaml:"Maximum Search Threads"` // [default : 4]
}

type ConfigFile interface {
	Save() error
	SaveAs(string) error
	SetDefault() error
	Path() string
	AppName() string
	IsCustom() bool
}

type Config interface {
	ConfigFile
	SearchRoot() string
	WhiteList() []string
	BlackList() []string
	BlackListEnabled() bool
	WhiteListEnabled() bool
	UpdateCycle() int
	MaximumSearchThreads() int
}

// //////////NEW-SAVE-LOAD////////////
// New - autogenerated constructor of config file
func New() Config {
	cfg := configuration{}
	cfg.path = configmanager.DefaultConfigDir(appName) + "config.yaml"
	cfg.app = appName
	return &cfg
}

// Save - autogenerated constructor of config file
func (cfg *configuration) Save() error {
	data := []byte(header())
	bt, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("can't marshal config: %v", err.Error())
	}
	if err := os.MkdirAll(filepath.Dir(cfg.path), 0777); err != nil {
		return fmt.Errorf("can't create directory")
	}
	f, err := os.OpenFile(cfg.path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return fmt.Errorf("can't open file: %v", err.Error())
	}
	defer f.Close()
	f.Truncate(0)
	data = append(data, bt...)
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("can't write file: %v", err.Error())
	}
	return nil
}

// SaveAs - autogenerated saver of alternative config file
func (cfg *configuration) SaveAs(path string) error {
	cfg.path = path
	cfg.isCustom = true
	return cfg.Save()
}

// Load - Load default config
func Load() (Config, error) {
	path := stdConfigPath()
	cfg, err := loadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("can't load default config: %v", err.Error())
	}
	cfg.isCustom = true
	return cfg, nil
}

// LoadCustom - Loader custom config
func LoadCustom(path string) (Config, error) {
	cfg, err := loadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("can't load custom config: %v", err.Error())
	}
	cfg.isCustom = true
	return cfg, nil
}

// loadConfig - autogenerated loader config file
func loadConfig(path string) (*configuration, error) {
	bt, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%v", err.Error())
	}
	cfg := &configuration{}
	err = yaml.Unmarshal(bt, cfg)
	if err != nil {
		return nil, fmt.Errorf("%v", err.Error())
	}
	return cfg, nil
}

////////////HELPERS////////////

// Path - return filepath of current config
func (cfg *configuration) Path() string {
	return cfg.path
}

// IsCustom - return true if config is custom
func (cfg *configuration) IsCustom() bool {
	return cfg.isCustom
}

// AppName - return true if config is custom
func (cfg *configuration) AppName() string {
	return cfg.app
}

func stdConfigDir() string {
	path, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	sep := string(filepath.Separator)
	path += sep
	return path + ".config" + sep + appName + sep
}

func stdConfigPath() string {
	return stdConfigDir() + "config.yaml"
}

////////////GETTERS////////////

// SearchRoot - autogenerated getter for 'Search Root' option
func (cfg *configuration) SearchRoot() string {
	return cfg.Root
}

// WhiteList - autogenerated getter for 'White List' option
func (cfg *configuration) WhiteList() []string {
	return cfg.WhLst
}

// BlackList - autogenerated getter for 'Black List' option
func (cfg *configuration) BlackList() []string {
	return cfg.BlLst
}

// BlackListEnabled - autogenerated getter for 'Black List Enabled' option
func (cfg *configuration) BlackListEnabled() bool {
	return cfg.BlLstEnbl
}

// WhiteListEnabled - autogenerated getter for 'White List Enabled' option
func (cfg *configuration) WhiteListEnabled() bool {
	return cfg.WhLstEnbl
}

// UpdateCycle - autogenerated getter for 'Update Cycle' option
// seconds
func (cfg *configuration) UpdateCycle() int {
	return cfg.UpdtCycl
}

// MaximumSearchThreads - autogenerated getter for 'Maximum Search Threads' option
func (cfg *configuration) MaximumSearchThreads() int {
	return cfg.MaxThreads
}

////////////SETTERS////////////

// SetSearchRoot - autogenerated setter for 'Search Root' option
func (cfg *configuration) SetSearchRoot(root string) {
	cfg.Root = root
}

// SetWhiteList - autogenerated setter for 'White List' option
func (cfg *configuration) SetWhiteList(whlst []string) {
	cfg.WhLst = whlst
}

// SetBlackList - autogenerated setter for 'Black List' option
func (cfg *configuration) SetBlackList(bllst []string) {
	cfg.BlLst = bllst
}

// SetBlackListEnabled - autogenerated setter for 'Black List Enabled' option
func (cfg *configuration) SetBlackListEnabled(bllstenbl bool) {
	cfg.BlLstEnbl = bllstenbl
}

// SetWhiteListEnabled - autogenerated setter for 'White List Enabled' option
func (cfg *configuration) SetWhiteListEnabled(whlstenbl bool) {
	cfg.WhLstEnbl = whlstenbl
}

// SetUpdateCycle - autogenerated setter for 'Update Cycle' option
func (cfg *configuration) SetUpdateCycle(updtcycl int) {
	cfg.UpdtCycl = updtcycl
}

// SetMaximumSearchThreads - autogenerated setter for 'Maximum Search Threads' option
func (cfg *configuration) SetMaximumSearchThreads(maxthreads int) {
	cfg.MaxThreads = maxthreads
}

func (cfg *configuration) SetDefault() error {
	cfg.UpdtCycl = 15
	cfg.MaxThreads = 4
	return cfg.Save()
}

func header() string {
	hdr := ""
	hdr += `################################################################################` + "\n"
	hdr += `#                 This file was generated by configbuilder app                 #` + "\n"
	hdr += `#                     Check formatting rules before editing                    #` + "\n"
	hdr += `#                 https://docs.fileformat.com/programming/yaml/                #` + "\n"
	hdr += `################################################################################` + "\n"
	hdr += `# expected location: C:\Users\Admin\.config\dirtracker\config.yaml` + "\n"
	hdr += `# app name         : dirtracker` + "\n"
	hdr += `` + "\n"
	return hdr
}
