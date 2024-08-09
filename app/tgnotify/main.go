package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/Galdoba/ffstuff/app/tgnotify/config"
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

var cfg config.Config

func init() {

	err := fmt.Errorf("config not loaded")
	cfg, err = config.Load()
	if err != nil {
		fmt.Println(err.Error())
		panic(0)
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

				token := cfg.ApiToken()
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
				chats := cfg.ChatChannels()
				if _, ok := chats[chatKey]; ok != true {
					return fmt.Errorf("no key '%v' found in config file", chatKey)
				}

				chatID, topic, err := ChatIDAndChatTopic(chatKey)
				if err != nil {
					return err
				}

				msg := tgbotapi.NewMessage(chatID, message)
				msg.ParseMode = tgbotapi.ModeHTML
				if topic > -1 {
					msg.ReplyToMessageID = int(topic)
				}

				_, errS := bot.Send(msg)
				if errS != nil {
					return fmt.Errorf("send message: %v", errS.Error())
				}

				return nil
			},
		},
		{
			Name:  "config",
			Usage: "print current config",
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
		{ //TODO
			Name:  "add",
			Usage: "add chat key to config from url",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "key",
					Usage:    "set key for new chat",
					Required: true,
					Aliases:  []string{"k"},
				},
				&cli.StringFlag{
					Name:     "link",
					Usage:    "parse chat data from value",
					Required: true,
					Aliases:  []string{"l"},
				},
			},
			Action: func(c *cli.Context) error {
				if len(c.Args().Slice()) != 0 {
					return fmt.Errorf("action 'add_chat' must not use arguments. \ncheck if key or link spaces and not encaplated with quotes" + ` (")`)
				}
				newKey := c.String("k")
				chats := cfg.ChatChannels()
				for availableKey := range chats {
					if newKey == availableKey {
						return fmt.Errorf("can't add key: '%v' already present", newKey)
					}
				}
				link := c.String("l")
				knownPrefixes := []string{`https://web.telegram.org/a/#`, `https://t.me/c/`}
				chatDataLine := ""
				for _, prefix := range knownPrefixes {
					if !strings.HasPrefix(link, prefix) {
						continue
					}
					data := strings.TrimPrefix(link, prefix)
					data = strings.ReplaceAll(data, "/", "_")

					dataParts := strings.Split(data, "_")

					if len(dataParts) > 0 {
						chatDataLine = dataParts[0]
						if !strings.HasPrefix(chatDataLine, "-100") {
							chatDataLine = "-100" + chatDataLine
						}
					}
					if len(dataParts) > 1 {

						chatDataLine += "_" + dataParts[1]
					}

				}
				if chatDataLine == "" {
					return fmt.Errorf("parsing failed\nenshure value of a flag '--link' is url of telegram chat")
				}
				cfg.ChatChannels()[newKey] = chatDataLine
				cfg.Save()

				// programConfig.ChatData[newKey] = chatDataLine
				// bts, errM := json.MarshalIndent(programConfig, "", "  ")
				// if errM != nil {
				// 	return errM
				// }
				// f, ee := os.OpenFile(configPath, os.O_WRONLY, 0777)
				// if ee != nil {
				// 	return ee
				// }
				// f.Truncate(0)
				// if _, err := f.Write(bts); err != nil {
				// 	return err
				// }
				return nil
			},
		},
		{
			Name:      "update_esplay_channel",
			Usage:     "none",
			ArgsUsage: "no arguments allowed",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "to_chat",
					Usage:    "where signal to send to",
					Required: true,
					Aliases:  []string{"tc"},
				},
			},
			Action: func(c *cli.Context) error {

				//https://t.me/c/1338947033/144079
				token := cfg.ApiToken()
				fmt.Println(token)
				bot, err := tgbotapi.NewBotAPI(token)
				if err != nil {
					fmt.Println(token)
					return fmt.Errorf("create bot api: %v", err.Error())
				}

				chatKey := c.String("to_chat")
				chats := cfg.ChatChannels()
				if _, ok := chats[chatKey]; ok != true {
					return fmt.Errorf("no key '%v' found in config file", chatKey)
				}

				chatID, _, err := ChatIDAndChatTopic(chatKey)
				if err != nil {
					return err
				}

				chatConfig := tgbotapi.NewForward()
				chatConfig.ChatID = chatID
				chatConfig.SuperGroupUsername = chatKey

				chatData, err := bot.GetChat(chatConfig)
				if err != nil {
					fmt.Println("chatErr", err.Error())
				}
				fmt.Println(chatData.PinnedMessage.Text)

				// editMsgSignal := tgbotapi.NewEditMessageText(chatID, topic, "edited message")
				// editMsgSignal.ParseMode = tgbotapi.ModeHTML

				// _, errS := bot.Send(deleteMsgSignal)
				// if errS != nil {
				// 	return fmt.Errorf("send message: %v", errS.Error())
				// }

				return nil
			},
		},
		// { //TODO
		// 	Name:  "delete",
		// 	Usage: "delete chat key from config",
		// 	Action: func(c *cli.Context) error {
		// 		chats := cfg.ChatChannels()
		// 		keys := c.Args().Slice()
		// 		if len(keys) < 1 {
		// 			return fmt.Errorf("action 'delete' uses arguments for keys")
		// 		}

		// 		userOutput := "keys deleted: 0"
		// 		deleted := 0
		// 		for _, k := range keys {
		// 			if _, ok := chats[k]; !ok {
		// 				println(fmt.Sprintf("key '%v' is not found", k))
		// 				continue
		// 			}
		// 			delete(chats, k)
		// 			deleted++
		// 			userOutput = fmt.Sprintf("keys deleted: %v", deleted)
		// 		}

		// 		bts, errM := json.MarshalIndent(programConfig, "", "  ")
		// 		if errM != nil {
		// 			return errM
		// 		}
		// 		f, ee := os.OpenFile(configPath, os.O_WRONLY, 0777)
		// 		if ee != nil {
		// 			return ee
		// 		}
		// 		f.Truncate(0)
		// 		if _, err := f.Write(bts); err != nil {
		// 			return err
		// 		}
		// 		println(userOutput)
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

func assertNoError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
