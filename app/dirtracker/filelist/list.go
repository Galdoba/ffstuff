package filelist

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/scanner"
	"gopkg.in/yaml.v3"
)

const (
	isDir  = "Directory"
	isFile = "File"
)

type FileList struct {
	paths    []fpath
	stats    map[string]int
	config   *config
	compiled string
}

type fpath struct {
	dir  string
	name string
}

type config struct {
	Root                string
	WhiteListEnabled    bool
	WhiteList           []string
	BlackListEnabled    bool
	BlackList           []string
	UpdateCycle_seconds int
}

type Config interface {
	GetRoot() string
	GetWhiteListEnabled() bool
	GetWhiteList() []string
	GetBlackListEnabled() bool
	GetBlackList() []string
	GetUpdateCycle_seconds() int
}

func (cf *config) GetRoot() string {
	return cf.Root
}
func (cf *config) GetWhiteListEnabled() bool {
	return cf.WhiteListEnabled
}
func (cf *config) GetWhiteList() []string {
	return cf.WhiteList
}
func (cf *config) GetBlackListEnabled() bool {
	return cf.BlackListEnabled
}
func (cf *config) GetBlackList() []string {
	return cf.BlackList
}
func (cf *config) GetUpdateCycle_seconds() int {
	return cf.UpdateCycle_seconds
}

func New(cfgData []byte) (*FileList, error) {
	fl := &FileList{}
	fl.stats = make(map[string]int)
	conf := &config{}
	err := yaml.Unmarshal(cfgData, conf)
	if err != nil {
		return fl, err
	}
	conf.normalizePaths()
	fl.config = conf
	if err := fl.Update(); err != nil {
		return fl, err
	}
	return fl, nil
}

func (fl *FileList) Update() error {
	fl.paths = []fpath{}
	l, err := scanner.Scan(fl.config.Root, "")
	if err != nil {
		return err
	}
	for _, pth := range l {
		if err := fl.AddEntry(pth); err != nil {
			return err
		}
	}
	fl.Compile()
	return nil
}

func (cfg *config) normalizePaths() {
	cfg.Root = strings.Replace(cfg.Root, "\\", "/", -1)
	for i := range cfg.WhiteList {
		cfg.WhiteList[i] = strings.Replace(cfg.WhiteList[i], "\\", "/", -1)
	}
	for i := range cfg.BlackList {
		cfg.BlackList[i] = strings.Replace(cfg.BlackList[i], "\\", "/", -1)
	}
}

func (fl *FileList) AddEntry(entry string) error {
	entry = strings.Replace(entry, "\\", "/", -1)
	prts := strings.Split(entry, "/")
	srtName := prts[len(prts)-1]
	dir := strings.TrimSuffix(entry, srtName)
	fp := fpath{dir, srtName}
	if err := fp.ensureDir(); err != nil {
		return err
	}
	switch fp.name {
	default:
		fl.stats[isFile]++
	case "":
		fl.stats[isDir]++
	}
	fl.paths = append(fl.paths, fp)
	return nil
}

func (fl *FileList) Stats() (int, int) {
	return fl.stats[isDir], fl.stats[isFile]
}
func (fl *FileList) NextUpdate() int {
	return fl.config.GetUpdateCycle_seconds()
}

//Compile - Выводит отдает список директорий/файлов исходя из следующих правил:
//если белый список более 0 - отдаем только то что попадает под белый список
//если черный список более 0 - отдаем только то что попадает не под черный список
func (fl *FileList) Compile() error {
	actual := []fpath{}
	if fl.config.WhiteListEnabled && fl.config.BlackListEnabled {
		return fmt.Errorf("can't have both White and Black lists enabled\ncheck file: {USER}/.config/dirtracker/dirtracker.config")
	}
	switch {
	case fl.config.WhiteListEnabled:
		for _, pth := range fl.config.WhiteList {
			actual = append(actual, fpath{pth, ""})
			for _, path := range fl.paths {
				if path.dir != pth {
					continue
				}
				actual = append(actual, path)
			}
		}
		//TODO: нужна логика поведения при включенном черном списке
		//case fl.config.BlackListEnabled:
	}
	if len(actual) == 0 {
		fl.compiled = ""
		return fmt.Errorf("Compiled empty list")
	}
	compiled := ""
	activeDir := ""
	activeDirectories := []string{}
	fileMaps := make(map[string][]string)
	for _, act := range actual {
		if activeDir != act.dir {
			activeDir = act.dir
			activeDirectories = append(activeDirectories, act.dir)
		}
		if act.name != "" {
			fileMaps[act.dir] = append(fileMaps[act.dir], act.name)
		}
	}
	for _, directory := range activeDirectories {
		compiled += fmt.Sprintf("%v\nfiles found %v:\n", directory, len(fileMaps[directory]))
		for _, file := range fileMaps[directory] {
			compiled += "  " + file + "\n"
		}
	}
	fl.compiled = compiled
	return nil
}

func (fl *FileList) String() string {
	return fl.compiled
}

func (fp *fpath) ensureDir() error {
	chk := fp.dir + fp.name
	st, err := os.Stat(chk)
	if err != nil {
		return fmt.Errorf("func (fp *fpath) ensureDir(): os.Stat(%v): returned\n %v", chk, err.Error())
	}
	if st.IsDir() {
		fp.dir = fp.dir + fp.name
		fp.name = ""
	}
	return nil
}
