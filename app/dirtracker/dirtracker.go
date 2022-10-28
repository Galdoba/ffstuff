package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Galdoba/ffstuff/app/dirtracker/filelist"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/utils"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/yaml.v3"
)

const (
	program = "dirtracker"
)

type confFile struct {
	Root                string
	WhiteListEnabled    bool
	WhiteList           []string
	BlackListEnabled    bool
	BlackList           []string
	UpdateCycle_seconds int
}

func main() {
	app := cli.NewApp()
	app.Version = "v 0.0.0"
	app.Name = "dirtracker"
	app.Usage = "Отслеживает файлы в указанных директориях"
	app.Description = "Должен крутиться постоянно чтобы считывать изменения и вести статистику"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "update",
			Usage: "если активен, то до начала выполнения любой команды - обновится csv с рабочей таблицей",
		},
		cli.StringFlag{
			Name:  "tofile",
			Usage: "указывает адрес файла в который будет добавляться вывод терминала",
			Value: "",
		},
		cli.BoolFlag{
			Name:  "process, p",
			Usage: "Если активен, программа запустит ffmpeg с полученной командной строкой",
		},
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
			conf := confFile{
				Root:                "\\\\nas\\buffer\\IN\\",
				WhiteListEnabled:    false,
				WhiteList:           []string{"DIR1\\", "DIR2\\"},
				BlackListEnabled:    false,
				BlackList:           []string{"DIR1\\", "DIR2\\"},
				UpdateCycle_seconds: 5,
			}
			dt, err := yaml.Marshal(conf)
			if err := cfg.Write(dt); err != nil {
				return err
			}
			fmt.Println("Default data filled...")
			return nil
		}
		fmt.Println("Config detected...")
		//yaml.Unmarshal(cfg.Data, confFile{})
		if err := yaml.Unmarshal(cfg.Data, confFile{}); err != nil {
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
			UsageText:   "TODO: прописать основную логику",
			Description: "TODO: подробное описание команды",
			ArgsUsage:   "TODO: подробное описание как пользовать аргументы",
			Action: func(c *cli.Context) error {
				//Непосредственно команда
				/*
					Составить Список имеющихся файлов
					Вывести Список
				*/
				//TODO: Вывести параметры root и exeptions в конфиг
				cfgData, err := config.ReadFrom(program)
				if err != nil {
					return err
				}
				fl, _ := filelist.New(cfgData)
				d, f := fl.Stats()
				//FillConfig()
				utils.ClearScreen()
				fmt.Printf("found %v files in %v directories\n", f, d)
				output := ""
				/*



				 */
				for {
					fl.Update()
					if err := fl.Compile(); err != nil {
						output = "No files compiled"
						utils.ClearScreen()
					} else {
						if output != fl.String() {
							utils.ClearScreen()
							//	fmt.Println(fl.String())
						}
						output = fl.String()
						utils.ClearScreen()
						fmt.Println(output)
					}
					updCyc := fl.NextUpdate()
					for i := updCyc; i > -1; i-- {
						fmt.Printf("Next update in %v seconds...             \r", i)
						time.Sleep(time.Second)
					}

				}
				fmt.Println(output)
				return nil
			},
		},
	}

	args := os.Args
	if err := app.Run(args); err != nil {
		fmt.Printf("application returned error: %v", err.Error())
	}
	exit := ""
	val := survey.ComposeValidators()
	promptInput := &survey.Input{
		Message: "Enter для завершения",
	}
	survey.AskOne(promptInput, &exit, val)

}
