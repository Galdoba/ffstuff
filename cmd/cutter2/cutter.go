package main

// import (
// 	"errors"
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/Galdoba/ffstuff/clipmaker"
// 	"github.com/Galdoba/ffstuff/ediread"
// 	"github.com/Galdoba/ffstuff/fldr"
// 	gcli "github.com/Galdoba/ffstuff/pkg/cli"
// 	"github.com/Galdoba/ffstuff/pkg/config"
// 	"github.com/Galdoba/ffstuff/pkg/namedata"
// 	"github.com/urfave/cli/v2"
// )

// var configMap map[string]string

// var edlPaths []string
// var inFolder string
// var preMuxFolder string

// func init() {
// 	conf, err := config.ReadProgramConfig("ffstuff")
// 	if err != nil {
// 		println(err.Error())
// 		os.Exit(1)
// 		//fmt.Println(err)
// 	}
// 	configMap = conf.Field
// 	if err != nil {
// 		switch err.Error() {
// 		case "Config file not found":
// 			println(fmt.Sprintf("Expecting config file in:\n", conf.Path))
// 			//fmt.Print("Expecting config file in:\n", conf.Path)
// 			os.Exit(1)
// 		}
// 	}
// 	inFolder = fldr.InPath()
// 	preMuxFolder = fldr.MuxPath()
// }

// func main() {
// 	app := cli.NewApp()
// 	app.Version = "v 0.0.0"
// 	app.Name = "cutter"
// 	app.Commands = []*cli.Command{

// 		{
// 			Name:        "cut",
// 			Aliases:     []string{},
// 			Usage:       "",
// 			UsageText:   "",
// 			Description: "",
// 			ArgsUsage:   "",
// 			Category:    "",
// 			BashComplete: func(*cli.Context) {
// 				fmt.Println("Start bashcomplete action")
// 			},
// 			Before: func(c *cli.Context) error {
// 				fmt.Println("Start before action")
// 				edlFound := 0

// 				for _, filepath := range c.Args().Tail() {
// 					f, err := os.Stat(filepath)
// 					switch {
// 					default:
// 						return errors.New("unknown before action error: " + err.Error())
// 					case err == nil:
// 						edlFound++
// 						fmt.Println(filepath)
// 						fmt.Println(f.Name(), "is valid file")
// 						edlPaths = append(edlPaths, f.Name())
// 					case cannotFindFile(err):
// 						fmt.Println("Error: Can't find file specified:", filepath)
// 						fmt.Println("Solution: skip file")
// 					}
// 				}
// 				fmt.Println(c.Args())
// 				if c.Bool("testflag") {
// 					fmt.Println("testflag is active")
// 				}
// 				if edlFound == 0 {
// 					fmt.Println("no valid edl-files detected")
// 					return errors.New("Before Action end error")
// 				}
// 				return nil
// 			},
// 			After: func(*cli.Context) error {
// 				fmt.Println("Start after action")
// 				return nil
// 			},

// 			Action: func(c *cli.Context) error {
// 				fmt.Println("Start main action")
// 				fmt.Println("List of edls:", edlPaths)
// 				errors := []string{}
// 				sourceDir := ""
// 				targetDir := ""
// 				for i, edlFile := range edlPaths {
// 					fmt.Printf("File %v/%v", i+1, len(edlPaths))
// 					edi, err := ediread.NewEdlData(edlFile)
// 					if err != nil {
// 						errors = append(errors, err.Error())
// 						//fmt.Println(err.Error())
// 						println(err.Error())
// 					}
// 					path := namedata.RetrieveDirectory(edlFile)
// 					targetDir = path
// 					if c.String("targetfolder") != "" {
// 						tf := c.String("targetfolder")
// 						os.Mkdir(tf, 0700)
// 						sourceDir = path
// 						targetDir = tf
// 					}
// 					if targetDir == "\\" {
// 						targetDir = ""
// 					}
// 					if sourceDir == "\\" {
// 						sourceDir = ""
// 					}
// 					cliTasks := []gcli.Task{}
// 					clipMap := clipmaker.NewClipMap()
// 					for _, clipData := range edi.Entry() {
// 						//fmt.Println(clipData)
// 						//////////////////////////////////
// 						f, err := os.OpenFile(targetDir+"cutter_commands.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
// 						if err != nil {
// 							panic(err)
// 						}
// 						defer f.Close()
// 						// if _, err = f.WriteString(clipData + "\n"); err != nil {
// 						// 	fmt.Println(err)
// 						// }
// 						//////////////////////////////////
// 						cl, err := clipmaker.NewClip(clipData)
// 						if err != nil {
// 							//fmt.Println(err)
// 							println(err.Error())
// 						}
// 						clipMap[cl.Index()] = cl
// 						newTask := gcli.NewTask(clipmaker.CutClipD(cl, sourceDir, targetDir))
// 						cliTasks = append(cliTasks, newTask)
// 						if _, err = f.WriteString(newTask.String() + "\n"); err != nil {
// 							//fmt.Println(err)
// 							println(err.Error())
// 						}
// 					}
// 					cliTasks = sortTasks(cliTasks)

// 					for _, task := range cliTasks {
// 						fmt.Print("RUN:", task, "\n")
// 						taskErr := task.Run()
// 						if taskErr != nil {
// 							errors = append(errors, taskErr.Error())
// 							//fmt.Println(taskErr.Error())
// 							println(taskErr.Error())
// 						}
// 					}
// 				}
// 				//MAINACTION:
// 				/*
// 					обязательные входящие данные: EDL-file
// 					пример: cutter cut filename.edl

// 				*/
// 				return nil
// 			},
// 			// OnUsageError: func(context *cli.Context, err error, isSubcommand bool) error {
// 			// 	fmt.Println("Start on error action")
// 			// 	return nil
// 			// },
// 			//Subcommands: []*cli.Command{},
// 			Flags: []cli.Flag{
// 				&cli.StringFlag{
// 					Name:  "targetfolder, tf",
// 					Usage: "If used sets output folder to cut files to (default: same folder as edl-file's)",
// 					Value: "",
// 				},
// 			},
// 			SkipFlagParsing:        false,
// 			HideHelp:               false,
// 			Hidden:                 false,
// 			UseShortOptionHandling: false,
// 			HelpName:               "",
// 			CustomHelpTemplate:     "",
// 		},
// 	}
// 	args := os.Args
// 	if err := app.Run(args); err != nil {
// 		println(err.Error())
// 		os.Exit(1)
// 		//fmt.Println(err.Error())
// 	}
// }

// func cannotFindFile(err error) bool {
// 	if strings.Contains(err.Error(), "The system cannot find the file specified.") {
// 		return true
// 	}
// 	return false
// }

// //ставит резку аудио перед резкой видео.
// func sortTasks(unsorted []gcli.Task) []gcli.Task {
// 	sorted := []gcli.Task{}
// 	for _, task := range unsorted {
// 		if strings.Contains(task.String(), "_ACLIP_") {
// 			sorted = append(sorted, task)
// 		}
// 	}
// 	for _, task := range unsorted {
// 		if strings.Contains(task.String(), "_VCLIP_") {
// 			sorted = append(sorted, task)
// 		}
// 	}
// 	if len(sorted) != len(unsorted) {
// 		fmt.Println("Can not sort")
// 		return unsorted
// 	}
// 	return sorted
// }
