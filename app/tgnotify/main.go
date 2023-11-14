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
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		data, err := json.MarshalIndent(defaultConfig(), "", "  ")
		if err != nil {
			println(err.Error())
			//os.Exit(1)
		}
		f, _ := os.OpenFile(configPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		f.Write(data)
		println("config template created")
		println("fill TOKEN and ChatID")
		os.Exit(0)
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
					Name:     "to_chat",
					Usage:    "where message to send to",
					Required: true,
					Aliases:  []string{"tc"},
				},
				&cli.StringFlag{
					Name:     "message",
					Usage:    "set message text",
					Required: true,
					Aliases:  []string{"m", "text"},
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

				// chatID := int64(programConfig.ChatID)
				// chatID_data := strings.Split(c.String("chat"), "_")
				// switch len(chatID_data) {
				// case 1:

				// }
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
					signanure := usr.Name
					if signanure == "" {
						signanure = usr.Username
					}
					message = "from user: " + signanure + "\n" + message
				}

				chatKey := c.String("to_chat")
				if _, ok := programConfig.ChatData[chatKey]; ok != true {
					return fmt.Errorf("no key '%v' found in config file", chatKey)
				}

				chatID, topic, err := ProcessInfo(chatKey)
				if err != nil {
					return err
				}

				msg := tgbotapi.NewMessage(chatID, message)
				msg.ParseMode = tgbotapi.ModeHTML
				if topic > -1 {
					msg.ReplyToMessageID = int(topic)
				}

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
				fmt.Println(programConfig.String())
				return nil
			},
		},
		// { TODO
		// 	Name:  "add_chat",
		// 	Usage: "add chat key to config from url",
		// 	Action: func(c *cli.Context) error {
		// 		fmt.Println("Config path :", configPath)
		// 		fmt.Println("  Chat Data :", programConfig.ChatData)
		// 		fmt.Println("      Token :", programConfig.Token)
		// 		return nil
		// 	},
		// },
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

/*
tgnotyfier send -t "--------------------" -m "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." -ps "PS: Владыка, услышь меня!"
*/
