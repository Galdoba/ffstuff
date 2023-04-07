package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Galdoba/ffstuff/pkg/config"
	"gopkg.in/yaml.v3"
)

type grabberConfig struct {
	Description       string
	External_Log_path string
	Local_Log_path    string
	//Actions
	Actions []Action
}

type Action struct {
	ActionName string
	Triggers   []string
}

func CreateDefaultConfig() error {
	gc := grabberConfig{}
	gc.Description = "config file for 'grabber.exe'"

	gc.External_Log_path = "TODO"
	gc.Local_Log_path = "TODO"
	gc.Actions = []Action{
		{
			ActionName: "MOVE_CURSOR_UP",
			Triggers:   []string{"UP"},
		},
		{
			ActionName: "MOVE_CURSOR_DOWN",
			Triggers:   []string{"DOWN"},
		},
		{
			ActionName: "TOGGLE_SELECTION_STATE",
			Triggers:   []string{"SPACE"},
		},
		{
			ActionName: "DROP_SELECTIONS",
			Triggers:   []string{"~", "BACKSPACE"},
		},
		{
			ActionName: "MOVE_SELECTED_TOP",
			Triggers:   []string{"ENTER", "Ctrl+T"},
		},
		{
			ActionName: "MOVE_SELECTED_UP",
			Triggers:   []string{"W", "PgUp"},
		},
		{
			ActionName: "MOVE_SELECTED_DOWN",
			Triggers:   []string{"S", "PgDn"},
		},
		{
			ActionName: "DECIDION_CONFIRM",
			Triggers:   []string{"ENTER"},
		},
		{
			ActionName: "DECIDION_DENY",
			Triggers:   []string{"ESC"},
		},
		{
			ActionName: "DOWNLOAD_PAUSE",
			Triggers:   []string{"P"},
		},
	}
	fileBts, err := yaml.Marshal(gc)
	if err != nil {
		return err
	}

	cDir, cFile := config.ConfigPathManual(programName)
	confPath := fmt.Sprintf("%v\\%v", cDir, cFile)
	fmt.Println("will go here:", confPath)
	//confPath, err = filepath.Abs(confPath)

	if err := os.MkdirAll(cDir, 0777); err != nil {
		return err
	}
	fmt.Printf("'%v' created\n", cDir)
	//panic(confPath)
	// read the whole file at once
	_, err = os.OpenFile(confPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("'%v' opened\n", confPath)
	// write the whole body at once
	err = ioutil.WriteFile(confPath, fileBts, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadConfig(path string) (*grabberConfig, error) {
	fl, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	gc := &grabberConfig{}
	err = yaml.Unmarshal(fl, gc)
	if err != nil {
		return nil, err
	}
	if gc.Validate() != nil {
		return nil, err
	}

	return gc, nil
}

func (gc *grabberConfig) Validate() error {
	if len(gc.Actions) < 5 {
		return fmt.Errorf("")
	}
	return nil
}
