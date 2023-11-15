package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Galdoba/devtools/gpath"
	"github.com/urfave/cli/v2"
)

/*
run


*/

const (
	programName = "tabviewer"
)

var dataPath string
var configPath string

func init() {
	dataPath = gpath.StdPath("Datafile.csv", []string{".ffstuff", "data", programName}...)
	if err := gpath.Touch(dataPath); err != nil {
		panic(err.Error())
	}

	configPath = gpath.StdPath(programName+".json", []string{".config", programName}...)
	if err := gpath.Touch(configPath); err != nil {
		panic(err.Error())
	}
	data, err := os.ReadFile(configPath)
	fmt.Println(len(data))
	if len(data) == 0 {
		programConfig = defaultConfig()
		data, err = json.MarshalIndent(programConfig, "", "  ")
		if err != nil {
			panic("can't create default config: " + err.Error())
		}
		f, err := os.OpenFile(configPath, os.O_WRONLY, 0777)
		if err != nil {
			panic(err.Error())
		}
		f.Write(data)
		defer f.Close()
		println(fmt.Sprintf("default config created at %v: ", configPath))
	}
	data, err = os.ReadFile(dataPath)
	if len(data) == 0 {
		print(fmt.Sprintf("no data in %v\nupdating. . .   ", dataPath))
		err = UpdateTable()
		if err != nil {
			println("fatal error")
			panic(err.Error())
		}
		println("ok")
	}
	// if err != nil {
	// 	switch {
	// 	default:
	// 		fmt.Println("Неизвестная ошибка при проверки наличия конфига:")
	// 		println(err.Error())
	// 		panic(0)
	// 	case strings.Contains(err.Error(), "The system cannot find the file specified"), strings.Contains(err.Error(), "The system cannot find the path specified"):
	// 		fmt.Println("Config file not found")
	// 		err := os.MkdirAll(strings.TrimSuffix(configPath, programName+".json"), 0777)
	// 		if err != nil {
	// 			panic(err.Error())
	// 		}

	// 		f, err := os.Create(configPath)
	// 		if err != nil {
	// 			panic(err.Error())
	// 		}
	// 		defer f.Close()
	// 		_, err = f.Write(data)
	// 		if err != nil {
	// 			panic(err.Error())
	// 		}
	// 		fmt.Println("ok")
	// 	}
	// }
	err = json.Unmarshal(data, &programConfig)
	if err != nil {
		panic(err.Error() + "asdasd")
	}
	programConfig.path = configPath
}

func main() {

	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = programName
	app.Usage = "отображает/редактирует csv файл"
	app.Flags = []cli.Flag{}

	//ДО НАЧАЛА ДЕЙСТВИЯ
	//	tb := newTableData(gconfig.DefineProgramDirectory(programName) + "taskSpreadsheet2.csv")

	//p := tea.NewProgram(tb)
	app.Before = func(c *cli.Context) error {
		return nil
	}
	app.Commands = []*cli.Command{
		{
			Name:  "config",
			Usage: "print current config",
			Action: func(c *cli.Context) error {

				fmt.Println(programConfig.String())
				return nil
			},
		},
		{ //TODO
			Name:  "update",
			Usage: "add chat key to config from url",

			Action: func(c *cli.Context) error {

				return UpdateTable()
			},
		},
	}
	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		return nil
	}

	args0 := os.Args
	if err := app.Run(args0); err != nil {
		fmt.Printf("\napplication returned error: %v\n", err.Error())
		os.Exit(3)
	}
	// if _, err := p.Run(); err != nil {
	// 	fmt.Println("critical error:", err.Error())
	// 	os.Exit(2)
	// }

}
