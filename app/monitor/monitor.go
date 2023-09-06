package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Galdoba/devtools/directory"
	"github.com/Galdoba/ffstuff/internal/terminal"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/urfave/cli/v2"

	tsize "github.com/kopoli/go-terminal-size"

	"gopkg.in/yaml.v3"
)

const (
	program = "monitor"
)

type ConfFile struct {
	Roots               map[string][]string
	Storage             string
	UpdateCycle_seconds int
	Max_threads         int
}

var Conf *ConfFile

var opsys string
var storagePath string
var storageFile string

func main() {
	storageFile = "session.info"
	//dv := []string{"buffer"}
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
		&cli.StringSliceFlag{
			Name:      "roots",
			Usage:     "defines whitch root to print (prints first by default)",
			FilePath:  "",
			Required:  false,
			Hidden:    false,
			TakesFile: false,
			//Value:     &cli.StringSlice{""},
		},
	}
	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {
		opsys = runtime.GOOS
		switch opsys {
		default:
			fmt.Printf("this is unsupported Operating System: %v", opsys)
			os.Exit(1)
		case "windows":
			storagePath = config.DataDirectory(program)
		case "linux":
			storagePath = config.DataDirectory(program)
		}

		//зачищаем остатки данных с прошлой сессии
		cfg := config.File{}
		// if err := os.RemoveAll(storagePath); err != nil {
		// 	return err
		// }
		if err := os.MkdirAll(storagePath, os.ModePerm); err != nil {
			return err
		}
		dataStore, err := os.OpenFile(storagePath+storageFile, os.O_CREATE, 0511)
		if err != nil {
			return err
		}
		defer dataStore.Close()
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
		runtime.GOMAXPROCS(runtime.NumCPU() - 1)
		return nil
	}
	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		return nil
	}
	app.Commands = []*cli.Command{
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
				list, err := scanRoots(c)
				if err != nil {
					return fmt.Errorf("scanRoots(c): %v", err.Error())
				}
				for _, l := range list {
					fmt.Println(l)
				}
				if err := updateStoredInfo(list); err != nil {
					return err
				}

				sendToBot := false
				switch sendToBot {
				case true:
				case false:
				}
				return nil
			},

			Subcommands: []*cli.Command{},
			Flags: []cli.Flag{

				&cli.BoolFlag{
					Name:  "files",
					Usage: "show files only",
				},
				&cli.BoolFlag{
					Name:  "dirs",
					Usage: "show directories only",
				},
			},
			SkipFlagParsing:        false,
			HideHelp:               false,
			Hidden:                 false,
			UseShortOptionHandling: false,
			HelpName:               "",
			CustomHelpTemplate:     "",
		},
		{
			Name:        "track",
			Usage:       "monitor track [--loop x] [--buffer path]",
			UsageText:   "keep updated and formated info about content on the screen",
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

				loop := 1
				lastScr := ""
				for loop > 0 {
					list, err := scanRoots(c)
					if err != nil {
						return err
					}
					if err := updateStoredInfo(list); err != nil {
						return err
					}

					s, _ := tsize.GetSize()
					if s.Width < 40 {
						return fmt.Errorf("console width is to small: %v (minimum 40)", s.Width)
					}
					scr := onScreenBW(s.Width)
					if scr != lastScr {
						lastScr = scr
						terminal.Clear()
						fmt.Println(scr)
					}
					time.Sleep(time.Second * time.Duration(Conf.UpdateCycle_seconds))

				}
				sendToBot := false
				switch sendToBot {
				case true:
				case false:
				}
				return nil
			},

			Subcommands: []*cli.Command{},
			Flags: []cli.Flag{

				&cli.StringFlag{
					Name:     "width",
					Category: "Output",
					Usage:    "S: [20-80]; M: [81-130]; L: [131+]",
					Value:    "S",
					Aliases:  []string{"w"},
				},
				// 	&cli.BoolFlag{
				// 		Name:  "dirs",
				// 		Usage: "show directories only",
				// 	},
			},
			SkipFlagParsing:        false,
			HideHelp:               false,
			Hidden:                 false,
			UseShortOptionHandling: false,
			HelpName:               "",
			CustomHelpTemplate:     "",
		},
		{
			Name:        "config",
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
		fmt.Printf("application returned error: %v\n", err.Error())
	}
	// exit := ""
	// val := survey.ComposeValidators()
	// promptInput := &survey.Input{
	// 	Message: "Enter для завершения",
	// }
	// survey.AskOne(promptInput, &exit, val)

}

func scanRoots(c *cli.Context) ([]string, error) {
	rootsUnsorted := Conf.Roots
	roots := make(map[string][]string)
	list := []string{}
	sep := string(filepath.Separator)
	// dataStore, err := os.OpenFile(storagePath+storageFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	// if err != nil {
	// 	return list, err
	// }
	// defer dataStore.Close()

	validRoots := c.StringSlice("roots")
	for k, root := range rootsUnsorted {
		if len(validRoots) == 0 {
			roots[k] = root
		}
		for _, k2 := range validRoots {
			if k == k2 {
				roots[k] = root
			}
		}
	}
	for _, rootList := range roots {

		for _, v := range rootList {
			dirs := []string{}
			fls := []string{}
			dir, files, err := directory.List(v)
			if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
				return list, err
			}
			if err != nil {
				return list, fmt.Errorf("directory.List(%v): %v", v, err.Error())
			}
			list = append(list, dir)
			for _, fl := range files {
				cntnt := dir + fl
				f, _ := os.Stat(cntnt)
				switch f.IsDir() {
				case true:
					if !c.Bool("files") {
						dirs = append(dirs, cntnt+sep)
					}
				case false:
					if !c.Bool("dirs") {
						fls = append(fls, cntnt)

						//	dataStore.Write([]byte(cntnt + "  \n"))
					}
				}
			}
			list = append(list, dirs...)
			list = append(list, fls...)
		}
	}
	return list, nil
}
