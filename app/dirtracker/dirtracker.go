package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/Galdoba/ffstuff/app/dirtracker/config"
	"github.com/Galdoba/ffstuff/app/dirtracker/filelist"
	"github.com/Galdoba/utils"
	"github.com/urfave/cli/v2"
)

const (
	program = "dirtracker"
)

var cfg config.Config

func main() {
	app := cli.NewApp()
	app.Version = "v 0.2.2"
	app.Name = "dirtracker"
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
		return nil
	}
	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		//Придумаем что-нибудь
		return nil
	}
	app.Commands = []*cli.Command{
		//start - запустить программу
		//ShowStats - поразать глубокую аналитику
		{
			Name:        "start",
			Usage:       "запустить программу",
			UsageText:   "Собирает необходимую информацию из config.file после чего формирует общий список и пропускает его через фильтры для формирования короткого списка",
			Description: "TODO: подробное описание команды",
			ArgsUsage:   "Аргументов не имеет\nВ планах локальный режим и указание файла в который должен писаться отчет",
			Action: func(c *cli.Context) error {

				runtime.GOMAXPROCS(runtime.NumCPU() - 1) // занимаем все ядра кроме одного чтобы система не висла в случае избыточного кол-ва тредов

				fl, _ := filelist.New(cfg.SearchRoot())
				output := ""
			mainLoop:
				for {
					//fmt.Printf("Updating...                                   \n")
					atempt := 1
					for atempt <= 100 {
						if err := fl.Update(cfg.MaximumSearchThreads()); err != nil {
							fmt.Print("\rTry ", atempt, " ", err.Error()) //на случай если будет ошибка обновления списка
							time.Sleep(time.Second)
						} else {
							break
						}
						atempt++
						if atempt > 10 {
							return fmt.Errorf("to many atempts to update list")
						}
					}

					shortList := filelist.Compile(fl.FullList(), cfg.WhiteList(), cfg.WhiteListEnabled(), cfg.BlackList(), cfg.BlackListEnabled())

					res, err := filelist.Format(shortList, cfg.WhiteList(), cfg.WhiteListEnabled())
					if err != nil {
						if err.Error() == "no files found" {
							utils.ClearScreen()
							fmt.Println("NO FILES FOUND            ")
							time.Sleep(time.Second * 5)
							continue mainLoop
						}
					}
					utils.ClearScreen()
					//stats := fl.Stats()
					//fmt.Printf("Found %v files in %v directories with %v errors\n", stats["file"], stats["dir"], stats["err"])
					sendToBot := false
					if output != res {
						if output != "" {
							sendToBot = true
						}
					}
					output = res
					fmt.Println(output)

					switch sendToBot {
					case true:
						//fmt.Println("Sending To Bot")

					case false:
						//fmt.Println("NOT Sending To Bot")
					}

					updCyc := cfg.UpdateCycle()
					for i := updCyc; i > -1; i-- {
						switch i {
						default:
							//	fmt.Printf("Next update in %v seconds...             \r", i)
						}
						time.Sleep(time.Second)
					}

				}
				//fmt.Println(output)
				//return nil
			},
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
		fmt.Printf("application returned error: %v", err.Error())
	}
	// exit := ""
	// val := survey.ComposeValidators()
	// promptInput := &survey.Input{
	// 	Message: "Enter для завершения",
	// }
	// survey.AskOne(promptInput, &exit, val)

}
