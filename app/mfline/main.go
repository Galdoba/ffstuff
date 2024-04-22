package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/app/mfline/cmd"
	"github.com/Galdoba/ffstuff/app/mfline/config"
	"github.com/urfave/cli/v2"
)

const (
	CONFIG = "cfg"
)

func main() {
	app := cli.NewApp()

	app.Version = "v 0.1.3"
	app.Usage = "Parse media stream data from file\nRequires ffprobe to work"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:      "use_config",
			Usage:     "use alternative config",
			TakesFile: false,
			Action: func(*cli.Context, string) error {
				return nil
			},
		},
	}

	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {
		cfg := &config.Config{}
		err := fmt.Errorf("config not loaded")
		cfg, err = config.Load(app.Name)
		if err != nil {
			fmt.Println("config err:", err.Error())

			switch {
			case strings.Contains(err.Error(), "The system cannot find "), strings.Contains(err.Error(), "no such file or directory"):
				cfg, err = config.NewConfig(c.App.Name)
				if err = cfg.Save(); err != nil {
					return fmt.Errorf("can't setup config: %v", err.Error())
				}
				fmt.Printf("default config created at %v: restart %v\n", cfg.Location, cfg.AppName)
				os.Exit(0)
			default:
				panic("wut?")
				return nil
			}
		}
		if _, err := os.ReadDir(cfg.StorageDir); err != nil {
			return fmt.Errorf("can't read storage dir: %v", err.Error())
		}

		if cfg.WriteLogs {
			f, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				return fmt.Errorf("can't write to log file: %v", err.Error())
			}
			defer f.Close()
		}
		for _, dir := range cfg.TrackDirs {
			if dir == "[TEMPLATE]" {
				continue
			}
			if _, err := os.ReadDir(dir); err != nil {
				return fmt.Errorf("can't read tracked dir: %v", err.Error())
			}
		}
		return nil
	}
	app.Commands = []*cli.Command{
		cmd.Sync(),
		cmd.Config(),
		cmd.Show(),
		cmd.ScanStreams(),
		cmd.FullScan(),
		cmd.AudioStreamsData(),
	}

	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		return nil
	}
	args := os.Args
	if err := app.Run(args); err != nil {
		errOut := fmt.Sprintf("%v error: %v", app.Name, err.Error())
		println(errOut)
		os.Exit(1)
	}

}
