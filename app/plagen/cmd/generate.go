package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/app/plagen/config"
	"github.com/Galdoba/ffstuff/app/plagen/internal/action"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func Custom() *cli.Command {

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
			fmt.Println(cfg.Destination)
			tasks := []*action.GenerationTask{}
			//CHECK OPTION

			dest := ""
			frmt := c.String("format")
			width := 0
			height := 0
			dur := c.Int("duration")
			switch frmt {
			default:
				return fmt.Errorf("unknown value 'format=%v': expect '4K', 'HD' or 'SD'", frmt)
			case "4K":
				width = 3140
				height = 2160
			case "HD":
				width = 1920
				height = 1080
			case "SD":
				width = 720
				height = 576
			}
			if c.String("target") == "" {
				dest = action.StdDestinationDir("undefined")
			}

			//ВЫБИРАЕМ ОСНОВУ
			source := c.String("s")

			_, err := os.Stat(source)
			if err != nil {
				switch {
				case os.IsNotExist(err):
					files := action.DetectSources()
					for _, fname := range files {
						base := filepath.Base(fname)
						if strings.HasPrefix(base, source+"__") {
							source = fname
							err = nil
							break
						}
					}
					if err != nil {
						return fmt.Errorf("source '%v' have bad name is not exist", source)
					}
				}
			}
			fmt.Println(source)

			validation := action.SourceValid(source, width, height, dur)
			if !validation.Valid {
				return fmt.Errorf("source '%v' invalid: %v", source, validation.MSG)
			}
			base := filepath.Base(source)
			prts := strings.Split(base, "__")
			if len(prts) > 1 {
				dest = prts[0] + "/"
			}
			layout := c.String("layout")
			fmt.Printf("|%v|\n", layout)
			layout = strings.TrimPrefix(layout, "[")
			layout = strings.TrimSuffix(layout, "]")
			layouts := strings.Fields(layout)
			fmt.Println(root + dest)
			for _, parsed := range layouts {
				audio, srtNum, err := action.ParseLayout(parsed)
				if err != nil {
					return err
				}
				tasks = append(tasks, action.NewTask(source, audio, srtNum))
			}
			for _, tsk := range tasks {
				fmt.Println(tsk)
			}
			// res, _ := Select("Select video source", options...)
			// videoSource := cfg.VideoPaths[res]
			// fmt.Printf("video source: %v\n", videoSource)
			// //ВЫБИРАЕМ ФОРМАТ

			// options = []string{}
			// for _, k := range cfg.VideoFormats {
			// 	options = append(options, k)
			// }
			// res, _ = Select("Select video format", options...)
			// videoFormat := res
			// fmt.Printf("video format: %v\n", videoFormat)
			// //ВЫБИРАЕМ КОЛВО звука и языки
			// fmt.Println("Emulate stage 3")
			// langNum, err := InputInt()
			// fmt.Println(langNum, err, "-------------")

			// fmt.Println("Emulate stage 4")
			// //считаем
			// fmt.Println("Emulate stage 5")

			fmt.Println("generate END")
			return nil
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			return nil
		},
		Subcommands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "source",
				Required: true,
				Aliases:  []string{"s"},
			},

			&cli.StringFlag{
				Name:     "target",
				Required: true,
				Aliases:  []string{"t"},
			},
			&cli.StringFlag{
				Name:     "format",
				Required: true,
				Aliases:  []string{"f"},
			},

			&cli.IntFlag{
				Name:    "duration",
				Aliases: []string{"d"},
			},

			&cli.StringSliceFlag{
				Name:     "layout",
				Required: true,
				Aliases:  []string{"l"},
			},
		},
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

// func Input() (string, error) {
// 	hint := ""
// 	hint += "Справка: код каналов состоит из 2 частей с разделителем в виде '_'\n"
// 	hint += "- первая часть это последовательность из кол-ва каналов в стриме. \n"
// 	hint += "- вторая часть отвечает за количество субтитров. В случае их отсуствия не вводится.\n"
// 	hint += "Примеры:\n"
// 	hint += "	62_s1 - 2 аудио (первая 5.1, вторая стерео) + 1 субитры\n"
// 	hint += "662   - 2 аудио (первые две - 5.1, третья - стерео) без субтитров\n"
// 	hint += "6262_s2   - 4 аудио (5.1, стерео,5.1, стерео) + 2 комплекта субтитров"

// 	fmt.Println(hint)
// 	val := -1
// 	err := fmt.Errorf("not validated")
// 	validate := func(input string) error {
// 		val, err = strconv.Atoi(input)
// 		if err != nil {
// 			return errors.New("Invalid number")
// 		}
// 		if val < 0 {
// 			return errors.New("Ожидаем число больше нуля")
// 		}
// 		return nil
// 	}

// 	prompt := promptui.Prompt{
// 		Label:    "Введи код каналов:",
// 		Default:  "",
// 		Validate: validate,
// 		Templates: &promptui.PromptTemplates{
// 			Success: "",
// 		},
// 	}

// 	_, err = prompt.Run()

// 	if err != nil {
// 		fmt.Printf("Prompt failed %v\n", err)
// 		return -1, err
// 	}

// 	// fmt.Printf("You choose %q\n", result)
// 	return val, err
// }
