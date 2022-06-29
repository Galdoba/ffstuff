package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "v 0.0.0"
	app.Name = "director"
	app.Usage = "Scans media streams to tell the story..."
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{}
	args := []string{""}
	if len(args) == 0 {
		fmt.Println("No arguments provided")
		os.Exit(0)
	}
	if err := app.Run(args); err != nil {
		fmt.Println("Here is my error:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Here is my story: ...")

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
