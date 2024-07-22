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

var chanIndex int

func RipChannels() *cli.Command {
	flag = make(map[int]bool)
	cm := &cli.Command{
		Name: "channels",
		Action: func(c *cli.Context) error {
			args := c.Args().Slice()
			outputFiles := []string{}
			flags(c)
			fmt.Println("check flags: done")
			for fileNum, filePath := range args {
				fmt.Println("start", fileNum)
				chanNum := []int{}
				fmt.Print("check file: ", filePath)
				_, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				fmt.Print("   ok\n")
				mp := ump.NewProfile()
				err = mp.ScanBasic(filePath)
				fmt.Print("scan file: ", filePath)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				streamNumMap := make(map[string]int)
				for _, stream := range mp.Streams {
					if stream.Codec_type != audio {
						continue
					}
					streamNumMap[stream.Codec_type]++
					chanNum = append(chanNum, stream.Channels)
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
				base := filepath.Base(filePath)
				fc, maps, names := fcMapsNamesAUDIO(chanNum, base)
				if fc != "" {
					fc = strings.TrimSuffix(fc, ";")
				}
				ffmpegArgs := strings.Fields(fmt.Sprintf("-i " + filePath + " -filter_complex " + fc))
				for i := range maps {
					ffmpegArgs = append(ffmpegArgs, maps[i])
					ffmpegArgs = append(ffmpegArgs, names[i])
					outputFiles = append(outputFiles, names[i])
				}

				// ffmpegArgs := strings.Fields(fmt.Sprintf("-i "+filePath+" -filter_complex "+fc+" %v", maps))
				_, stderr, err := command.Execute("ffmpeg", command.CommandLineArguments("ffmpeg", strings.Join(ffmpegArgs, " ")),
					command.Set(command.TERMINAL_ON))
				if err != nil {
					fmt.Println(err.Error())
				}
				if stderr != "" {
					fmt.Println("++", stderr)
				}
				fmt.Println(outputFiles)
			}

			return nil
		},
		Flags: []cli.Flag{
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

// func fcAudioPan(n int) string {
// 	//[0:a:0]pan=mono|c0=c0[left]
// 	fc := ""
// 	for i := 0; i < n; i++ {
// 		fc += fmt.Sprintf("[0:a:%v]pan=mono|c0|c%v[a%v]; ", i, i, i)
// 	}
// 	return fc
// }

// func mapAudioChan(n int, source string, index_only bool) string {
// 	mp := ""
// 	base, _ := baseAndExt(source)
// 	for i := 0; i < n; i++ {
// 		baseLoc := fmt.Sprintf("%v_a%v", base, streamNumber(i))
// 		if index_only {
// 			baseLoc = "a" + streamNumber(i)
// 		}
// 		outName := baseLoc + ".wav"
// 		mp += fmt.Sprintf(` -map 0:a:%v -acodec copy %v`, i, outName)
// 	}
// 	fmt.Printf("|%v|\n", mp)
// 	return mp
// }

// func chanNumber(n int) string {
// 	switch flag[flag_io] {
// 	case true:
// 		n = chanIndex
// 		chanIndex++
// 	}
// 	s := ""
// 	if n < 10 {
// 		s += "0"
// 	}
// 	s += fmt.Sprintf("%v", n)
// 	return s
// }
