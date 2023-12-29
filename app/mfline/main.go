package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/app/mfline/ump"
	"github.com/urfave/cli/v2"
)

const (
	programName = "mfline"
)

//сканировать файл
//сканировать json
//
//mfprofile -q show -file file.mp4 -long -save file.json
//mfprofile show -json file.json

func main() {
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = programName
	app.Usage = "Parse media stream data from file\nRequires ffprobe to work"
	app.Flags = []cli.Flag{}
	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {
		//
		return nil
	}
	app.Commands = []*cli.Command{
		{
			Name:      "show",
			Usage:     "Create/print universal media profile",
			ArgsUsage: "",
			Flags: []cli.Flag{
				//-short
				&cli.BoolFlag{
					Name:        "short",
					Category:    "Output",
					Usage:       "print short profile line",
					Aliases:     []string{"s"},
					DefaultText: "true if all Output flags false",
				},
				//-long
				&cli.BoolFlag{
					Name:     "long",
					Category: "Output",
					Usage:    "print long profile line",
					Aliases:  []string{"l"},
				},
				&cli.BoolFlag{
					Name:     "audio_layout",
					Category: "Output",
					Usage:    "print audio layout line",
					Aliases:  []string{"a"},
				},
				//-warnings
				&cli.BoolFlag{
					Name:     "warning",
					Category: "Output",
					Usage:    "print list of warnings",
					Aliases:  []string{"w"},
				},
				//
				&cli.BoolFlag{
					Name:    "name",
					Usage:   "print name of the file",
					Aliases: []string{"n"},
				},
				&cli.BoolFlag{
					Name:    "split",
					Usage:   "print separation line between different files",
					Aliases: []string{"sp"},
				},
			},
			Action: func(c *cli.Context) error {
				args := c.Args().Slice()
				if len(args) < 1 {
					return fmt.Errorf("no arguments received\n'mfline --help show' for instructions")
				}
				srt := c.Bool("short")
				lng := c.Bool("long")
				aud := c.Bool("audio_layout")
				wrn := c.Bool("warning")
				nme := c.Bool("name")
				split := c.Bool("split")
				if !srt && !lng && !wrn && !aud {
					srt = true
				}
				for _, path := range args {
					scan := ump.NewProfile()
					switch strings.HasSuffix(path, ".json") {
					case true:
						err := scan.ConsumeJSON(path)
						if err != nil {
							fmt.Fprintf(os.Stderr, "can't consume json: %v\n", err.Error())
							continue
						}
					default:
						err := scan.ConsumeFile(path)
						if err != nil {
							fmt.Fprintf(os.Stderr, "can't consume file: %v\n", err.Error())
							continue
						}
					}

					if split {
						fmt.Fprintf(os.Stdout, "\n")
					}
					if nme {
						fmt.Fprintf(os.Stdout, "%v\n", path)
					}
					if srt {
						fmt.Fprintf(os.Stdout, "%v\n", scan.Short())
					}
					if lng {
						fmt.Fprintf(os.Stdout, "%v\n", scan.Long())
					}
					if aud {
						fmt.Fprintf(os.Stdout, "%v\n", scan.AudioLayout())
					}
					if wrn {
						for _, w := range scan.Warnings() {
							fmt.Fprintf(os.Stdout, "%v\n", w)
						}
					}
				}

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
