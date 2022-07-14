package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/devtools/cli/command"
	actioncombine "github.com/Galdoba/ffstuff/app/demuxer/actions/combine"
	actiondemux "github.com/Galdoba/ffstuff/app/demuxer/actions/demux"
	actiontrailer "github.com/Galdoba/ffstuff/app/demuxer/actions/trailer"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
)

func main() {
	console, err := command.New(
		command.CommandLineArguments("ffmpeg"),
		command.Set(command.TERMINAL_ON),
		command.Set(command.BUFFER_ON),
	)
	if err != nil {
		fmt.Printf("console error: %v", err.Error())
	}
	resultLine := ""
	app := cli.NewApp()
	app.Version = "v 0.1.1"
	app.Name = "demuxer"
	app.Usage = "Анализирует файлы чтобы составить команду ffmpeg для демукса"
	app.Description = "combine: Ready\n	 demux: TODO\n  fullhelp: PLACEHOLDER\n   trailer: InWork"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "update",
			Usage: "если активен, то до начала выполнения любой команды - обновится csv с рабочей таблицей",
		},
		cli.StringFlag{
			Name:  "tofile",
			Usage: "указывает адрес файла в который будет добавляться вывод терминала",
			Value: "",
		},
		cli.BoolFlag{
			Name:  "process, p",
			Usage: "Если активен, программа запустит ffmpeg с полученной командной строкой",
		},
	}
	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {
		if c.GlobalString("tofile") != "" {
			console.AddInstruction(command.WriteToFile(c.GlobalString("tofile")))
		}
		if c.GlobalBool("update") {
			sheet, err := spreadsheet.New()
			if err != nil {
				return err
			}
			if err := sheet.Update(); err != nil {
				return err
			}
		}
		return nil
	}
	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		if c.GlobalBool("process") {
			console.AddInstruction(command.CommandLineArguments("ffmpeg", resultLine))
			console.AddInstruction(command.Set(command.TERMINAL_ON))
			if err := console.Run(); err != nil {
				return err
			}
		}

		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:        "combine",
			ShortName:   "",
			Usage:       "собирает звук 5.1 из 6 файлов wav",
			UsageText:   "Требует 6 .wav файлов в аргументов, для того чтобы сшить из них 5.1.\n   Файлы должны содержать маркеры каналов: .L./.R./.C./.LFE./.Ls./.Rs. в своём имени",
			Description: "TODO: подробное описание команды",
			ArgsUsage:   "TODO: подробное описание как пользовать аргументы",
			Action: func(c *cli.Context) error {
				comLine, err := actioncombine.Run(c) //Собираем команду
				if err != nil {
					color.Red("Error:") //
					return err
				}
				fmt.Println("\nProcessing command:")
				color.HiCyan(fmt.Sprintf("ffmpeg %v   \n  \n", comLine))
				resultLine = comLine
				return nil
			},
		},
		{
			Name:        "demux",
			ShortName:   "",
			Usage:       "составляет ffmpeg строку для разложения файла на дорожки в стандартной кодировке",
			UsageText:   "исходный файл должен быть аргументом",
			Description: "TODO: подробное описание команды",
			ArgsUsage:   "TODO: подробное описание как пользовать аргументы",

			Action: func(c *cli.Context) error {
				return actiondemux.Run(c)
			},
		},
		{
			Name:      "fullhelp",
			ShortName: "",
			Usage:     "Предоставляет полное описание всего что умеет программа",
			Action: func(c *cli.Context) error {
				fmt.Println("Run fullhelp")
				fmt.Println("End fullhelp")
				return nil
			},
		},
		{
			Name:        "trailer",
			ShortName:   "",
			Usage:       "составляет ffmpeg строку для разложения файла на дорожки в стандартной кодировке для Трейлера",
			UsageText:   "demuxer -update -process trailer [FILE.mp4] (-destination FOLDER) (-archive FOLDER)",
			Description: "TODO: подробное описание команды",
			ArgsUsage:   "TODO: подробное описание как пользовать аргументы",

			Action: func(c *cli.Context) error {
				line, err := actiontrailer.Run(c)
				fmt.Println(line, err)
				return err
			},
		},
	}

	args := os.Args
	if err := app.Run(args); err != nil {
		fmt.Printf("application returned error: %v", err.Error())
	}
	exit := ""
	val := survey.ComposeValidators()
	promptInput := &survey.Input{
		Message: "Enter для завершения",
	}
	survey.AskOne(promptInput, &exit, val)

}

/*
fmpeg -i Сквозь_огонь_Through_the_fire.mkv
  INPUT 0: Сквозь_огонь_Through_the_fire.mkv
  Duration: 01:52:36.12, start: 0.000000, bitrate: 18720 kb/s
    0:0 (eng) Video: h264 (Main), yuv420p(progressive), 1920x1080 [SAR 1:1 DAR 16:9], 25 fps, 25 tbr, 1k tbn, 50 tbc (default)
    0:1 (rus) Audio: mp3, 48000 Hz, stereo, fltp, 320 kb/s (default)
    0:2 (rus) Audio: ac3, 48000 Hz, 5.1(side), fltp, 448 kb/s
    0:3 (fre) Audio: mp3, 48000 Hz, stereo, fltp, 320 kb/s
    0:4 (fre) Audio: ac3, 48000 Hz, 5.1(side), fltp, 448 kb/s
*/
/*
clear
&& mkdir -p /mnt/aakkulov/ROOT/IN/_MEGO_DISTRIBUTION/_DONE/Skvoz_ogon
&& mkdir -p /mnt/aakkulov/ROOT/EDIT/_mego_distribushn/
&& mv /home/aakkulov/IN/Сквозь_огонь_Through_the_fire.mkv /home/aakkulov/IN/_IN_PROGRESS/
&& fflite -r 25 -i /home/aakkulov/IN/_IN_PROGRESS/Сквозь_огонь_Through_the_fire.mkv
	-filter_complex "[0:a:1]aresample=48000,atempo=25/(25)[arus]"
	-map [arus] -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 /mnt/aakkulov/ROOT/EDIT/_mego_distribushn/Skvoz_ogon_AUDIORUS51.m4a
	-map 0:v:0 -c:v libx264 -preset medium -crf 16 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 /mnt/aakkulov/ROOT/EDIT/_mego_distribushn/Skvoz_ogon_HD.mp4
&& touch /mnt/aakkulov/ROOT/EDIT/_mego_distribushn/Skvoz_ogon.ready
&& mv /home/aakkulov/IN/_IN_PROGRESS/Сквозь_огонь_Through_the_fire.mkv /home/aakkulov/IN/_DONE/
&& at now + 10 hours <<< "mv /home/aakkulov/IN/_DONE/Сквозь_огонь_Through_the_fire.mkv /mnt/aakkulov/ROOT/IN/_MEGO_DISTRIBUTION/_DONE/Skvoz_ogon"
&& clear
&& touch /home/aakkulov/IN/TASK_COMPLETE_Сквозь_огонь_Through_the_fire.mkv.txt
*/
/*
fflite -r 25 -i /home/aakkulov/IN/_IN_PROGRESS/Сквозь_огонь_Through_the_fire.mkv
	-filter_complex "[0:a:1]aresample=48000,atempo=25/(25)[arus]"
	-map [arus] -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 /mnt/aakkulov/ROOT/EDIT/_mego_distribushn/Skvoz_ogon_AUDIORUS51.m4a
	-map 0:v:0 -c:v libx264 -preset medium -crf 16 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 /mnt/aakkulov/ROOT/EDIT/_mego_distribushn/Skvoz_ogon_HD.mp4



*/
