package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Galdoba/devtools/directory"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v3"
)

const (
	program = "monitor"
)

type ConfFile struct {
	Roots               map[string][]string
	UpdateCycle_seconds int
	Max_threads         int
}

var Conf *ConfFile

func main() {
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = "monitor"
	app.Usage = "Отслеживает файлы в указанных директориях"
	app.Description = "После сбора информации об имеющихся файлах и папках ниже корня готовит отчеты для вывода в файл/на терминал"
	app.Flags = []cli.Flag{
		// cli.BoolFlag{
		// 	Name:  "update",
		// 	Usage: "если активен, то до начала выполнения любой команды - обновится csv с рабочей таблицей",
		// },
	}
	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {
		//Убедиться что есть файлы статистики. если нет то создать
		cfg := config.File{}
		if !config.Exists(program) {
			fmt.Println("Config not detected...")
			cfg, err := config.ConstructManual(program)
			if err != nil {
				return err
			}
			fmt.Println("Filling default Data...")
			roots := make(map[string][]string)
			roots["buffer"] = []string{`\\192.168.31.4\buffer\IN\_DONE\`, `\\192.168.31.4\buffer\IN\_IN_PROGRESS\`, `\\192.168.31.4\buffer\IN\`}
			conf := ConfFile{
				Roots:               roots,
				UpdateCycle_seconds: 5,
				Max_threads:         5,
			}
			dt, err := yaml.Marshal(conf)
			if err := cfg.Write(dt); err != nil {
				return err
			}
			fmt.Println("Default data filled...")
			return nil
		}
		if err := yaml.Unmarshal(cfg.Data, ConfFile{}); err != nil {
			return err
		}
		return nil
	}
	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		//Придумаем что-нибудь
		return nil
	}
	app.Commands = []cli.Command{
		//start - запустить программу
		//ShowStats - поразать глубокую аналитику
		{
			Name:        "print",
			Usage:       "monitor print [rootkey] [-destination]...",
			UsageText:   "print list of content in tracked directories",
			Description: "TODO: подробное описание команды",
			ArgsUsage:   "Аргументов не имеет\nВ планах локальный режим и указание файла в который должен писаться отчет",
			Category:    "Primary",
			Action: func(c *cli.Context) error {
				cfgData, err := config.ReadFrom(program)
				if err != nil {
					return fmt.Errorf("config.ReadFrom: %v", err)
				}
				if err := yaml.Unmarshal(cfgData, &Conf); err != nil {
					return fmt.Errorf("yaml.Unmarshal: %v", err)
				}
				if len(Conf.Roots) < 1 {
					return fmt.Errorf("no Roots set in config.file")
				}
				runtime.GOMAXPROCS(runtime.NumCPU() - 1)
				roots := Conf.Roots
				list := []string{}
				sep := string(filepath.Separator)
				for _, rootList := range roots {
					for _, v := range rootList {
						dirs := []string{}
						fls := []string{}
						dir, files, err := directory.List(v)
						if err != nil {
							return fmt.Errorf("directory.List(%v): %v", v, err.Error())
						}
						list = append(list, dir)
						for _, fl := range files {
							cntnt := dir + fl
							f, _ := os.Stat(cntnt)
							switch f.IsDir() {
							case true:
								if !c.Bool("files") {
									dirs = append(dirs, ".."+sep+fl+sep)
								}
							case false:
								if !c.Bool("dirs") {
									fls = append(fls, fl)
								}
							}
						}
						list = append(list, dirs...)
						list = append(list, fls...)
					}
				}
				for _, l := range list {
					fmt.Println(l)
				}
				sendToBot := false
				switch sendToBot {
				case true:
				case false:
				}
				return nil
			},

			Subcommands: []cli.Command{},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "files",
					Usage: "show files only",
				},
				cli.BoolFlag{
					Name:  "dirs",
					Usage: "show directories only",
				},
			},
			SkipFlagParsing:        false,
			SkipArgReorder:         false,
			HideHelp:               false,
			Hidden:                 false,
			UseShortOptionHandling: false,
			HelpName:               "",
			CustomHelpTemplate:     "",
		},
		{
			Name:        "config",
			ShortName:   "",
			Usage:       "Показывает информацию о текущих настройках",
			UsageText:   "ТУДУ: описание как использовать команду",
			Description: "TODO: подробное описание команды",
			ArgsUsage:   "",
			Action: func(c *cli.Context) error {

				cfgData, err := config.ReadFrom(program)
				if err != nil {
					return err
				}
				fmt.Printf("Current config is:\n")
				fmt.Println("--------------------------------------------------------------------------------")
				fmt.Println(string(cfgData))
				fmt.Println("--------------------------------------------------------------------------------")
				fmt.Println("File location: ", config.Filepath(program))
				return nil
			},
		},
	}

	args := os.Args
	if err := app.Run(args); err != nil {
		fmt.Printf("application returned error: %v", err.Error())
	}
	// exit := ""
	// val := survey.ComposeValidators()
	// promptInput := &survey.Input{
	// 	Message: "Enter для завершения",
	// }
	// survey.AskOne(promptInput, &exit, val)

}
