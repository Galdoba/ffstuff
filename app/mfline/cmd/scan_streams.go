package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/ffstuff/app/mfline/config"
	"github.com/Galdoba/ffstuff/app/mfline/ump"
	"github.com/urfave/cli/v2"
)

func ScanStreams() *cli.Command {
	cfg := &config.Config{}
	cm := &cli.Command{
		Name:      "scan",
		Usage:     "scan file ('mfline scan --help' for discription)",
		UsageText: "mfline scan subcommand --source PATH1 --destination PATH2 [flags]...\nno arguments expected",
		Description: strings.Join([]string{
			"scan file with designated target and write all data about 'source' file is stored written 'destination' file",
			"each stage other than basic REQUIRE basic stage to be completed",
			"scan stages are:",
			"  basic      - get common data from file header",
			"  interlace  - scan video streams with ffmpeg's 'idet' filter",
			"  silence    - scan audio streams for segments with loudness less than X Lufs (TODO)",
		}, "\n"),
		ArgsUsage: "args used this way",
		//Category:    "Scan",
		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {
			cfg, _ = config.Load(c.App.Name)
			return nil
		},
		After: func(*cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			return fmt.Errorf("bad usage\n  command 'scan' must be used in conjunction with subcommands:\n  enter 'mfline scan --help' for description")
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			return err
		},
		Subcommands: []*cli.Command{
			{
				Name:        "basic",
				Aliases:     []string{},
				Usage:       "get common data from file header",
				UsageText:   "mfline scan basic --source FILE [--destination DIR]\n\ndestination priority:\n -flag        (special case: 'local' = runtime dir)\n -config value\n -runtime dir\n\nno arguments expected\n\nmfline scan basic -? - for detailed args usage",
				Description: "DESCR",
				ArgsUsage:   "ARGS",
				Category:    "",
				Action: func(c *cli.Context) error {
					overwrite := c.Bool("overwrite")
					//CHECK SOURCE
					sourceFile := c.String("source")
					if err := checkSource(c.String("source")); err != nil {
						return err
					}
					//CHECK DESTINATION
					dest := c.String("destination")
					if dest == "" {
						dest = cfg.StorageDir
					}
					destInfo, err := os.Stat(dest)
					if os.IsNotExist(err) {
						return fmt.Errorf("destination is not exist: %v", dest)
					}
					if !destInfo.IsDir() {
						return fmt.Errorf("destination must be a directory: %v", dest)
					}
					if err != nil {
						return fmt.Errorf("os.Stat: %v", err.Error())
					}
					//COMENCE BASIC SCAN
					mp := ump.NewProfile()
					err = mp.ConsumeFile(sourceFile)
					if err != nil {
						return err
					}
					if mp.ConfirmScan(ump.ScanBasic, overwrite) != nil {
						return err
					}
					//SAVE TO TARGET FILE
					bt, err := mp.MarshalJSON()
					if err != nil {
						return err
					}
					fname := filepath.Base(sourceFile)
					trgInfo, _ := os.Stat(dest + fname + ".json")

					if !overwrite && trgInfo != nil {
						return fmt.Errorf("previous scan data exist: overwrite forbidden")
					}
					f, err := os.OpenFile(dest+fname+".json", os.O_CREATE|os.O_WRONLY, 0777)
					if err != nil {
						return fmt.Errorf("can't save target file: %v", err.Error())
					}
					defer f.Close()
					f.Truncate(0)
					_, err = f.Write(bt)

					return err
				},
				OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
					fmt.Println("usage error")
					fmt.Println(err.Error())
					fmt.Println("sub =", isSubcommand)
					return nil
				},
				Subcommands: []*cli.Command{},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "source",
						//Category:  "path",
						Usage:     "path to file which will be scanned   (required)",
						Required:  true,
						Aliases:   []string{},
						TakesFile: true,
					},

					&cli.StringFlag{
						Name: "destination",
						//Category:  "path",
						Usage:     "path to file data will be written to           \n\nother:",
						Required:  false,
						TakesFile: true,
					},
					&cli.BoolFlag{
						Name:               "overwrite",
						Usage:              "rewrite destination file",
						DisableDefaultText: true,
						Aliases:            []string{"o"},
					},
					&cli.BoolFlag{
						Name:               "args",
						Usage:              "show detailed args usage",
						Aliases:            []string{"?"},
						DisableDefaultText: true,
					},
				},
				SkipFlagParsing:        false,
				HideHelp:               false,
				HideHelpCommand:        false,
				Hidden:                 false,
				UseShortOptionHandling: false,
				HelpName:               "",
			},
			{
				Name:        "interlace",
				Aliases:     []string{},
				Usage:       "use idet to scan file for interlace data",
				UsageText:   "mfline scan interlace --source FILE [--destination DIR]\n\ndestination priority:\n -flag        (special case: 'local' = runtime dir)\n -config value\n -runtime dir\n\nno arguments expected\n\nmfline scan basic -? - for detailed args usage",
				Description: "DESCR",
				ArgsUsage:   "ARGS",
				Category:    "",
				Action: func(c *cli.Context) error {
					overwrite := c.Bool("overwrite")
					//CHECK SOURCE
					sourceFile := c.String("source")
					if err := checkSource(sourceFile); err != nil {
						return err
					}

					//CHECK DESTINATION
					dest := c.String("destination")
					if dest == "" {
						dest = cfg.StorageDir
					}
					destInfo, err := os.Stat(dest)
					if os.IsNotExist(err) {
						return fmt.Errorf("destination is not exist: %v", dest)
					}
					if !destInfo.IsDir() {
						return fmt.Errorf("destination must be a directory: %v", dest)
					}
					if err != nil {
						return fmt.Errorf("os.Stat: %v", err.Error())
					}
					fs, err := os.ReadDir(dest)
					if err != nil {
						return err
					}
					validMP := ump.NewProfile()
					path := ""
					basicScanConfirmed := false
					for _, f := range fs {
						if f.IsDir() {
							continue
						}
						mp := ump.NewProfile()
						err = mp.ConsumeJSON(dest + f.Name())
						if err != nil {
							fmt.Fprintf(os.Stderr, "read JSON: %v", err.Error())
							return err
						}

						if mp.Format.Filename == sourceFile {
							validMP = mp
							path = validMP.Format.Filename
							for _, scanMark := range mp.ScansCompleted {
								if scanMark == ump.ScanBasic {
									basicScanConfirmed = true
								}
							}
							continue
						}
					}
					if path == "" || !basicScanConfirmed {
						return fmt.Errorf("path = '%v' for source '%v'\n", path, sourceFile)
					}
					for _, scans := range validMP.ScansCompleted {
						if scans == ump.ScanInterlace && !overwrite {
							return fmt.Errorf("previous scan data exist: overwrite forbidden \nuse -o flag to overwrite data")
						}
					}
					//COMENCE INTERLACE SCAN
					frames := 9999
					devnull := ""
					switch runtime.GOOS {
					case "linux":
						devnull = "/dev/null"
					case "windows":
						devnull = "NUL"
					}

					com := fmt.Sprintf("ffmpeg -hide_banner -filter:v idet -frames:v %v -an -f rawvideo -y %v -i %v", frames, devnull, path)
					fmt.Fprintf(os.Stderr, "run: %v\n", com)

					done := false
					var wg sync.WaitGroup
					process, err := command.New(command.CommandLineArguments(com),
						command.AddBuffer("buf"),
					)
					if err != nil {
						fmt.Fprintf(os.Stderr, "subprocess error1: %v", err.Error())
						return err
					}
					buf := process.Buffer("buf")
					wg.Add(1)
					go func() {
						err = process.Run()
						if err != nil {
							fmt.Fprintf(os.Stderr, "subprocess error2: %v", err.Error())
						}
						done = true
						wg.Done()
					}()
					for !done {
						time.Sleep(time.Millisecond * 500)
						//bts, _ := os.ReadFile("aaa.txt")

						ln := strings.Split(buf.String(), "\n")
						last := len(ln) - 1
						if last < 0 {
							last = 0
						}
						if strings.Contains(ln[last], "s/s speed=") {
							fmt.Fprintf(os.Stderr, "%v\r", ln[last])
						}
					}
					fmt.Fprintf(os.Stderr, "\n")
					wg.Wait()
					if err != nil {
						fmt.Fprintf(os.Stderr, "subprocess error3: %v", err.Error())
						return err
					}

					//ANALYZE REPORT
					idetReport := filterIdet(buf.String())
					sum := float64(idetReport["I"] + idetReport["P"])
					progressive := float64(idetReport["P"]) / sum
					progressive = float64(int(progressive*10000)) / 100
					validMP.Streams[0].Progressive_frames_pct = progressive
					if err := validMP.ConfirmScan(ump.ScanInterlace, overwrite); err != nil {
						fmt.Fprintf(os.Stderr, "subprocess error: %v", err.Error())
						return err
					}
					//SAVE TO TARGET FILE
					bt, err := validMP.MarshalJSON()
					if err != nil {
						fmt.Fprintf(os.Stderr, "subprocess error: %v", err.Error())
						return err
					}
					fname := filepath.Base(sourceFile)
					f, err := os.OpenFile(dest+fname+".json", os.O_CREATE|os.O_WRONLY, 0777)
					if err != nil {
						return fmt.Errorf("can't save target file: %v", err.Error())
					}
					defer f.Close()
					f.Truncate(0)
					_, err = f.Write(bt)
					if err != nil {
						return err
					}
					text := fmt.Sprintf("%v", validMP.Streams[0].Progressive_frames_pct) + "%" + " progressive"
					fmt.Fprint(os.Stdout, text)
					return err
				},
				OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
					fmt.Println("usage error")
					fmt.Println(err.Error())
					fmt.Println("sub =", isSubcommand)
					return nil
				},
				Subcommands: []*cli.Command{},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "source",
						//Category:  "path",
						Usage:     "path to file which will be scanned   (required)",
						Required:  true,
						Aliases:   []string{},
						TakesFile: true,
					},

					&cli.StringFlag{
						Name: "destination",
						//Category:  "path",
						Usage:     "path to file data will be written to           \n\nother:",
						Required:  false,
						TakesFile: true,
					},
					&cli.BoolFlag{
						Name:               "overwrite",
						Usage:              "rewrite destination file",
						DisableDefaultText: true,
						Aliases:            []string{"o"},
					},
					&cli.BoolFlag{
						Name:               "args",
						Usage:              "show detailed args usage",
						Aliases:            []string{"?"},
						DisableDefaultText: true,
					},
				},
				SkipFlagParsing:        false,
				HideHelp:               false,
				HideHelpCommand:        false,
				Hidden:                 false,
				UseShortOptionHandling: false,
				HelpName:               "",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "fl1",
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
			},
		},
		SkipFlagParsing:        false,
		HideHelp:               false,
		HideHelpCommand:        false,
		Hidden:                 false,
		UseShortOptionHandling: false,
		HelpName:               "",
		CustomHelpTemplate:     "",
	}
	return cm
}

func checkSource(sourceFile string) error {
	srcInfo, err := os.Stat(sourceFile)
	if os.IsNotExist(err) {
		return fmt.Errorf("source is not exist: %v", sourceFile)
	}
	if srcInfo.IsDir() {
		return fmt.Errorf("source must not be a directory: %v", sourceFile)
	}
	if err != nil {
		return fmt.Errorf("os.Stat: %v", err.Error())
	}
	return nil
}

func filterIdet(report string) map[string]int {
	rMap := make(map[string]int)
	lines := strings.Split(report, "\n")
	for _, ln := range lines {
		if !strings.Contains(ln, "Parsed_idet") {
			continue
		}
		if strings.Contains(ln, "Repeated ") {
			words := strings.Split(ln, " ")
			fn := 0
			for _, w := range words {
				n, err := strconv.Atoi(w)
				if err != nil {
					continue
				}
				fn++
				rMap["T"] += n
				switch fn {
				case 2, 3:
					rMap["I"] += n
				}
			}
		}
		if strings.Contains(ln, "Single ") {
			words := strings.Split(ln, " ")
			fn := 0
			for _, w := range words {
				n, err := strconv.Atoi(w)
				if err != nil {
					continue
				}
				fn++
				rMap["T"] += n
				switch fn {
				case 1, 2:
					rMap["I"] += n
				case 3:
					rMap["P"] += n
				}
			}
		}
		if strings.Contains(ln, "Multi ") {
			words := strings.Split(ln, " ")
			fn := 0
			for _, w := range words {
				n, err := strconv.Atoi(w)
				if err != nil {
					continue
				}
				fn++
				rMap["T"] += n
				switch fn {
				case 1, 2:
					rMap["I"] += n
				case 3:
					rMap["P"] += n
				}
			}
		}
	}
	return rMap
}
