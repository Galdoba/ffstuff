package actiontrailer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/ffstuff/pkg/spreadsheet"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet/tablemanager"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
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
	// conf, err := config.Read()
	// if err != nil {
	// 	return "", err
	// }
	args := c.Args()
	if len(args) < 1 {
		return "", fmt.Errorf("no arguments provided")
	}
	if len(args) > 1 {
		return "", fmt.Errorf("to many arguments provided")
	}
	//Declarations:
	//cur_dir := currentDir()
	//file := args[0]
	//work_dir := conf["InProgressDir"]
	//videoMapping := ""
	//atempo := ""
	//destination := ""
	//name := ""
	//vtag := ""
	//atag := ""
	//trl_done_dir := conf["TrlDoneDir"]

	//comLine := ""
	//mediaReport, err := probe.MediaFileReport(file, probe.MediaTypeTrailerHD)
	//atempo = mediaReport.FPS()

	trailerTask, err := selectTask()
	fmt.Println(trailerTask, err)

	// path := args[0]
	// mediaInfo, err := probe.MediaFileReport(path, probe.MediaTypeTrailerHD)
	// atempo := mediaInfo.FPS()

	return "", nil
}

func currentDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)

}

func selectTask() (tablemanager.Row, error) {
	sheet, err := spreadsheet.New()
	if err != nil {
		return tablemanager.Row{}, err
	}
	taskList := tablemanager.TaskListFrom(sheet)
	activeTasks := []tablemanager.Row{}
	valid := survey.ComposeValidators()
	options := []string{}
	selected := ""
	for _, task := range taskList.ReadyTrailers() {
		options = append(options, task.String())
		activeTasks = append(activeTasks, task)
	}
	selected, err = askSelection(valid, "Выберите Трейлер:", options)
	for _, task := range activeTasks {
		if task.String() == selected {
			return task, nil
		}
	}
	return tablemanager.Row{}, fmt.Errorf("Задача не найдена")
}

func askSelection(validator survey.Validator, message string, options []string) (string, error) {
	choose := ""
	promptSelect := &survey.Select{
		Message: message,
		Options: append(options, "Отмена"),
	}
	if err := survey.AskOne(promptSelect, &choose, validator); err != nil {
		return choose, err
	}
	if choose == "Отмена" {
		return "", fmt.Errorf("была выбрана `Отмена`")
	}

	return choose, nil
}

func askInput(val survey.Validator, message string) (string, error) {
	result := ""
	promptInput := &survey.Input{
		Message: message,
	}
	if err := survey.AskOne(promptInput, &result, val); err != nil {
		return result, err
	}
	return result, nil
}
