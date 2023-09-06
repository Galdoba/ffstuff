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
	program = "kval"
)

func main() {
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = "kval"
	app.Usage = "Оперирует коллекциями Ключ-Значение"
	app.Description = "CRUD - позволяет создавать/редактировать и удалять коллекции и отдельные пары Key-Val"
	app.Flags = []cli.Flag{
		// cli.BoolFlag{
		// 	Name:  "update",
		// 	Usage: "если активен, то до начала выполнения любой команды - обновится csv с рабочей таблицей",
		// },
	}
	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {
	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		//Придумаем что-нибудь
		return nil
	}
	app.Commands = []cli.Command{
/*
collection file theme call list 

kval print --list amedia
kval get --list amedia --key Pravednye_dzhemstouny
kval set --list amedia --key I_prosto_tak --val "./IN/@SCRIPTS/amedia_ar6e2.sh I_prosto_tak"
kval del --list amedia
kval new --list wink
kval set --list wink --key Perevozchik --val "./IN/@SCRIPTS/amedia_ar2.sh Perevozchik"
curl -OL https://github.com/ryanoasis/nerd-fonts/releases/latest/download/GoMonoNerdFont.tar.xz
*/
		{
			Name:        "new",
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

					shortList := filelist.Compile(fl.FullList(), Conf.WhiteList, Conf.WhiteListEnabled, Conf.BlackList, Conf.BlackListEnabled)

					res, err := filelist.Format(shortList, Conf.WhiteList, Conf.WhiteListEnabled)
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
