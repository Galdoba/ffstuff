package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/Galdoba/ffstuff/app/dirtracker/filelist"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/utils"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v3"
)

const (
	program = "dirtracker"
)

type ConfFile struct {
	Root                string
	WhiteListEnabled    bool
	WhiteList           []string
	BlackListEnabled    bool
	BlackList           []string
	UpdateCycle_seconds int
	Max_threads         int
}

var Conf *ConfFile

func main() {
	app := cli.NewApp()
	app.Version = "v 0.2.1"
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
		cfg := config.File{}
		if !config.Exists(program) {
			fmt.Println("Config not detected...")
			cfg, err := config.ConstructManual(program)
			if err != nil {
				return err
			}
			fmt.Println("Filling default Data...")
			conf := ConfFile{
				Root:                "\\\\nas\\buffer\\IN\\",
				WhiteListEnabled:    false,
				WhiteList:           []string{"DIR1\\", "DIR2\\"},
				BlackListEnabled:    false,
				BlackList:           []string{"DIR1\\", "DIR2\\"},
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
		fmt.Println("Config detected...")
		if err := yaml.Unmarshal(cfg.Data, ConfFile{}); err != nil {
			return err
		}
		fmt.Println("Content is valid...")
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
			Name:        "start",
			ShortName:   "",
			Usage:       "запустить программу",
			UsageText:   "Собирает необходимую информацию из config.file после чего формирует общий список и пропускает его через фильтры для формирования короткого списка",
			Description: "TODO: подробное описание команды",
			ArgsUsage:   "Аргументов не имеет\nВ планах локальный режим и указание файла в который должен писаться отчет",
			Action: func(c *cli.Context) error {
				cfgData, err := config.ReadFrom(program) //считываем содержание конфига
				if err != nil {
					return fmt.Errorf("config.ReadFrom: %v", err)
				}
				if err := yaml.Unmarshal(cfgData, &Conf); err != nil { //интерпритируем конфиг для дальнейшего пользования
					return fmt.Errorf("yaml.Unmarshal: %v", err)
				}
				runtime.GOMAXPROCS(runtime.NumCPU() - 1) // занимаем все ядра кроме одного чтобы система не висла в случае избыточного кол-ва тредов

				fl, _ := filelist.New(Conf.Root)
				output := ""
			mainLoop:
				for {
					//fmt.Printf("Updating...                                   \n")
					atempt := 1
					for atempt <= 100 {
						if err := fl.Update(Conf.Max_threads); err != nil {
							fmt.Print("Try ", atempt, " ", err.Error(), "\n") //на случай если будет ошибка обновления списка
							time.Sleep(time.Second)
						} else {
							fmt.Printf("Update Complete\r")
							break
						}
						atempt++
						if atempt > 10 {
							return fmt.Errorf("to many atempts to update list")
						}
					}

					shortList := filelist.Compile(fl.FullList(), Conf.WhiteList, Conf.WhiteListEnabled, Conf.BlackList, Conf.BlackListEnabled)

					res, err := filelist.Format(shortList, Conf.WhiteList, Conf.WhiteListEnabled)
					if err != nil {
						if err.Error() == "no files found" && output != "" {
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

					updCyc := Conf.UpdateCycle_seconds
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
