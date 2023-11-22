package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/translit"
	"github.com/urfave/cli/v2"
	"golang.design/x/clipboard"
)

var programName string
var configPath string
var programConfig *config

type config struct {
	ReadMe           string   `json:"ReadMe",omitempty`
	ListSeparatorArg string   `json:"ListSeparatorArg"`
	Rules            []string `json:"Rules"`
	Verbose          bool     `json:"Verbose Mode"`
	LogFile          string   `json:"LogFile"`
	LogLenght        int      `json:"LogLenght"`
}

func defaultConfig() *config {
	cfg := config{}
	cfg.ListSeparatorArg = "======"
	cfg.Rules = []string{"Rule_1"}
	cfg.Verbose = true
	cfg.LogFile = ""
	return &cfg
}

func init() {
	programName = "glisort"
	configPath = defineConfigPath()
	data, err := os.ReadFile(configPath)
	fmt.Println(configPath)

	if err != nil {
		switch {
		default:
			fmt.Println("Неизвестная ошибка при проверки наличия конфига:")
			println(err.Error())
			panic(0)
		case strings.Contains(err.Error(), "The system cannot find the file specified"), strings.Contains(err.Error(), "The system cannot find the path specified"):
			fmt.Println("Config file not found")
			err := os.MkdirAll(strings.TrimSuffix(configPath, "glisort.json"), 0777)
			if err != nil {
				panic(err.Error())
			}
			programConfig = defaultConfig()
			data, err = json.MarshalIndent(programConfig, "", "  ")
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Creating default config at %v: ", configPath)
			f, err := os.Create(configPath)
			if err != nil {
				panic(err.Error())
			}
			defer f.Close()
			_, err = f.Write(data)
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("ok")
		}
	}

	err = json.Unmarshal(data, &programConfig)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = programName
	app.Usage = "сортирует список используя указаное правило"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:     "verbose",
			Category: "Global",
			Usage:    "send operation info to StdOut",
			Value:    false,
			Aliases:  []string{"vb"},
			Action: func(*cli.Context, bool) error {
				return nil
			},
		},
	}
	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {

		return nil
	}
	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		return nil
	}
	app.Commands = []*cli.Command{
		{
			Name:  "run",
			Usage: "Sort list",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "method",
					Usage:    "REQUIRED: Choose sorting rule name",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "clip",
					Usage:   "Send output to clipboard",
					Aliases: []string{"c"},
				},
				&cli.StringFlag{
					Name:  "file",
					Usage: "Send output to file (rewrite if file exists)",
				},
				&cli.StringFlag{
					Name:  "file_append",
					Usage: "--                  (append to file if file exists)",
				},
				&cli.StringFlag{
					Name:  "file_safe",
					Usage: "--                  (will not if file file exists)",
				},
				&cli.BoolFlag{
					Name:    "error",
					Usage:   "Send output to StdErr",
					Aliases: []string{"e"},
				},

				&cli.BoolFlag{
					Name:    "silent",
					Usage:   "DO NOT send output to StdOut",
					Aliases: []string{"s"},
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("EMULATE Picking a method")
				name := c.String("method")
				method := NewSortMethod(name)
				err := method.Compile(name)
				if err != nil {
					return err
				}
				fmt.Println("EMULATE SORTING")
				fmt.Println(method)
				args := c.Args().Slice()
				fmt.Println("INPUT:")
				inputNum := 1
				skipArgs := []string{}
				for _, inp := range args {
					path, err := filepath.Abs(inp)
					if err != nil {
						continue
					}
					bt, err := os.ReadFile(path)
					if err != nil {
						continue
					}
					data := strings.Split(string(bt), "\n")
					input_key := fmt.Sprintf("%v", inputNum)
					for _, d := range data {
						method.input[input_key] = append(method.input[input_key], d)
						//						fmt.Println("add", d, "to", input_key)
					}

					inputNum++
					skipArgs = append(skipArgs, inp)
				}
			argums:
				for _, inp := range args {
					for _, check := range skipArgs {
						if check == inp {
							fmt.Println("skip", check)
							continue argums
						}
					}
					fmt.Println(inp)
					if inp == method.ListSeparator {
						inputNum++
						continue
					}
					input_key := fmt.Sprintf("%v", inputNum)
					method.input[input_key] = append(method.input[input_key], inp)
					fmt.Println("add", inp, "to", input_key)
				}
				if err := method.Execute(); err != nil {
					fmt.Println("method Err", err.Error())
				}
				fmt.Println("OUTPUT:")
				for k, v := range method.input {
					fmt.Println(k, v)
				}
				//panic("TODO: нужно протестировать ввести работу с аргументами, написать экшен фильтрацтт по белому/черному списку")
				return nil
			},
		},
	}

	args0 := os.Args
	if err := app.Run(args0); err != nil {
		fmt.Printf("\napplication returned error: %v\n", err.Error())
	}
	os.Exit(3)
	// exit := ""
	clp := clipboard.Read(clipboard.FmtText)
	text := string(clp)
	args := strings.Split(text, "\n")
	list1 := []string{}
	list2 := []string{}
	delimPassed := false
	for _, arg := range args {
		switch delimPassed {
		case false:
			if arg == "======" {
				delimPassed = true
				continue
			}
			list1 = append(list1, arg)
		case true:
			list2 = append(list2, arg)
		}

	}
	// fmt.Println("l1", list1)
	// fmt.Println("l2", list2)

	/*
		{1} REGEXP "[0-3]+" ==> {1}

	*/

	//os.Exit(0)

	// list2 := []string{
	// 	`hd_2008_gollivudskaya_madam__ar2e2_d2311011913_prt231031232708_xbeuS7MSVk5_film.mp4`,
	// 	`hd_2017_energeticheskaya_revolyuciya_segodnya__ar6e2_sr_d2311011913_prt231031234451_xvCzHO0uRzA_film.mp4`,
	// 	`hd_2017_intervyu_s_dzhordzhem_martinom__ar2_d2311011913_prt231031225231_xUlXPdNi9G7_film.mp4`,
	// 	`hd_2019_igra_prestolov_posledniy_dozor__ar6e2_sr_d2311011913_prt231101004159_x7GiGvOvQAb_film.mp4`,
	// 	`hd_2020_hvatit_molchat__ar2e2_sr_d2311011913_prt231031230857_xEoLgJExEpk_film.mp4`,
	// }

	// list1 := []string{
	// 	`Голливудская мадам (Замена)`,
	// 	`Игра Престолов: Последний дозор (Замена)`,
	// 	`Интервью с Джорджем Мартином (Замена)`,
	// 	`Хватит молчать! (Замена)`,
	// 	`Энергетическая революция сегодня (Замена)`,
	// }
	output := []string{}
	for _, input := range list1 {
		tag := haveTag(input)
		tag_trnsl := strings.ToLower(translit.Transliterate(tag))
		nameTranslited := translit.Transliterate(input)
		nameTranslited = strings.TrimSuffix(nameTranslited, tag_trnsl)
		nameTranslited = strings.TrimSuffix(nameTranslited, "_")
		//fmt.Println("debug: ")
		//fmt.Println("debug: ", input, nameTranslited, ":::", tag)
		//fmt.Println(input)

		want := []string{}
		for _, check := range list2 {
			lowCheck := strings.ToLower(check)
			lowName := strings.ToLower(nameTranslited)
			if strings.Contains(lowCheck, lowName) { //&& strings.HasPrefix(lowCheck, tag_trnsl) {
				want = append(want, check)
			}
		}
		for _, w := range want {
			fmt.Println(w)
		}
		output = append(output, want...)
	}
	out := strings.Join(output, "\n")
	clipboard.Write(clipboard.FmtText, []byte(out))
}

func haveTag(name string) string {
	if strings.HasSuffix(name, " SD") {
		return "SD"
	}
	if strings.HasSuffix(name, " 4K") {
		return "4K"
	}
	if strings.HasSuffix(name, " 3D") {
		return "3D"
	}
	return "HD"
}

func defineConfigPath() string {
	userdir, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	sep := string(filepath.Separator)

	return fmt.Sprintf("%v%v.config%v%v%vglisort.json", userdir, sep, sep, programName, sep)
}

/*
FILTERS
FullMatch
PartialMatch


White Part Reg : [0-9]
Black Part Reg : "D"
FindBase




REVERSE

*/
