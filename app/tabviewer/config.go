package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/devtools/cli/command"
)

type config struct {
	path         string
	UpdateTicker int               `json:"Update cycle (seconds)"`
	Curl         string            `json:"CURL request"` //-s --use-ascii --proxy http://proxy.local:3128 https://docs.google.com/spreadsheets/d/1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg/gviz/tq?tqx=out:csv -k --output
	CSV_DataFile string            `json:"CSV path"`     //c:\Users\pemaltynov\.ffstuff\data\taskSpreadsheet.csv
	KeyLayout    map[string]string `json:"Key Layout,omitempty"`
	ActivePreset string            `json:"Active Preset"`
}

//curl --use-ascii --proxy http://proxy.local:3128 https://docs.google.com/spreadsheets/d/1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg/gviz/tq?tqx=out:csv -k --output c:\Users\pemaltynov\.ffstuff\data\taskSpreadsheet2.csv

func UpdateTable() error {
	_, err := command.RunSilent("curl", programConfig.Curl+programConfig.CSV_DataFile+".tmp")
	if err != nil {
		return err
	}
	newPath := programConfig.CSV_DataFile + ".tmp"
	println("emulate checking")
	oldPath := programConfig.CSV_DataFile
	return os.Rename(newPath, oldPath)

}

var programConfig *config

func (cfg *config) String() string {
	str := fmt.Sprintf("path: %v\n", cfg.path)
	str += fmt.Sprintf("data file: %v\n", cfg.CSV_DataFile)
	if cfg.UpdateTicker > 0 {
		str += fmt.Sprintf("Auto Update: %v\n", cfg.UpdateTicker)
	}
	return str
}

func defaultConfig() *config {
	cfg := config{}
	cfg.path = gconfig.DefineConfigPath(programName)
	cfg.UpdateTicker = 10
	cfg.Curl = "-s --use-ascii --proxy http://proxy.local:3128 https://docs.google.com/spreadsheets/d/1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg/gviz/tq?tqx=out:csv -k --output "
	cfg.CSV_DataFile = dataPath
	cfg.ActivePreset = "Default"
	return &cfg
}

//https://docs.google.com/spreadsheets/d/1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg/gviz/tq?tqx=out:csv&sheet=График работ 2.0
//https://docs.google.com/spreadsheets/d/1pIrdfVbUy5I9NF70USDMfgr_H8by3CwGJstFfWTTaug/gviz/tq?tqx=out:json&sheet=Factions
