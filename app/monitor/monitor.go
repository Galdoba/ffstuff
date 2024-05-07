package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Galdoba/devtools/directory"
	"github.com/Galdoba/ffstuff/app/monitor/config"
	"github.com/Galdoba/ffstuff/internal/terminal"
	"github.com/urfave/cli/v2"

	tsize "github.com/kopoli/go-terminal-size"
)

const (
	program = "monitor"
)

var cfg config.Config

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

		//зачищаем остатки данных с прошлой сессии
		err := fmt.Errorf("config not loaded")
		cfg, err = config.Load()
		if err != nil {
			cfg = config.New()
			cfg.SetDefault()
			if err := cfg.Save(); err != nil {
				fmt.Printf("initialisation failed: %v", err.Error())
				os.Exit(1)

			}
			fmt.Printf("config file generated at %v \n", cfg.Path())
			fmt.Println("restart application")
			os.Exit(0)

		}
		storagePath = cfg.DataStorage()

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
				// cfgData, err := config.ReadFrom(program)
				// if err != nil {
				// 	return fmt.Errorf("config.ReadFrom: %v", err)
				// }
				// if err := yaml.Unmarshal(cfgData, &Conf); err != nil {
				// 	return fmt.Errorf("yaml.Unmarshal: %v", err)
				// }
				if len(cfg.TrackRoots()) < 1 {
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
				// cfgData, err := os.ReadFile(cfg.Path())
				// if err != nil {
				// 	return fmt.Errorf("config.ReadFrom: %v", err)
				// }
				// if err := yaml.Unmarshal(cfgData, &Conf); err != nil {
				// 	return fmt.Errorf("yaml.Unmarshal: %v", err)
				// }
				if len(cfg.TrackRoots()) < 1 {
					return fmt.Errorf("no Roots set in config.file")
				}

				loop := 1
				lastScr := ""
				s, err := tsize.GetSize()
				if err != nil {
					return fmt.Errorf("width checking: %v", err.Error())
				}
				width := s.Width
				for loop > 0 {
					// s, _ := tsize.GetSize()
					// if s.Width < 40 {
					// 	if s.Width != 0 {
					// 		return fmt.Errorf("console width is to small: %v (minimum 40)", s.Width)
					// 	}
					// }
					list, err := scanRoots(c)
					if err != nil {
						return err
					}
					if err := updateStoredInfo(list); err != nil {
						return err
					}

					scr := onScreenBW(width)
					if scr != lastScr {
						lastScr = scr
						terminal.Clear()
						fmt.Println(scr)
					}
					time.Sleep(time.Second * time.Duration(cfg.UpdateCycle()))

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

				cfgData, err := os.ReadFile(cfg.Path())
				if err != nil {
					return err
				}
				fmt.Printf("Current config is:\n")
				fmt.Println("--------------------------------------------------------------------------------")
				fmt.Println(string(cfgData))
				fmt.Println("--------------------------------------------------------------------------------")
				fmt.Println("File location: ", cfg.Path())
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
	rootsUnsorted := cfg.TrackRoots()
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
