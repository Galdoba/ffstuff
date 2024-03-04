package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Galdoba/ffstuff/app/plagen/config"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func Generate() *cli.Command {

	return &cli.Command{
		Name:        "custom",
		Aliases:     []string{},
		Usage:       "sync scan data files is original files get prefix (auto feature)",
		UsageText:   "",
		Description: "",
		Args:        false,
		ArgsUsage:   "",
		Category:    "",
		BashComplete: func(*cli.Context) {
		},
		Before: func(*cli.Context) error {
			return nil
		},
		After: func(*cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			cfg, _ := config.Load(c.App.Name)
			//ВЫБИРАЕМ ОСНОВУ
			options := []string{}
			for k := range cfg.VideoPaths {
				options = append(options, k)
			}
			res, _ := Select("Select video source", options...)
			videoSource := cfg.VideoPaths[res]
			fmt.Printf("video source: %v\n", videoSource)
			//ВЫБИРАЕМ ФОРМАТ

			options = []string{}
			for _, k := range cfg.VideoFormats {
				options = append(options, k)
			}
			res, _ = Select("Select video format", options...)
			videoFormat := res
			fmt.Printf("video format: %v\n", videoFormat)
			//ВЫБИРАЕМ КОЛВО звука и языки
			fmt.Println("Emulate stage 3")
			langNum, err := InputInt()
			fmt.Println(langNum, err, "-------------")

			fmt.Println("Emulate stage 4")
			//считаем
			fmt.Println("Emulate stage 5")

			fmt.Println("generate END")
			return nil
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			return nil
		},
		Subcommands:            []*cli.Command{},
		Flags:                  []cli.Flag{},
		SkipFlagParsing:        false,
		HideHelp:               false,
		HideHelpCommand:        false,
		Hidden:                 true,
		UseShortOptionHandling: false,
		HelpName:               "",
		CustomHelpTemplate:     "",
	}
}

func Select(msg string, opts ...string) (string, error) {
	prompt := promptui.Select{
		Label: msg,
		Items: opts,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "?", err
	}

	fmt.Printf("You choose %q\n", result)
	return result, nil
}

func InputInt() (int, error) {
	val := -1
	err := fmt.Errorf("not validated")
	validate := func(input string) error {
		val, err = strconv.Atoi(input)
		if err != nil {
			return errors.New("Invalid number")
		}
		if val < 0 {
			return errors.New("Ожидаем число больше нуля")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Введи целоположительное число",
		Default:  "",
		Validate: validate,
		Templates: &promptui.PromptTemplates{
			Success: "",
		},
	}

	_, err = prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1, err
	}

	// fmt.Printf("You choose %q\n", result)
	return val, err
}
