package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/ffstuff/app/plagen/config"
	"github.com/Galdoba/ffstuff/app/plagen/internal/action"
	"github.com/urfave/cli/v2"
)

func Custom() *cli.Command {

	return &cli.Command{
		Name:        "custom",
		Aliases:     []string{},
		Usage:       "generate custom video",
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
			cfg, _ := config.Load()
			base := c.String("base")
			if cfg.VideoPaths[base] == "" {
				fmt.Printf("unknown base argument received!\n")
				baseConf := []string{}
				for b := range cfg.VideoPaths {
					baseConf = append(baseConf, b)
				}
				fmt.Printf("define this argument in config file or use one of these:\n")
				sort.Strings(baseConf)
				for _, b := range baseConf {
					fmt.Printf(" - %v\n", b)
				}
				return fmt.Errorf("bad argument")
			}
			frmt := strings.ToUpper(c.String("format"))
			switch frmt {
			default:
				return fmt.Errorf("unknown argument 'format=%v': expect '4K', 'HD', 'SD43' or 'SD169'", frmt)
			case "4K":
			case "HD":
			case "SD169":
			case "SD43":
			}
			expect := expectVideo(base, frmt)
			if !inCache(expect) {
				err := fmt.Errorf("no formated video")
				switch frmt {
				case "4K":
					_, _, err = command.Execute(fmt.Sprintf("ffmpeg -hide_banner -y -i %v -c:v libx264 -preset medium -crf 18 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 %v", cfg.VideoPaths[base], expect), command.Set(command.TERMINAL_ON))
				case "HD":
					_, _, err = command.Execute(fmt.Sprintf("ffmpeg -hide_banner -y -i %v -vf scale=1920:-2,setsar=1/1,unsharp=3:3:0.3:3:3:0,pad=1920:1080:-1:-1 -c:v libx264 -preset medium -crf 16 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 %v", cfg.VideoPaths[base], expect), command.Set(command.TERMINAL_ON))
				case "SD43":
					_, _, err = command.Execute(fmt.Sprintf("ffmpeg -hide_banner -y -i %v -vf scale=-2:576:0:0,unsharp=3:3:0.3:3:3:0,setsar=16/15,crop=720:576:152:0 -c:v libx264 -preset medium -crf 13 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 %v", cfg.VideoPaths[base], expect), command.Set(command.TERMINAL_ON))
				case "SD169":
					_, _, err = command.Execute(fmt.Sprintf("ffmpeg -hide_banner -y -i %v -vf scale=720:round(ih/(iw/1024)/2)*2,unsharp=3:3:0.3:3:3:0,pad=720:576:-1:-1,setsar=64/45 -c:v libx264 -preset medium -crf 13 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 %v", cfg.VideoPaths[base], expect), command.Set(command.TERMINAL_ON))
				}
				if err != nil {
					return fmt.Errorf("subprocess err: creating formatted video: %v", err.Error())
				}
			}

			layout := c.String("layout")
			err := fmt.Errorf("not parsed layout")
			audio := []string{}
			srtNum := 0
			audio, srtNum, err = action.ParseLayout(layout)
			if err != nil {
				return err
			}
			fmt.Println("generate command:")
			streams := 0
			out := "ffmpeg -hide_banner -y "
			out += fmt.Sprintf("-i %v ", expect)
			for _, a := range audio {
				out += fmt.Sprintf("-i %v ", cfg.AudioPaths[a])

			}
			for s := 0; s < srtNum; s++ {
				out += fmt.Sprintf("-i %v ", cfg.Subtitle)
			}
			out += "-codec copy -codec:s mov_text -map 0:v "
			for i := range audio {
				out += fmt.Sprintf("-map %v:a -metadata:s:a:%v language=und ", i+1, i)
				streams++
			}
			for s := 0; s < srtNum; s++ {
				out += fmt.Sprintf("-map %v:s -metadata:s:s:%v language=rus ", streams+1, s)
				streams++
			}
			sep := string(filepath.Separator)
			out += fmt.Sprintf("-t 5.57 %v%v%v%v_%v_%v.mp4", os.Getenv("AGELOGOPATH"), base, sep, base, frmt, makeTag(audio, srtNum))
			fmt.Println(out)
			_, _, err = command.Execute(out, command.Set(command.TERMINAL_ON))
			if err != nil {
				return fmt.Errorf("main process err: creating video: %v", err.Error())
			}
			return nil
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			return nil
		},
		Subcommands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "source",
				Aliases: []string{"s"},
			},
			&cli.StringFlag{
				Name:    "base",
				Aliases: []string{"b"},
			},
			&cli.StringFlag{
				Name:     "format",
				Required: true,
				Aliases:  []string{"f"},
			},
			&cli.StringFlag{
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

func expectVideo(base, frmt string) string {
	sep := string(filepath.Separator)
	return os.Getenv("AGELOGOPATH") + "originals" + sep + ".cache" + sep + base + "_" + frmt + ".mp4"
}

func inCache(path string) bool {
	f, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer f.Close()
	if err != nil {
		return false
	}
	return true

}

func makeTag(audio []string, srtNum int) string {
	t := "a"
	for _, a := range audio {
		t += a
	}
	if srtNum > 0 {
		t += "_s" + fmt.Sprintf("%v", srtNum)
	}
	return t
}
