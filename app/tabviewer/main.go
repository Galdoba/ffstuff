package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Galdoba/devtools/gpath"
	tsize "github.com/kopoli/go-terminal-size"
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
	configPath = gpath.StdPath(programName+".json", []string{".config", programName}...)
	dataPath = gpath.StdPath("Datafile.csv", []string{".ffstuff", "data", programName}...)

	for _, err := range []error{
		checkConfig(),
		checkDataFile(),
	} {
		if err != nil {
			panic(err.Error())
		}
	}

}

func checkConfig() error {
	if err := gpath.Touch(configPath); err != nil {
		return fmt.Errorf("can't confirm config path: " + err.Error())
	}
	data, err := os.ReadFile(configPath)
	if len(data) == 0 {
		programConfig = defaultConfig()
		data, err = json.MarshalIndent(programConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("can't create default config: " + err.Error())
		}
		f, err := os.OpenFile(configPath, os.O_WRONLY, 0777)
		if err != nil {
			return fmt.Errorf("can't open config: " + err.Error())
		}
		_, err = f.Write(data)
		if err != nil {
			return fmt.Errorf("can't write config: " + err.Error())
		}
		defer f.Close()
		println(fmt.Sprintf("default config created: %v", configPath))
	}
	err = json.Unmarshal(data, &programConfig)
	if err != nil {
		return fmt.Errorf("can't unmarhal config: %v", err.Error())
	}
	return nil
}

func checkDataFile() error {
	if err := gpath.Touch(dataPath); err != nil {
		return fmt.Errorf("can't confirm dataPath: %v", err.Error())
	}
	programConfig.path = configPath
	data, err := os.ReadFile(dataPath)
	if err != nil {
		return fmt.Errorf("can't read dataPath: %v", err.Error())
	}
	if len(data) == 0 {
		print(fmt.Sprintf("no data in %v\nupdating. . .   ", dataPath))
		err = UpdateTable()
		if err != nil {
			return fmt.Errorf("can't update %v: %v", dataPath, err.Error())
		}
		println("ok")
	}
	return nil
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
		{ //config
			Name:  "config",
			Usage: "Print current config",
			Action: func(c *cli.Context) error {
				_, err := json.Marshal(programConfig)
				if err != nil {
					return err
				}
				fmt.Println(programConfig)
				return nil
			},
		},
		{ //update
			Name:  "update",
			Usage: "Update current tabledata with curl",
			Action: func(c *cli.Context) error {
				println("Updating...")
				err := UpdateTable()
				if err != nil {
					return err
				}
				println("ok")
				return nil
			},
		},
		{ //run
			Name:  "run",
			Usage: "Show table",
			Action: func(c *cli.Context) error {
				tablefile, err := os.OpenFile(dataPath, os.O_RDWR, 0777)
				if err != nil {
					return fmt.Errorf("can't read data file: %v", err.Error())
				}
				defer tablefile.Close()
				csvReader := csv.NewReader(tablefile)
				data, err := csvReader.ReadAll()

				columnLen := columnSizes(data)

				for _, line := range data {
					fmt.Println(FormatLineSize(line, columnLen))

				}

				sz, err := tsize.GetSize()
				width := sz.Width
				for i, line := range data {
					if i > 50 {
						fmt.Println(FormatLine(line, width))
					}
				}
				fmt.Println(columnLen)
				return nil
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
