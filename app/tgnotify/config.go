package main

import (
	"github.com/Galdoba/ffstuff/pkg/gconfig"
)

type config struct {
	path   string
	ChatID int64  `json:"ChatID"` //-s --use-ascii --proxy http://proxy.local:3128 https://docs.google.com/spreadsheets/d/1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg/gviz/tq?tqx=out:csv -k --output
	Token  string `json:"Token"`
}

//curl --use-ascii --proxy http://proxy.local:3128 https://docs.google.com/spreadsheets/d/1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg/gviz/tq?tqx=out:csv -k --output c:\Users\pemaltynov\.ffstuff\data\taskSpreadsheet2.csv

// func UpdateTable() error {
// 	command.RunSilent("curl", "")
// 	// comm, err := command.New(
// 	// 	command.CommandLineArguments("curl "+sp.curl+sp.csvPath),
// 	// 	command.Set(command.BUFFER_OFF),
// 	// 	command.Set(command.TERMINAL_ON),
// 	// )
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// fmt.Println("Updating Spreadsheet:")
// 	// comm.Run()
// 	// if err := sp.fillCSVData(); err != nil {
// 	// 	return fmt.Errorf("sp Update(): sp.fillCSVData() = %v", err.Error())
// 	// }
// 	// fmt.Println("Update Status: ok")
// 	return nil
// }

var programConfig *config

func defaultConfig() *config {
	cfg := config{}
	cfg.path = gconfig.DefineConfigPath(programName)

	return &cfg
}
