package actiontrailer

import (
	"fmt"
	"os"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/ffstuff/app/demuxer/handle"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/mdm/probe"
	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet/tablemanager"
	"github.com/urfave/cli"
)

/*
# Трейлер HD:
MOVE \\nas\buffer\IN\FILE \\nas\buffer\IN\_IN_PROGRESS\ && \
fflite -r 25 -i \\nas\buffer\IN\_IN_PROGRESS\FILE \
-filter_complex "[0:a:0]aresample=48000,atempo=25/24[audio]" \
-map [audio]    @alac0 \\nas\buffer\IN\_IN_PROGRESS\NAME_TRL_AUDIORUS20.m4a \
-map 0:v:0      @crf10 \\nas\buffer\IN\_IN_PROGRESS\NAME_TRL_HD.mp4 \
&& MD \\nas\ROOT\IN\@TRAILERS\_DONE\NAME_TRL && MOVE \\nas\buffer\IN\_IN_PROGRESS\FILE \\nas\ROOT\IN\@TRAILERS\_DONE\NAME_TRL && MOVE \\nas\buffer\IN\_IN_PROGRESS\NAME_TRL_* \\nas\ROOT\EDIT\@trailers_temp\ && exit

# Трейлер HD DOWNSCALE:
MOVE \\nas\buffer\IN\FILE \\nas\buffer\IN\_IN_PROGRESS\ && \
fflite -r 25 -i \\nas\buffer\IN\_IN_PROGRESS\FILE \
-filter_complex "[0:a:0]aresample=48000,atempo=25/24[audio]; [0:v:0]scale=1920:-2,setsar=1/1,unsharp=3:3:0.3:3:3:0,pad=1920:1080:-1:-1[video_hd]" \
-map [audio]    @alac0 \\nas\buffer\IN\_IN_PROGRESS\NAME_TRL_AUDIORUS20.m4a \
-map [video_hd] @crf10 \\nas\buffer\IN\_IN_PROGRESS\NAME_TRL_HD.mp4 \
&& MD \\nas\ROOT\IN\@TRAILERS\_DONE\NAME_TRL && MOVE \\nas\buffer\IN\_IN_PROGRESS\FILE \\nas\ROOT\IN\@TRAILERS\_DONE\NAME_TRL && MOVE \\nas\buffer\IN\_IN_PROGRESS\NAME_TRL_* \\nas\ROOT\EDIT\@trailers_temp\ && exit



fflite @crf10
> ffmpeg -an -vcodec libx264 -preset medium -crf 10 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1
$ fflite @alac0
> ffmpeg -vn -acodec alac -compression_level 0 -map_metadata -1 -map_chapters -1
ПЛАН:
анализировать трейлер
выбрать интересующие стримы с которыми будем работать
формируем команду
	выбираем видео канал
		узнаем atempo
		узнаем размер
			узнаем надо ли даунскейлить
	выбираем звуковой канал
		выбираем язык
		выбираем канальность
прерываем если ошибка
перемещаем трейлер в папку inprogress
выполняем команду
перемещаем результаты в trailers_temp
перемещаем трейлер в trailers_done
завершаем программу
составить маршруты движения файлов
составить переменные для -filter_complex video
составить переменные для -filter_complex audio
запустить
*/

const (
	start = iota
)

type director struct {
	err          error
	currentStage int
}

/*
MOVE [CURRENT_DIR][FILE] [WORKFOLDER]
ffmpeg -y -r 25 -i [WORKFOLDER][FILE]
-filter complex [0:v:0][VIDEOMAPPING][video];[0:a:0]aresample=48000,atempo=[ATEMPO][audio]
-map [video] -an -vcodec libx264 -preset medium -crf 10 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 [DESTINATION][NAME][VTAG]
-map [audio] -vn -acodec alac -compression_level 0 -map_metadata -1 -map_chapters -1 [DESTINATION][NAME][ATAG]
MD [TRAILER_DONE_DIR][NAME]
MOVE [FILE] [TRAILER_DONE_DIR][NAME]
[CURRENT_DIR]      - определяем по ходу
[FILE]             - аргумент
[WORKFOLDER]       - константа
[ATEMPO]           - аналитика файла
[VIDEOMAPPING]     - аналитика файла
[DESTINATION]      - аналитика таблицы
[NAME]             - аналитика таблицы
[VTAG]             - определяем по ходу
[ATAG]             - определяем по ходу
[TRAILER_DONE_DIR] - константа
*/

func Run(c *cli.Context) (string, error) {
	conf, err := config.Read()
	handle.Error(err)
	args := c.Args()
	if len(args) < 1 {
		return "", fmt.Errorf("no arguments provided")
	}
	if len(args) > 1 {
		return "", fmt.Errorf("to many arguments provided")
	}
	//Declarations:
	cur_dir := currentDir()
	fmt.Println("cur_dir", cur_dir)
	file := args[0]
	com, _ := command.New(
		command.CommandLineArguments("fflite", "-i ", file),
		command.Set(command.TERMINAL_ON),
	)
	com.Run()
	work_dir := conf["InProgressDir"]
	fmt.Println("work_dir", work_dir)
	//videoMapping := ""
	//atempo := ""
	//destination := ""
	//name := ""
	//vtag := ""
	//atag := ""
	trl_done_dir := conf["TrlDoneDir"]
	fmt.Println("trl_done_dir", trl_done_dir)
	//comLine := ""
	// mediaReport, err := probe.MediaFileReport(file, probe.MediaTypeTrailerHD)
	// handle.Error(err)
	//atempo = mediaReport.FPS()

	task, err := selectTask()
	handle.Error(err)

	mediaInfo, err := probe.NewReport(file)

	//MOVE [CURRENT_DIR][FILE] [WORKFOLDER]
	res := fmt.Sprintf("%v\\%v %v ", cur_dir, file, work_dir)
	//ffmpeg -y -r 25 -i [WORKFOLDER][FILE]
	res += fmt.Sprintf("&& fflite -y -r 25 -i %v\\%v ", work_dir, file)
	//-filter complex [0:v:0][VIDEOMAPPING][video];[0:a:0]aresample=48000,atempo=[ATEMPO][audio]
	videomapping, atempo := VideoMapping(file) //TODO: нужна адекватная общая функция дающая строку для -filter_complex от всего файла.
	res += fmt.Sprintf("-filter_complex [0:v:0]%v[video];[0:a:0]aresample=48000,atempo=25/(%v)[audio]", videomapping, atempo)
	//-map [video] -an -vcodec libx264 -preset medium -crf 10 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 [DESTINATION][NAME][VTAG]
	name := namedata.TransliterateForEdit(task.Name())
	vtag := "_HD_TRL"
	res += fmt.Sprintf(" -map [video] -an -vcodec libx264 -preset medium -crf 10 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 %v%v%v.mp4 ", trl_done_dir, name, vtag)
	//-map [audio] -vn -acodec alac -compression_level 0 -map_metadata -1 -map_chapters -1 [DESTINATION][NAME][ATAG]
	aStream := mediaInfo.Audio()
	atag := analyzeAudio(aStream)
	res += fmt.Sprintf(" -map [audio] -vn -acodec alac -compression_level 0 -map_metadata -1 -map_chapters -1 %v%v%v ", trl_done_dir, name, atag)
	return res, nil
}

func currentDir() string {
	dir, err := os.Getwd()
	handle.Error(err)
	return dir
}

func selectTask() (tablemanager.TaskData, error) {
	sheet, err := spreadsheet.New()
	handle.Error(err)
	taskList := tablemanager.TaskListFrom(sheet)
	activeTasks := []tablemanager.TaskData{}
	options := []string{}
	for _, task := range taskList.ReadyTrailers() {
		options = append(options, task.String())
		activeTasks = append(activeTasks, task)
	}
	selected := handle.SelectionSingle("Выберите Трейлер:", options...)
	for _, task := range activeTasks {
		if task.String() == selected {
			return task, nil
		}
	}
	return tablemanager.TaskData{}, fmt.Errorf("Задача не найдена")
}

func VideoMapping(path string) (string, string) {
	mr, err := probe.NewReport(path)
	handle.Error(err)

	for _, issue := range mr.Issues() {
		switch issue {
		default:
			fmt.Println("Issue:", issue)
		}
	}
	return "setsar=1/1", mr.FPS()
}

func analyzeAudio(streamData []probe.AudioData) string {
	if len(streamData) == 1 {
		return "_TRL_AUDIORUS" + atag(streamData[0].ChanLayout())
	}
	if len(streamData) < 1 {
		return "no audio"
	}
	combine := ""

	for _, st := range probe.SelectAudio(streamData) {
		combine += atag(st.ChanLayout()) + "|"
	}

	return combine

}

func atag(layout string) string {
	switch layout {
	default:
		return "unknown"
	case "stereo":
		return "20"
	case "5.1":
		return "51"
	}
}
