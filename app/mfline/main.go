package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/gconfig"
	"github.com/Galdoba/ffstuff/app/mfline/cmd"
	"github.com/urfave/cli/v2"
)

const (
	programName    = "mfline"
	opt_storageDir = "some dir"
)

var configuration *gconfig.Config

func init() {
	conf, err := gconfig.Load(programName)
	if err != nil {
		fmt.Printf("can't initiate %v: %v\n", programName, err.Error())
		if strings.Contains(err.Error(), " The system cannot find") {
			fmt.Printf("creating default config:")
			conf, err = gconfig.NewConfig(programName, gconfig.Default())
			if err != nil {
				fmt.Printf("can't create default config: %v\n", err.Error())
				os.Exit(1)
			}

			conf.Option_STR[opt_storageDir] = "created"
			if conf.Save() != nil {
				fmt.Printf("can't create default config: %v\n", err.Error())
				os.Exit(1)
			}
			fmt.Printf("    ok\n")
			fmt.Printf("restart mfline")
			os.Exit(0)
		}
		os.Exit(1)
	}
	configuration = conf
}

func main() {
	app := cli.NewApp()

	app.Version = "v 0.1.0"
	app.Name = programName
	app.Usage = "Parse media stream data from file\nRequires ffprobe to work"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "test",
			Category:    "",
			DefaultText: "",
			FilePath:    "",
			Usage:       "",
			Required:    false,
			Hidden:      false,
			HasBeenSet:  false,
			Value:       "",
			Destination: new(string),
			Aliases:     []string{},
			EnvVars:     []string{},
			TakesFile:   false,
			Action: func(*cli.Context, string) error {
				return nil
			},
		},
	}
	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {
		//
		return nil
	}
	app.Commands = []*cli.Command{
		// {
		// 	Name:      "show",
		// 	Usage:     "Create/print universal media profile",
		// 	ArgsUsage: "",
		// 	Flags: []cli.Flag{
		// 		//-short
		// 		&cli.BoolFlag{
		// 			Name:        "short",
		// 			Category:    "Common Output:",
		// 			Usage:       "print short profile line",
		// 			Aliases:     []string{"s"},
		// 			DefaultText: "true if all Output flags false",
		// 		},
		// 		//-long
		// 		&cli.BoolFlag{
		// 			Name:     "long",
		// 			Category: "Common Output:",
		// 			Usage:    "print long profile line",
		// 			Aliases:  []string{"l"},
		// 		},
		// 		&cli.BoolFlag{
		// 			Name:     "audio_layout",
		// 			Category: "Common Output:",
		// 			Usage:    "print audio layout line",
		// 			Aliases:  []string{"a"},
		// 		},
		// 		//-warnings
		// 		&cli.BoolFlag{
		// 			Name:     "warning",
		// 			Category: "Common Output:",
		// 			Usage:    "print list of warnings",
		// 			Aliases:  []string{"w"},
		// 		},
		// 		//
		// 		&cli.BoolFlag{
		// 			Name:    "name",
		// 			Usage:   "print name of the file",
		// 			Aliases: []string{"n"},
		// 		},
		// 		&cli.BoolFlag{
		// 			Name:    "split",
		// 			Usage:   "print separation line between different files",
		// 			Aliases: []string{"sp"},
		// 		},
		// 		&cli.StringSliceFlag{
		// 			Name:     "include_stream",
		// 			Category: "Filtered Output:",
		// 			//Usage:    "usage: \n      print info on exact stream if available\n      uses ffmpeg notation: '0:v:0', '0:a:2'\n      valid examples: '0:a:1', '[0:a:1]', 'a1', 4 (this will return info on stream number 4)\n      special case: 'all' - return all stream info\n      special case: '0' - return list of all available keys to stdout",
		// 			Usage: strings.Join([]string{
		// 				"usage:",
		// 				"print info on exact stream if available",
		// 				"uses ffmpeg notation: '0:v:0', '0:a:2'",
		// 				"key examples: '0:a:1', '[0:a:1]', 'a1', '4' (return info on stream 4)",
		// 				"special case: 'all' - return all stream info",
		// 				"special case: '0'   - return list of all available keys to stdout",
		// 			}, "\n      "),
		// 			Aliases: []string{"is"},
		// 		},
		// 		&cli.BoolFlag{},
		// 	},
		// 	Action: func(c *cli.Context) error {
		// 		args := c.Args().Slice()
		// 		if len(args) < 1 {
		// 			return fmt.Errorf("no arguments received\n'mfline --help show' for instructions")
		// 		}
		// 		srt := c.Bool("short")
		// 		lng := c.Bool("long")
		// 		aud := c.Bool("audio_layout")
		// 		wrn := c.Bool("warning")
		// 		nme := c.Bool("name")
		// 		split := c.Bool("split")
		// 		stream_keys := c.StringSlice("include_stream")
		// 		if !srt && !lng && !wrn && !aud && len(stream_keys) == 0 {
		// 			srt = true
		// 		}
		// 		for _, path := range args {
		// 			scan := ump.NewProfile()
		// 			switch strings.HasSuffix(path, ".json") {
		// 			case true:
		// 				err := scan.ConsumeJSON(path)
		// 				if err != nil {
		// 					fmt.Fprintf(os.Stderr, "can't consume json: %v\n", err.Error())
		// 					continue
		// 				}
		// 			default:
		// 				err := scan.ConsumeFile(path)
		// 				if err != nil {
		// 					fmt.Fprintf(os.Stderr, "can't consume file: %v\n", err.Error())
		// 					continue
		// 				}
		// 			}

		// 			if split {
		// 				fmt.Fprintf(os.Stdout, "\n")
		// 			}
		// 			if nme {
		// 				fmt.Fprintf(os.Stdout, "%v\n", path)
		// 			}
		// 			if srt {
		// 				fmt.Fprintf(os.Stdout, "%v\n", scan.Short())
		// 			}
		// 			if lng {
		// 				fmt.Fprintf(os.Stdout, "%v\n", scan.Long())
		// 			}
		// 			if aud {
		// 				fmt.Fprintf(os.Stdout, "%v\n", scan.AudioLayout())
		// 			}
		// 			if wrn {
		// 				for _, w := range scan.Warnings() {
		// 					fmt.Fprintf(os.Stdout, "%v\n", w)
		// 				}
		// 			}
		// 			if len(stream_keys) > 0 {
		// 				info := scan.ByStream()
		// 				errors := []string{}
		// 				printed := 0

		// 				for _, key := range stream_keys {
		// 					if val, ok := info[key]; ok {
		// 						fmt.Fprintf(os.Stdout, "%v\n", val)
		// 						printed++
		// 					} else {
		// 						errors = append(errors, fmt.Sprintf("error: no data on key '%v'", key))
		// 					}
		// 				}
		// 				for _, err := range errors {
		// 					fmt.Fprintf(os.Stderr, "%v\n", err)
		// 				}
		// 				if len(errors) > 0 && printed == 0 {
		// 					err := fmt.Sprintf("possible keys: ")

		// 					for k := range info {
		// 						err += fmt.Sprintf("'%v', ", k)
		// 					}
		// 					err = strings.TrimSuffix(err, ", ")
		// 					fmt.Fprintf(os.Stderr, "%v\n", err)
		// 				}
		// 				//fmt.Fprintf(os.Stdout, "DEBUG: %v\n", info)
		// 			}
		// 		}

		// 		return nil
		// 	},
		// },
		cmd.Show(),
		cmd.ScanStreams(),
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
