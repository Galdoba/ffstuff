package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/Galdoba/ffstuff/cmd/grabber/ui"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/ffstuff/pkg/scanner"
	"github.com/Galdoba/ffstuff/pkg/sortnames"
	"github.com/Galdoba/utils"
	"github.com/martinlindhe/notify"
	"github.com/urfave/cli"
)

/*
TZ:
>> grab only [path]				-- забрать только указанные пути
>> grab filename.ready			-- забрать все связанное с ready файлом
>> grab help (-h)				-- вывести на экран помогалку 							--help
>> grab new (-n)				-- забрать все новое (предварительное сканирование)		--new
>> grab -v						-- забрать только если одобряет инчекер					--valid
>> grab -p						-- забрать только звук и прокси							--proxy
>> grab -fc						-- забрать только если одобряет fflite @check0			--fflitecheck0

пред проверки:
-папка куда копировать
-отсуствие файла с таким же именем и размером
-наличие свободного места для копии

пост проверки:
-копия равна по имени и размеру с источником

*/

var configMap map[string]string

//var logger glog.Logger
var username string
var programName string

func init() {
	programName = "grabber"
	fmt.Println("Initialisation...")

	fmt.Println("Reading config file...")
	cDir, cFile := config.ConfigPathManual(programName)
	gc, err := ReadConfig(cDir + "\\" + cFile)
	if err != nil {
		fmt.Println("Error:", err.Error())
		offerToCreateDefaultConfig(cDir)
		err = CreateDefaultConfig()
		if err != nil {
			fmt.Println("Error:", err.Error())
		} else {
			fmt.Println("Default config file created")
			fmt.Println("Restart the program")
			os.Exit(0)
		}
	}

	//conf, err := config.ReadProgramConfig(programName)
	//configMap = conf.Field
	// if err != nil {
	// 	switch err.Error() {
	// 	case "Config file not found":
	// 		cDir, cFile := config.ConfigPathManual(programName)
	// 		fmt.Printf("Expecting config file in: %v\n", cDir+"\\"+cFile)
	// 		answer, err := askSelection("Create default config file?", []string{"YES", "NO"})
	// 		panicIfErr(err)
	// 		switch answer {
	// 		case "YES":
	// 			_, err := config.ConstructManual(programName)
	// 			panicIfErr(err)
	// 			conf, err = config.ReadProgramConfig(programName)
	// 			//генерируем поля для автоконфига
	// 			config.AddCommentManual(programName, "Logging:")
	// 			config.SetFieldManual(programName, "External log", "TODO")
	// 			config.SetFieldManual(programName, "local log", "TODO")
	// 			config.AddCommentManual(programName, "Actions:")
	// 			config.SetFieldManual(programName, "MOVE_CURSOR_UP", "UP")
	// 			config.SetFieldManual(programName, "MOVE_CURSOR_DOWN", "DOWN")
	// 			config.SetFieldManual(programName, "TOGGLE_SELECTION_STATE", "SPACE")
	// 			fmt.Println("Restart the program")
	// 			os.Exit(1)
	// 		case "NO":
	// 			fmt.Println("Can not run program without config")
	// 			os.Exit(0)
	// 		}
	// 	}
	// }
	configMap = make(map[string]string)
	configMap["External_Log_path"] = gc.External_Log_path
	for _, action := range gc.Actions {
		for index, key := range action.Triggers {
			indexedKey := fmt.Sprintf("%v_%v", action.ActionName, index)
			configMap[indexedKey] = key
		}
	}
	for k, v := range configMap {
		fmt.Println(k, v)
	}
	fmt.Println("Config file reading complete...")
	currentUser, userErr := user.Current()
	if userErr != nil {
		fmt.Printf("Initialisation failed: %v", userErr.Error())
	}
	fmt.Print("Username: ")
	username = currentUser.Name
	fmt.Print(username + "\n")
	fmt.Println("Initialisation complete.")
}

func offerToCreateDefaultConfig(cDir string) {
	answ, err := askSelection(fmt.Sprintf("Create/overwrite config at %v\\ ?", cDir), []string{"YES", "NO"})
	if err != nil {
		panic(err.Error())
	}
	if answ == "NO" {
		fmt.Println("Program can not run without config")
		os.Exit(0)
	}
}

func main() {
	//searchRoot := configMap[constant.SearchRoot]
	//searchMarker := configMap[constant.SearchMarker]
	//dest := configMap[constant.InPath] + "IN_" + utils.DateStamp() + "\\"
	//logPath := configMap[constant.MuxPath] + "MUX_" + utils.DateStamp() + "\\logfile.txt"
	//logger = glog.New(logPath, glog.LogLevelINFO)
	//logPath := configMap["EXTERNAL_LOG"]
	//logger = glog.New(logPath, glog.LogLevelINFO)
	destination := ""
	app := cli.NewApp()
	app.Version = "v 0.0.3"
	app.Name = "grabber"
	app.Usage = "dowloads files and sort it to working directories"
	app.Flags = []cli.Flag{

		&cli.StringFlag{
			Name:        "destination, dest",
			Usage:       "folder where grabber will drop files to",
			Required:    true,
			Destination: new(string),
		},
	}

	app.Commands = []cli.Command{
		////////////////////////////////////
		{
			Name:        "take",
			ShortName:   "",
			Aliases:     []string{},
			Usage:       "grabs files to destination folder",
			UsageText:   "grabber --dest [FOLDER] take [FILE_1] [FILE_2] ... [FILE_N]",
			Description: "TODO:Descr",
			ArgsUsage:   "TODO:ArgsUsage",
			Category:    "Operation",

			Action: func(c *cli.Context) error {
				paths := c.Args()

				if len(paths) == 0 {
					fmt.Println("No arguments provided")
					return nil
				}
				dest := ""

				switch c.GlobalString("destination") {
				default:
					dest = c.GlobalString("destination")
					configMap["dest"] = dest
					fmt.Println("destination set as: " + c.GlobalString("destination"))
					//TODO: отучить от необходимости ставить слэшь для аргумента
				case "":
					dest = destination
				}
				list := []string{}
				for i, path := range paths {
					fmt.Printf("%v	argument: %v\n", i, path)
					//continue
					if isReadyfile(path) {
						assosiated, err := scanner.ListAssosiated(path)
						if err != nil {
							return fmt.Errorf("scanner.ListAssosiated(%v): %v", path, err.Error())
						}
						list = append(list, assosiated...)
					}
					list = append(list, path)
				}
				fmt.Printf("sorting...\n")
				fmt.Println(list)
				list = sortnames.GrabberOrder(list)
				fmt.Printf("starting main loop...\n")
				if err := ui.StartMainloop(configMap, list); err != nil {
					return err
				}
				deleteReadyFiles(dest, list)
				fmt.Printf("simulating report making...\n")
				return nil
			},
		},
	}
	args := os.Args
	if len(args) < 2 {
		args = append(args, "help") //Принудительно зовем помощь если нет других аргументов
	}
	err := app.Run(args)
	fmt.Printf("simulating errors processing...\n")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		notify.Notify("Grabber", "Done", "No errors", `c:\Users\pemaltynov\.config\grabber\grabbe2r.png`)
	}

	fmt.Printf("simulating graceful exit...\n")
}

func checkMediaFile(path string, keys ...string) error {
	f, err := os.Stat(path)
	if err != nil {
		return err
	}
	for _, key := range keys {
		switch key {
		default:
		case ".ready", ".mp4", ".m4a":
			if !strings.Contains(f.Name(), key) {
				return fmt.Errorf("%v is not a '%v' file", path, key)
			}
		}
	}
	return nil
}

func deleteReadyFiles(dest string, list []string) {
	usr, _ := user.Current()
	for _, path := range list {
		if !strings.HasSuffix(path, ".ready") {
			continue
		}
		short := namedata.RetrieveShortName(path)
		os.Remove(dest + short)
		os.Remove(path + "_" + usr.Name)
	}

}

func downloadAssociatedWith(l glog.Logger, paths []string, destination string) error {
	marker := ""
	for _, path := range paths {
		if strings.Contains(path, ".ready") {
			marker = strings.TrimSuffix(path, ".ready") + "." + username
			err := os.Rename(path, marker)
			//fmt.Println("Rename", path)
			if err != nil {
				return err
			}
			continue
		}
		// if err := grabber.Download(logger, path, destination); err != nil {
		// 	return err
		// }
	}

	return nil
}

func isReadyfile(path string) bool {
	data := strings.Split(path, ".")
	if data[len(data)-1] == "ready" {
		return true
	}
	return false
}

func ensureValidOrder(sl []string) []string {
	valid := []string{}
	for _, val := range sl {
		if strings.Contains(val, ".ready") {
			valid = append(valid, val)
		}
	}
	for _, val := range sl {
		valid = utils.AppendUniqueStr(valid, val)
	}
	return valid
}
