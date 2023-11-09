package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/gconfig"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

/*
run


*/

const (
	programName = "gcsvviewer"
)

func init() {

	configPath := gconfig.DefineConfigPath(programName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		switch {
		default:
			fmt.Println("Неизвестная ошибка при проверки наличия конфига:")
			println(err.Error())
			panic(0)
		case strings.Contains(err.Error(), "The system cannot find the file specified"), strings.Contains(err.Error(), "The system cannot find the path specified"):
			fmt.Println("Config file not found")
			err := os.MkdirAll(strings.TrimSuffix(configPath, programName+".json"), 0777)
			if err != nil {
				panic(err.Error())
			}
			programConfig = defaultConfig()
			data, err = json.MarshalIndent(programConfig, "", "  ")
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Creating default config at %v: ", configPath)
			f, err := os.Create(configPath)
			if err != nil {
				panic(err.Error())
			}
			defer f.Close()
			_, err = f.Write(data)
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("ok")
		}
	}

	err = json.Unmarshal(data, &programConfig)
	if err != nil {
		panic(err.Error())
	}

}

func main() {

	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = programName
	app.Usage = "отображает/редактирует csv файл"
	app.Flags = []cli.Flag{}

	//ДО НАЧАЛА ДЕЙСТВИЯ
	tb := newTableData(gconfig.DefineProgramDirectory(programName) + "taskSpreadsheet2.csv")
	fmt.Println(len(tb.data))
	p := tea.NewProgram(tb)
	app.Before = func(c *cli.Context) error {

		return nil
	}
	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		return nil
	}

	//args0 := os.Args
	// if err := app.Run(args0); err != nil {
	// 	fmt.Printf("\napplication returned error: %v\n", err.Error())
	// 	os.Exit(3)
	// }
	if _, err := p.Run(); err != nil {
		fmt.Println("critical error:", err.Error())
		os.Exit(2)
	}

}