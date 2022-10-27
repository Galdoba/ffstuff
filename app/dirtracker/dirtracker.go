package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/dirtracker/filelist"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
)

func main() {
	app := cli.NewApp()
	app.Version = "v 0.0.0"
	app.Name = "dirtracker"
	app.Usage = "Отслеживает файлы в указанных директориях"
	app.Description = "Должен крутиться постоянно чтобы считывать изменения и вести статистику"
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
		//Убедиться что есть файлы статистики. если нет то создать
		return nil
	}
	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		//Придумаем что-нибудь
		return nil
	}
	app.Commands = []cli.Command{
		//start - запустить программу
		//ShowStats - поразать глубокую аналитику
		{
			Name:        "start",
			ShortName:   "",
			Usage:       "запустить программу",
			UsageText:   "TODO: прописать основную логику",
			Description: "TODO: подробное описание команды",
			ArgsUsage:   "TODO: подробное описание как пользовать аргументы",
			Action: func(c *cli.Context) error {
				//Непосредственно команда
				/*
					Составить Список имеющихся файлов
					Вывести Список
				*/
				fl := filelist.New(`\\nas\buffer\IN\`)
				fl.AddExeption(`\\nas\buffer\IN\.Trash-1000`)
				fl.AddExeption(`\\nas\buffer\IN\@GOBLIN`)
				fl.AddExeption(`\\nas\buffer\IN\@KURAZH_BAMBEY`)
				fl.AddExeption(`\\nas\buffer\IN\@SCRIPTS`)
				fl.AddExeption(`\\nas\buffer\IN\@TEST`)
				fl.AddExeption(`\\nas\buffer\IN\@WINK_SRC`)
				fl.AddExeption(`\\nas\buffer\IN\_AWAIT`)
				//fl.AddExeption(`\\nas\buffer\IN\_DONE`)
				//fl.AddExeption(`\\nas\buffer\IN\_IN_PROGRESS`)
				fl.AddExeption(`\\nas\buffer\IN\_REJECTED`)
				fl.AddExeption(`\\nas\buffer\IN\bats`)
				fl.AddExeption(`\\nas\buffer\IN\CASES`)
				fl.AddExeption(`\\nas\buffer\IN\Doramy`)
				fl.AddExeption(`\\nas\buffer\IN\errors`)
				fl.AddExeption(`\\nas\buffer\IN\ffmpegs`)
				fl.AddExeption(`\\nas\buffer\IN\Torrent`)
				fl.Print()
				return nil
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
