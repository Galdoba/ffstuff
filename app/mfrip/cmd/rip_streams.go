package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/ffstuff/app/mfline/ump"
	"github.com/urfave/cli/v2"
)

const (
	flag_ov = iota
	flag_oa
	flag_os
	flag_io
	flag_q
	flag_abs
	video = "video"
	audio = "audio"
)

var flag map[int]bool

func RipStreams() *cli.Command {
	flag = make(map[int]bool)
	cm := &cli.Command{
		Name:  "streams",
		Usage: "Rip streams from media file",
		Action: func(c *cli.Context) error {
			args := c.Args().Slice()
			outputFiles := []string{}
			flags(c)
			fmt.Println("check flags: done")
			for fileNum, filePath := range args {

				fmt.Print("check file: ", filePath)
				_, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				fmt.Print("   ok\n")
				// fileName := filepath.Base(filePath)
				// fileExt := filepath.Ext(filePath)
				// fileBase := strings.TrimSuffix(fileName, fileExt)

				mp := ump.NewProfile()
				err = mp.ScanBasic(filePath)
				fmt.Print("scan file: ", filePath)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				streamNumMap := make(map[string]int)
				for _, stream := range mp.Streams {
					streamNumMap[stream.Codec_type]++
				}

				fmt.Print("   ok\n")
				for _, stype := range []string{video, audio} {
					reaction := ""
					switch stype {
					case video:
						if flag[flag_ov] {
							reaction = " (skip)"
						}
					case audio:
						if flag[flag_ov] {
							reaction = " (skip)"
						}
					}
					fmt.Printf("%v: %v%v\n", stype, streamNumMap[stype], reaction)

				}
				fmt.Println("")
				// fc := fcVideo(streamNumMap[video]) + fcAudio(streamNumMap[audio])
				// if fc != "" {
				// 	fc = `"` + strings.TrimSuffix(fc, "; ") + `"`
				// }
				// mVid := mapVideo(streamNumMap[video], filePath, flag[flag_io])
				mAud := mapAudio(streamNumMap[audio], filePath, flag[flag_io])
				//vOutName := ""
				fmt.Println(filePath)
				ffmpegArgs := strings.Fields(fmt.Sprintf(`-i `+filePath+` %v`, mAud))
				fmt.Printf("|%v|\n", ffmpegArgs)

				_, stderr, err := command.Execute("ffmpeg", command.CommandLineArguments("ffmpeg", strings.Join(ffmpegArgs, " ")),
					command.Set(command.TERMINAL_ON))
				if err != nil {
					fmt.Println(err.Error())
				}
				if stderr != "" {
					fmt.Println("++")
				}
				fmt.Println(fileNum, filePath)
			}
			fmt.Println("fl1:", c.String("ov"))
			fmt.Println(outputFiles)

			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:               "ov",
				DisableDefaultText: true,
				Usage:              "omite video",
			},
			&cli.BoolFlag{
				Name:               "oa",
				DisableDefaultText: true,
				Usage:              "omite audio",
			},
			&cli.BoolFlag{
				Name:               "os",
				DisableDefaultText: true,
				Usage:              "omite subtitles",
			},
			&cli.BoolFlag{
				Name:               "index_only",
				DisableDefaultText: true,
				Aliases:            []string{"io"},
				Usage:              "new files keep only index in their name (do not reset for multiple arguments)",
			},
			&cli.BoolFlag{
				Name:               "queite",
				DisableDefaultText: true,
				Aliases:            []string{"q"},
				Usage:              "do not print new file paths",
			},
			&cli.BoolFlag{
				Name:               "absolute_paths",
				DisableDefaultText: true,
				Aliases:            []string{"abs"},
				Usage:              "new file paths are absolute",
			},
		},
	}
	return cm
}

func flags(c *cli.Context) {
	flag[flag_ov] = c.Bool("ov")
	flag[flag_oa] = c.Bool("oa")
	flag[flag_os] = c.Bool("os")
	flag[flag_io] = c.Bool("io")
	flag[flag_q] = c.Bool("q")
	flag[flag_abs] = c.Bool("abs")
}

func fcVideo(n int) string {
	fc := ""
	for i := 0; i < n; i++ {
		fc += fmt.Sprintf("[0:v:%v][v%v]; ", i, i)
	}
	return fc
}

func fcMapsNamesAUDIO(chNums []int, baseOut string) (string, []string, []string) {
	fc := ""
	cLoc := 0
	outLoc := []string{}
	outTag := []string{}
	for i, n := range chNums {
		for c := 0; c < n; c++ {
			chNumCloc := chNum(cLoc, true)
			chNumC := chNum(c, false)
			add := fmt.Sprintf("[0:a:%v]pan=mono|c0=c%v[a%v];", i, chNumC, chNumCloc)
			fc += add
			outLoc = append(outLoc, chNumCloc)
			oTag := fmt.Sprintf("a%vch%v", i, c)
			outTag = append(outTag, oTag)
			cLoc++
		}
	}
	fc = strings.TrimSuffix(fc, ";")
	maps := []string{}
	names := []string{}
	base := baseOut
	if baseOut != "" {
		base += "_"
	}
	base = strings.TrimSuffix(base, "_")
	for i, n := range outLoc {
		maps = append(maps, fmt.Sprintf("-map [a%v]", n))
		names = append(names, fmt.Sprintf("%v.%v.wav", base, outTag[i]))
	}

	return fc, maps, names
}

func chNum(n int, local bool) string {
	if flag[flag_io] {
		n = chanIndex
		if !local {
			chanIndex++
		}
	}
	s := ""
	if n < 10 {
		s += "0"
	}
	s += fmt.Sprintf("%v", n)
	return s

}

func mapVideo(n int, source string, index_only bool) string {
	mp := ""
	base, ext := baseAndExt(source)
	for i := 0; i < n; i++ {
		baseLoc := fmt.Sprintf("%v_v%v", base, streamNumber(i))
		if index_only {
			baseLoc = "v" + streamNumber(i)
		}
		outName := baseLoc + ext
		mp += fmt.Sprintf(" -map 0:v:%v -c copy %v", i, outName)
	}
	return mp
}

func mapAudio(n int, source string, index_only bool) string {
	mp := ""
	base, _ := baseAndExt(source)
	for i := 0; i < n; i++ {
		baseLoc := fmt.Sprintf("%v_a%v", base, streamNumber(i))
		if index_only {
			baseLoc = "a" + streamNumber(i)
		}
		outName := baseLoc + ".wav"
		mp += fmt.Sprintf(` -map 0:a:%v -acodec copy %v`, i, outName)
	}
	fmt.Printf("|%v|\n", mp)
	return mp
}

func streamNumber(n int) string {
	s := ""
	if n < 10 {
		s += "0"
	}
	s += fmt.Sprintf("%v", n)
	return s
}

func baseAndExt(source string) (string, string) {
	fileName := filepath.Base(source)
	fileExt := filepath.Ext(source)
	fileBase := strings.TrimSuffix(fileName, fileExt)
	return fileBase, fileExt
}
