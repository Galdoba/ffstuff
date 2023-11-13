package main

// token := "6963211405:AAEQqcDDAnueLx0iE3Imu_fIj0i6mCLw8qo"
// 	chatID := int64(-4083924452)
// 	message := "aaaббб"
// 	bot, err := tgbotapi.NewBotAPI(token)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	msg := tgbotapi.NewMessage(chatID, message)
// 	msg.ParseMode = tgbotapi.ModeHTML
// 	_, err = bot.Send(msg)
// 	//fmt.Println(ms)
// 	if err != nil {
// 		panic(err.Error())
// 	}

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"

	"github.com/Galdoba/ffstuff/pkg/gconfig"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/urfave/cli/v2"
)

/*
run


*/

var configPath string

const (
	programName = "tgnotify"
)

func init() {

	configPath = gconfig.DefineConfigPath(programName)
	exs, err := fileExists(configPath)

	if !exs {
		if err != nil {
			panic(err.Error)
		}
		errEx := fmt.Sprintf("config file not exist: %v", configPath)
		println(errEx)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		switch {
		default:
			errText := fmt.Sprintf("unexpected config error: %v", err.Error())
			println(errText)
			os.Exit(1)
		}
	}

	err = json.Unmarshal(data, &programConfig)
	if err != nil {
		errText := fmt.Sprintf("can't unmarshal config data: %v", err.Error())
		println(errText)
		os.Exit(1)
	}

}

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, fmt.Errorf("file may or may not exist: %v", err.Error())
	}

}

func main() {

	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = programName
	app.Usage = "send message to telegram channel"
	app.Flags = []cli.Flag{}

	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {

		return nil
	}
	app.Commands = []*cli.Command{

		{
			Name:      "send",
			Usage:     "send text message",
			ArgsUsage: "no arguments allowed",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "sign",
					Usage:   "add user name to text",
					Aliases: []string{"s"},
				},
				&cli.StringFlag{
					Name:     "message",
					Usage:    "set message text",
					Required: true,
					Aliases:  []string{"m"},
				},
				&cli.StringFlag{
					Name:    "title",
					Usage:   "set text title",
					Aliases: []string{"t"},
				},
				&cli.StringFlag{
					Name:    "postscript",
					Usage:   "set text post scriptum",
					Aliases: []string{"ps"},
				},
			},
			Action: func(c *cli.Context) error {
				message := c.String("message")
				if message == "" {
					return fmt.Errorf("message MUST not be empty")

				}
				if len(c.Args().Slice()) != 0 {
					return fmt.Errorf("action 'send' must not use arguments. \ncheck if text of the message have spaces and not encaplated with quotes" + ` (")`)
				}
				chatID := int64(programConfig.ChatID)
				token := programConfig.Token
				bot, err := tgbotapi.NewBotAPI(token)
				if err != nil {
					return fmt.Errorf("create bot api: %v", err.Error())
				}
				title := c.String("title")
				if title != "" {
					message = title + "\n" + message

				}
				ps := c.String("ps")
				if ps != "" {
					message = message + "\n" + ps
				}

				if c.Bool("sign") {
					usr, err := user.Current()
					if err != nil {
						println("can't get user name")
					}
					message = "from user: " + usr.Name + "\n" + message
				}

				msg := tgbotapi.NewMessage(chatID, message)
				msg.ParseMode = tgbotapi.ModeHTML
				_, err = bot.Send(msg)
				if err != nil {
					return fmt.Errorf("send message: %v", err.Error())
				}
				return nil
			},
		},
		{
			Name:  "config",
			Usage: "print current config",
			Action: func(c *cli.Context) error {
				fmt.Println("Config path :", configPath)
				fmt.Println("     ChatID :", programConfig.ChatID)
				fmt.Println("      Token :", programConfig.Token)
				return nil
			},
		},
	}

	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		return nil
	}
	args := os.Args
	if err := app.Run(args); err != nil {
		errOut := fmt.Sprintf("%v error: %v", programName, err.Error())
		println(errOut)
		os.Exit(1)
	}

}
