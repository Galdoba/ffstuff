package statfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type StatFile struct {
	path string
	Min  int            `json:"Minimum Uses"`
	Max  int            `json:"Maximum Uses"`
	Used map[string]int `json:"Used Times"`
}

var Default_DIR = defaultDir()

func defaultDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("can't get home dir:\n  " + err.Error())
	}
	sep := string(filepath.Separator)
	return home + sep + ".galdoba" + sep + "statfiles" + sep
}

func Create(dir, name string, useCases ...string) (*StatFile, error) {
	file := dir + name + ".json"
	if _, err := os.Open(file); err == nil {
		return nil, fmt.Errorf("can't create statfile '%v': file exist %v", name, file)
	}
	sf := StatFile{}
	sf.path = file
	sf.Used = make(map[string]int)
	for _, cntnt := range useCases {
		sf.Used[cntnt] = 0
	}
	if err := Save(&sf); err != nil {
		return nil, err
	}
	return &sf, nil
}

func Save(sf *StatFile) error {
	dir := filepath.Dir(sf.path)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return fmt.Errorf("can't create dir '%v': %v", dir, err)
	}
	f, err := os.Create(sf.path)
	if err != nil {
		return fmt.Errorf("can't create file '%v': %v", sf.path, err)
	}
	defer f.Close()
	sf.updateMinMax()
	bt, err := json.MarshalIndent(&sf, "", "  ")
	if err != nil {
		return fmt.Errorf("can't marshal statfile: %v", err)
	}
	if _, err := f.Write(bt); err != nil {
		return fmt.Errorf("can't write to statfile '%v': %v", sf.path, err)
	}
	return nil
}

func Load(name string) (*StatFile, error) {
	name = name + ".json"
	_, err := os.Open(name)
	if err != nil {
		defaultName := defaultDir() + name
		_, err = os.Open(defaultDir() + name)
		if err != nil {
			return nil, fmt.Errorf("can't open file '%v'", name)
		}
		name = defaultName
	}
	bt, errRead := os.ReadFile(name)
	if errRead != nil {
		return nil, fmt.Errorf("can't read file '%v'", name)
	}
	sf := &StatFile{}
	sf.Used = make(map[string]int)
	if err := json.Unmarshal(bt, sf); err != nil {
		return nil, fmt.Errorf("can't unmarshal statfile '%v': %v", name, err)
	}
	sf.path = name
	return sf, nil
}

func (sf *StatFile) AddWeight(name string, w int) error {
	if val, ok := sf.Used[name]; ok {
		sf.Used[name] = val + w
		sf.updateMinMax()
		return Save(sf)
	}
	return fmt.Errorf("no option '%v'", name)
}

func (sf *StatFile) AddOption(opts ...string) error {
	for _, name := range opts {
		if _, ok := sf.Used[name]; ok {
			continue
		}
		sf.Used[name] = 0
	}
	return Save(sf)
}

func (sf *StatFile) RemoveOption(opts ...string) error {
	for k, _ := range sf.Used {
		delete(sf.Used, k)
	}
	return Save(sf)
}

func (sf *StatFile) ResetWeight(opts ...string) error {
	for _, name := range opts {
		if _, ok := sf.Used[name]; ok {
			sf.Used[name] = 0
		}
	}
	return Save(sf)
}

func (sf *StatFile) ResetWeightAll() error {
	for k, _ := range sf.Used {
		sf.Used[k] = 0
	}
	return Save(sf)
}

func (sf *StatFile) updateMinMax() {
	mn, mx := 0, 0
	first := true
	for _, v := range sf.Used {
		if first {
			mn, mx = v, v
			first = false
		}
		if v > mx {
			mx = v
		}
		if v < mn {
			mn = v
		}
	}

	sf.Min = mn
	sf.Max = mx
}

type weighted struct {
	opt string
	w   int
}

func (sf *StatFile) List() []string {
	list := []weighted{}
	output := []string{}
	for k, v := range sf.Used {
		list = append(list, weighted{k, v})
	}
	for i := sf.Max; i >= sf.Min; i-- {
		for _, wt := range list {
			if wt.w == i {
				output = append(output, wt.opt)
			}
		}
	}
	return output
}

/*
AddWeight (name string, weight int) error
Use(string) error
Clear(names ...string) error
ClearAll() error

Delete(name string)

List() []string
*/
