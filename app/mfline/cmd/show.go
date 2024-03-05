package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Galdoba/ffstuff/app/mfline/config"
	"github.com/Galdoba/ffstuff/app/mfline/ump"
	"github.com/urfave/cli/v2"
)

func Show() *cli.Command {
	cfg := &config.Config{}
	return &cli.Command{
		Name:      "show",
		Usage:     "Print/format basic scan level universal media profile",
		ArgsUsage: "",
		Flags: []cli.Flag{
			//-short
			&cli.BoolFlag{
				Name:        "short",
				Category:    "Common Output:",
				Usage:       "print short profile line",
				Aliases:     []string{"s"},
				DefaultText: "true if all Output flags false",
			},
			//-long
			&cli.BoolFlag{
				Name:     "long",
				Category: "Common Output:",
				Usage:    "print long profile line",
				Aliases:  []string{"l"},
			},
			&cli.BoolFlag{
				Name:     "audio_layout",
				Category: "Common Output:",
				Usage:    "print audio layout line",
				Aliases:  []string{"a"},
			},
			//-warnings
			&cli.BoolFlag{
				Name:     "warning",
				Category: "Common Output:",
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
			&cli.BoolFlag{
				Name:    "duration",
				Usage:   "print duration in seconds",
				Aliases: []string{"d"},
			},
			&cli.StringSliceFlag{
				Name:     "include_stream",
				Category: "Filtered Output:",
				//Usage:    "usage: \n      print info on exact stream if available\n      uses ffmpeg notation: '0:v:0', '0:a:2'\n      valid examples: '0:a:1', '[0:a:1]', 'a1', 4 (this will return info on stream number 4)\n      special case: 'all' - return all stream info\n      special case: '0' - return list of all available keys to stdout",
				Usage: strings.Join([]string{
					"usage:",
					"print info on exact stream if available",
					"uses ffmpeg notation: '0:v:0', '0:a:2'",
					"key examples: '0:a:1', '[0:a:1]', 'a1', '4' (return info on stream 4)",
					"special case: 'all' - return all stream info",
					"special case: '0'   - return list of all available keys to stdout",
				}, "\n      "),
				Aliases: []string{"is"},
			},
			&cli.BoolFlag{},
		},
		Before: func(c *cli.Context) error {
			cfg, _ = config.Load(c.App.Name)
			return nil
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
			dur := c.Bool("duration")
			split := c.Bool("split")
			stream_keys := c.StringSlice("include_stream")
			if !srt && !lng && !wrn && !aud && !dur && len(stream_keys) == 0 {
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
					storage := cfg.StorageDir
					fs, err := os.ReadDir(storage)
					if err != nil {
						return fmt.Errorf("can't read storage directory")
					}
					fname := filepath.Base(path)
					foundJSON := false
					for _, f := range fs {
						if f.IsDir() {
							continue
						}
						if f.Name() != fname+".json" {
							continue
						}
						err := scan.ConsumeJSON(storage + fname + ".json")
						if err != nil {
							fmt.Fprintf(os.Stderr, "can't consume json: %v\n", err.Error())
							continue
						}
						foundJSON = true
						break
					}
					if !foundJSON {
						err := scan.ConsumeFile(path)
						if err != nil {
							fmt.Fprintf(os.Stderr, "can't consume file: %v\n", err.Error())
							continue
						}
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
				if dur {
					fmt.Fprintf(os.Stdout, "%v", scan.Format.Duration)
				}
				if wrn {
					for _, w := range scan.Warnings() {
						fmt.Fprintf(os.Stdout, "%v\n", w)
					}
				}
				if len(stream_keys) > 0 {
					info := scan.ByStream()
					errors := []string{}
					printed := 0

					if err := specialCases(stream_keys, info); err != nil {
						if strings.Contains(err.Error(), "complete") {
							return nil
						}
						return err
					}

					// if sliceHas(stream_keys, "0") {
					// 	foundKeys := []string{}
					// 	for k := range info {
					// 		foundKeys = append(foundKeys, k)
					// 	}
					// 	sortedKeys, sortErr := sortStreamKeys(foundKeys)
					// 	if sortErr != nil {
					// 		return sortErr
					// 	}
					// 	l := len(sortedKeys)
					// 	switch l {
					// 	default:
					// 		fmt.Fprintf(os.Stderr, "stream keys detected: %v\n", l)
					// 	case 0:
					// 		fmt.Fprintf(os.Stderr, "stream keys detected: %v", l)
					// 		return fmt.Errorf("have no media streams")
					// 	}
					// 	fmt.Fprintf(os.Stdout, "%v", strings.Join(foundKeys, " "))
					// 	return nil
					// }

					for _, key := range stream_keys {
						if val, ok := info[key]; ok {
							fmt.Fprintf(os.Stdout, "%v\n", val)
							printed++
						} else {
							errors = append(errors, fmt.Sprintf("error: no data on key '%v'", key))
						}
					}
					for _, err := range errors {
						fmt.Fprintf(os.Stderr, "%v\n", err)
					}
					if len(errors) > 0 && printed == 0 {
						err := "possible keys: "
						stream_keys = []string{}
						for k := range info {
							stream_keys = append(stream_keys, k)
						}
						sort.Strings(stream_keys)
						err += strings.Join(stream_keys, " ")
						fmt.Fprintf(os.Stderr, "%v\n", err)
					}
					//fmt.Fprintf(os.Stdout, "DEBUG: %v\n", info)
				}
			}

			return nil
		},
	}

}

// func sliceHas(sl []string, s string) bool {
// 	for _, st := range sl {
// 		if st == s {
// 			return true
// 		}
// 	}
// 	return false
// }

func sortStreamKeys(keys []string) ([]string, error) {
	video := []string{}
	audio := []string{}
	data := []string{}
	subtitle := []string{}
	out := []string{}
	for _, k := range keys {
		switch {
		default:
			return nil, fmt.Errorf("can't sort keys: unknown key detected: '%v'", k)
		case strings.Contains(k, ":v:"):
			video = append(video, k)
		case strings.Contains(k, ":a:"):
			audio = append(audio, k)
		case strings.Contains(k, ":d:"):
			data = append(data, k)
		case strings.Contains(k, ":s:"):
			subtitle = append(subtitle, k)
		}
	}
	sort.Strings(video)
	sort.Strings(audio)
	sort.Strings(data)
	sort.Strings(subtitle)
	out = append(out, video...)
	out = append(out, audio...)
	out = append(out, data...)
	out = append(out, subtitle...)
	return out, nil
}

func specialCases(keys []string, info map[string]string) error {
	for _, k := range keys {
		switch k {
		case "0":
			return specialCase0(info)
		case strings.ToLower("all"):
			return specialCaseAll(info)
		}
	}
	return nil
}

func specialCase0(info map[string]string) error {
	foundKeys := []string{}
	for k := range info {
		foundKeys = append(foundKeys, k)
	}
	sortedKeys, sortErr := sortStreamKeys(foundKeys)
	if sortErr != nil {
		return sortErr
	}
	l := len(sortedKeys)
	switch l {
	default:
		fmt.Fprintf(os.Stderr, "stream keys detected: %v\n", l)
	case 0:
		fmt.Fprintf(os.Stderr, "stream keys detected: %v", l)
		return fmt.Errorf("have no media streams")
	}
	fmt.Fprintf(os.Stdout, "%v", strings.Join(foundKeys, " "))
	return fmt.Errorf("special case '0' complete")
}

func specialCaseAll(info map[string]string) error {
	foundKeys := []string{}
	for k := range info {
		foundKeys = append(foundKeys, k)
	}
	sortedKeys, sortErr := sortStreamKeys(foundKeys)
	if sortErr != nil {
		return sortErr
	}
	for _, k := range sortedKeys {
		if data, ok := info[k]; ok {
			fmt.Fprintf(os.Stdout, "%v\n", data)
		} else {
			fmt.Fprintf(os.Stderr, "unknown key [%v]\n", k)
			return fmt.Errorf("unknown stream key detected")
		}
	}
	return fmt.Errorf("special case 'all' complete")
}
