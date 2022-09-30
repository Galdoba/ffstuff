package actiondemux

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/Galdoba/ffstuff/pkg/mdm/probe"

	"github.com/Galdoba/ffstuff/pkg/mdm/format"

	"github.com/Galdoba/ffstuff/pkg/spreadsheet/tablemanager"

	"github.com/Galdoba/devtools/cli/command"

	"github.com/Galdoba/ffstuff/app/demuxer/handle"

	"github.com/urfave/cli"
)

/*
ПРИМЕРЫ ПРИМЕНЕНИЯ
demuxer -tofile file.txt -update demux -i film.mp4
	-tofile file.txt - терминал будет писаться в указанный файл

Для демукса требуются:
1. Исходник(и)
2. Информация по заданию (данные из таблицы)
3. Задание (ввод в ручную что это фильм/трейлер/сериал)

ПЛАН:
1. Собираем данные:
	1.1 Подтверждаем исходник(и)
	1.2 Запрашиваем задание
	1.3 Читаем таблицу.
	1.4 ДЕБАГ: Выводим имена и пути предпологаемых результатов.


*/

const (
	taskTypeFILM    = "Фильм"
	taskTypeTRAILER = "Трейлер (НЕ РАБОТАЕТ)"
	taskTypeSERIES  = "Сериал  (НЕ РАБОТАЕТ)"
)

var inputBuffer []string
var inputPaths []string

func Run(c *cli.Context) error {
	fmt.Println("RUN Precheck")
	if err := Precheck(c); err != nil {
		return err
	}
	fmt.Println("Precheck complete")
	args := c.Args()
	for _, arg := range args {
		taskType := handle.SelectionSingle("Что в исходнике?", []string{taskTypeFILM, taskTypeTRAILER}...)
		// task := tablemanager.TaskData{}
		fmt.Println("в исходнике: ", taskType)
		task, err := DefineTask(taskType)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		///////////////////////

		targetDir := `\\nas\ROOT\EDIT\` + tablemanager.ProposeTargetDirectory(handle.TaskListFull(), task)
		outBaseName := task.OutputBaseName()
		fmt.Printf(`DEMUX DIR  : %v%v`+"\n", targetDir, outBaseName)
		if err := confirmDir(targetDir); err != nil {
			return err
		}
		archive := tablemanager.ProposeArchiveDirectory(task)
		fmt.Printf("ARCHIVE DIR: %v\n", archive)
		if err := confirmDir(archive); err != nil {
			return err
		}
		tFormat := &format.TargetFormat{}
		switch {
		default:
			tFormat, _ = format.SetAs(format.FilmHD)
		case strings.Contains(task.Name(), " SD"):
			tFormat, _ = format.SetAs(format.FilmSD)
		case strings.Contains(task.Name(), " 4K"):
			tFormat, _ = format.SetAs(format.Film4K)
		}
		fmt.Println(" ")
		fmt.Println("Check Input:")
		for _, issue := range Issues(tFormat, task) {
			fmt.Println("WARNING: " + issue)
		}
		rep, _ := probe.NewReport(arg)
		audioStreams := rep.Audio()
		needSelect := false
		if len(audioStreams) > 2 {
			needSelect = true
		}
		for _, as := range audioStreams {
			switch as.ChanLayout() {
			case "stereo", "5.1", "5.1(side)":
			default:
				needSelect = true
			}
		}
		selectedAudio := audioStreams
		if needSelect {
			selectedAudio = probe.SelectAudio(rep.Audio())
		}
		langData := []string{}
		langAdded := make(map[string]int)

		//checkMono(ad []probe.AudioData)

		for a, str := range selectedAudio {
			fmt.Println(" ")
			lang := str.FCmapKey() + "__AUDIO"
			lang += handle.SelectionSingle(fmt.Sprintf("Язык стрима [0:%v]: %v", a, str.String()), []string{"RUS", "ENG", "QQQ"}...)
			switch str.ChanLayout() {
			default:
				lang += "ХЗ"
			case "stereo":
				lang += "20"
			case "5.1", "5.1(side)":
				lang += "51"
			}
			if langAdded[lang] > 0 {
				lang += fmt.Sprintf("_%v", langAdded[lang])
			}
			langAdded[lang]++
			langData = append(langData, lang)
		}

		fmt.Println("||||||||||||||||")
		proxy := ""
		proxy = handle.SelectionSingle("Делать Прокси?", "ДА", "НЕТ")

		ffmpegCMD, errff := formAmediaFFmpegComLine(arg, langData, task, proxy)
		if errff != nil {
			fmt.Println(errff)
		}
		//fmt.Println(errff)
		//fmt.Println("fflite", ffmpegCMD)
		//распределяем звуковые тэги по потокам ...ок
		//составляем фильтр комплекс ..............ок
		//формируем команду .......................ок
		//дополняем транспортные команды
		/*
			clear && \
			mv  ~/IN/FILE ~/IN/_IN_PROGRESS/ && \
			fflite -r 25 -i ~/IN/_IN_PROGRESS/FILE \
			-filter_complex "[0:a:0]aresample=48000,atempo=25/24[audio]" \
			-map [audio]    @alac0 NAME_AUDIORUS20.m4a \
			-map 0:v:0      @crf16 NAME_HD.mp4 \
			&& touch NAME.ready \
			&& at now + 10 hours <<< "mv ~/IN/_IN_PROGRESS/FILE OUTPATH"
		*/

		fullCommand := formFullCommandFilm(arg, ffmpegCMD, targetDir, outBaseName)
		processOS := handle.SelectionSingle("На какой операционной системе будет просчет?", "Linux", "Windows")
		switch processOS {
		default:
			return fmt.Errorf("выбрана неизвестная ОС")
		case "Windows":
		case "Linux":
			fullCommand = handle.ConvertToLinux(fullCommand)
		}

		fmt.Println(" ")
		fmt.Printf("Скопируйте строку в терминал машины, что будет заниматься просчетом:\n")
		fmt.Println(color.HiGreenString(fmt.Sprintf("\n%v\n", fullCommand)))
		fmt.Println(" ")
		//echo %DATE% %TIME% >\\nas\buffer\IN\%v.ready
		//&& mv /home/pemaltynov/IN/_IN_PROGRESS/Потеря_надежды_Hope_Lost_2015_2.mkv /home/pemaltynov/IN/_DONE/
		//уточняем машину на которой будет считаться
		if err := handle.ArchiveDelay(arg, archive); err != nil {
			return err
		}
	}
	return nil
}

func Precheck(c *cli.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided")
	}
	fmt.Println(" ")
	for _, arg := range args {
		fmt.Println(arg)
	}
	fmt.Println(" ")
	selected := handle.SelectionSingle("Перечень исходных файлов корректен?", []string{"ДА", "НЕТ"}...)
	if selected != "ДА" {
		return fmt.Errorf("User abort")
	}

	for _, arg := range args {
		for _, s := range []string{" ", "(", ")"} {
			if strings.Contains(arg, s) {
				fmt.Printf("Внимание: Имя файла содержит недопустимый символ: '%v'\n", s)
				return fmt.Errorf("недопустимый символ '%v'", s)
			}
		}

		com, err := command.New(
			command.CommandLineArguments("fflite", "-i "+arg),
			command.Set(command.TERMINAL_OFF),
			command.Set(command.BUFFER_ON),
		)
		if err != nil {
			return err
		}
		fmt.Println(" ")
		com.Run()
		inputBuffer = strings.Split(com.StdOut()+"\n"+com.StdErr(), "\n")
	}
	for _, line := range inputBuffer {
		fmt.Println(color.HiYellowString(line))
	}

	return nil
}

func DefineTask(taskType string) (tablemanager.TaskData, error) {
	task := tablemanager.TaskData{}
	taskList := []tablemanager.TaskData{}
	switch taskType {
	default:
		return task, fmt.Errorf("DefineTask(taskType): taskType=%v (unknown)", taskType)
	case taskTypeFILM, taskTypeTRAILER:
		taskList = handle.SelectFromTable(taskType)
	}
	s := []string{}
	for _, t := range taskList {
		s = append(s, t.String())
	}
	taskStr := handle.SelectionSingle("Данные из таблицы: ", s...)
	for _, t := range taskList {
		if t.Match(taskStr) {
			task = t
		}
	}
	return task, nil
}

func Issues(tFormat *format.TargetFormat, task tablemanager.TaskData) []string {
	fmt.Println("DEBUG: Checking Issues")
	videoSizeValid := false
	issues := []string{}
	videoFound := 0
	audioFound := 0
	soundMap := make(map[string]int)
	for _, data := range inputBuffer {
		switch {
		case strings.Contains(data, `Video: `):
			videoFound++
			if strings.Contains(data, ` 1920x1080`) && (!strings.Contains(task.Name(), " SD") && !strings.Contains(task.Name(), " 4K")) {
				videoSizeValid = true
			}
			if strings.Contains(data, ` 720x576`) && strings.Contains(task.Name(), " SD") {
				videoSizeValid = true
			}
			if strings.Contains(data, ` 3840x2160`) && strings.Contains(task.Name(), " 4K") {
				videoSizeValid = true
			}
			if !videoSizeValid {
				issues = append(issues, "scaling needed")
			}
		case strings.Contains(data, `Audio: `):
			audioFound++
			if strings.Contains(data, ` stereo`) {
				soundMap["stereo"]++
			}
			if strings.Contains(data, ` 5.1`) {
				soundMap["5.1"]++
			}
			if strings.Contains(data, ` 5.1(side)`) {
				soundMap["5.1"]++
			}
			if strings.Contains(data, ` mono`) {
				soundMap["mono"]++
			}
			if audioFound != mapSum(soundMap) {
				soundMap["warning"]++
			}
		}
	}
	if soundMap["warning"] > 0 {
		issues = append(issues, fmt.Sprintf("%v audio streams require attention", soundMap["warning"]))
	}
	return issues
}

func mapSum(sm map[string]int) int {
	s := 0
	for _, v := range sm {
		s += v
	}
	return s
}

func dupeSlice(sl []string) []string {
	sR := []string{}
	for _, s := range sl {
		sR = append(sR, s)
		sR = append(sR, s)
	}
	return sR
}

func formAmediaFFmpegComLine(inputPath string, audioTags []string, task tablemanager.TaskData, proxy string) (string, error) {
	fcUsed := false
	rep, err := probe.NewReport(inputPath)
	if err != nil {
		return "", err
	}
	fps := rep.FPS()
	str := "-r 25"
	str += fmt.Sprintf(` -i \\nas\buffer\IN\_IN_PROGRESS\%v`, inputPath)
	str += fmt.Sprintf(` -filter_complex`)
	if proxy == "ДА" {
		fcUsed = true
		str += fmt.Sprintf(` "[0:v:0]scale=iw/2:ih, setsar=(1/1)*2[vidp]; `)
	}
	activeMap := []string{}
	stream2LangMAP := make(map[string]string)
	for i, audStr := range audioTags {
		if i == 0 && fcUsed == false {
			str += ` "`
		}
		fcUsed = true
		fmt.Println("audStr=", audStr)
		key := strings.Split(audStr, "__")
		actMap := fmt.Sprintf("aud%v", i)
		str += fmt.Sprintf(`[0%v]aresample=48000,atempo=25/(%v)[%v]; `, key[0], fps, actMap)
		activeMap = append(activeMap, actMap)
		fmt.Println("ADDED=", actMap)
		stream2LangMAP[actMap] = key[1]
		if proxy == "ДА" {
			str += fmt.Sprintf(`[0%v]aresample=48000,atempo=25/(%v)[%v]; `, key[0], fps, actMap+"_pr")
			activeMap = append(activeMap, actMap+"_pr")
			fmt.Println("ADDED=", actMap+"_pr")
			stream2LangMAP[actMap+"_pr"] = key[1] + "_proxy"
		}
		//str = strings.TrimSuffix(str, "; ")
		// if i == len(audioTags)-1 {

		// 	str += `"`
		// }
	}

	if !fcUsed {
		str = strings.Replace(str, " -filter_complex", "", -1)
	} else {
		str = strings.TrimSuffix(str, "; ")
		str += `"`
	}
	crf := "16"
	switch {
	default:
	case strings.Contains(task.Name(), " 4K"):
		crf = "18"
	case strings.Contains(task.Name(), " SD"):
		crf = "13"
	}
	str += fmt.Sprintf(` -map 0:v:0 -c:v libx264 -preset medium -crf %v -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 \\nas\ROOT\EDIT\%v%v_HD.mp4`, crf, tablemanager.ProposeTargetDirectory(handle.TaskListFull(), task), task.OutputBaseName())
	if proxy == "ДА" {
		str += fmt.Sprintf(` -map [vidp] -c:v libx264 -x264opts interlaced=1 -preset medium -pix_fmt yuv420p  -b:v 2000k -maxrate 2000k \\nas\ROOT\EDIT\%v%v_HD_proxy.mp4`, tablemanager.ProposeTargetDirectory(handle.TaskListFull(), task), task.OutputBaseName())
	}
	for s, key := range activeMap {
		fmt.Println("DEB", s, key, stream2LangMAP[key])
		switch {
		default:
			//fmt.Println("ERR", s, key)
			audTag := stream2LangMAP[key]
			str += fmt.Sprintf(` -map [%v] -c:a -vn -acodec alac -compression_level 0 -map_metadata -1 -map_chapters -1 \\nas\ROOT\EDIT\%v%v_%v.m4a`, key, tablemanager.ProposeTargetDirectory(handle.TaskListFull(), task), task.OutputBaseName(), audTag)
		case strings.Contains(key, "_pr"):
			audTag := stream2LangMAP[key]
			str += fmt.Sprintf(` -map [%v] -c:a ac3 -b:a 128k \\nas\ROOT\EDIT\%v%v_%v.ac3`, key, tablemanager.ProposeTargetDirectory(handle.TaskListFull(), task), task.OutputBaseName(), audTag)
		}

	}

	return str, nil
}

func formFullCommandFilm(inputFile, ffmpegCMD, targetDir, outBaseName string) string {
	fullCommand := ""
	fullCommand += fmt.Sprintf(`cls`)
	fullCommand += fmt.Sprintf(` && move \\nas\buffer\IN\%v \\nas\buffer\IN\_IN_PROGRESS\`, inputFile)
	fullCommand += fmt.Sprintf(` && fflite %v`, ffmpegCMD)
	fullCommand += fmt.Sprintf(` && echo 1>` + targetDir + outBaseName + ".ready")
	fullCommand += fmt.Sprintf(` && move \\nas\buffer\IN\_IN_PROGRESS\%v \\nas\buffer\IN\_DONE\`, inputFile)
	fullCommand += fmt.Sprintf(` && echo 1> \\nas\buffer\IN\z_TASK_COMPLETE_` + outBaseName + ".txt")
	//fullCommand += fmt.Sprintf(` && \\nas\buffer\IN\TASK_COMPLETE_` + outBaseName + ".txt")
	return fullCommand
}

/*
\\nas\buffer\IN\bats\checkForErrors.bat
/home/pemaltynov/IN/bats/checkForErrors.bat
*/

func confirmDir(d string) error {
	dir, err := os.Stat(d)
	if err != nil {
		switch {
		default:
			return fmt.Errorf("confirming directory: %v", err.Error())
		case strings.Contains(err.Error(), "system cannot find the file specified"):
			switch handle.SelectionSingle(fmt.Sprintf("WARNING:\n'%v' не существует!\nСоздать?", d), "ДА", "НЕТ") {
			case "ДА":
				os.MkdirAll(d, 0777)
				return nil
			case "НЕТ":
				return fmt.Errorf("нет папки назначения")
			}
		}
	}
	if !dir.Mode().IsDir() {
		return fmt.Errorf("%v is not a directory", d)
	}
	return nil
}

/*
func RunUNTESTED(c *cli.Context) error {
	args, err := CheckArguments(c)
	if err != nil {
		return err
	}
	streams, err := probe.StreamsData(args...)
	if err != nil {
		return err
	}
	vid, aud := probe.SeparateByTypes(streams)
	if len(vid)+len(aud) == 0 {
		return fmt.Errorf("file(s) contains no video and no audio streams")
	}
	for _, stream := range streams {
		fmt.Println(stream.PrintStreamData())
	}
	fmt.Println("TODO: Выбрать Схему желаемого результата")
	taskType := handle.SelectionSingle("Что делаем?", taskTypeFILM, taskTypeTRAILER, taskTypeSERIES)
	fmt.Println(taskType)
	fmt.Println("TODO: Выбрать позицию в таблице")
	fmt.Println("TODO: Выбрать интересующие видео дорожки")
	fmt.Println("TODO: Выбрать интересующие аудио дорожки")

	return nil
	fmt.Println("RUN Precheck")
	if err := Precheck(c); err != nil {
		return err
	}
	fmt.Println("Precheck complete")
	args = c.Args()
	for _, arg := range args {
		taskType := handle.SelectionSingle("Что в исходнике?", []string{taskTypeFILM, taskTypeTRAILER}...)
		// task := tablemanager.TaskData{}
		fmt.Println("в исходнике: ", taskType)
		task, err := DefineTask(taskType)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		///////////////////////

		targetDir := `\\nas\ROOT\EDIT\` + tablemanager.ProposeTargetDirectory(handle.TaskListFull(), task)
		outBaseName := task.OutputBaseName()
		fmt.Printf(`DEMUX DIR  : %v%v`+"\n", targetDir, outBaseName)
		if err := confirmDir(targetDir); err != nil {
			return err
		}
		archive := tablemanager.ProposeArchiveDirectory(task)
		fmt.Printf("ARCHIVE DIR: %v\n", archive)
		if err := confirmDir(archive); err != nil {
			return err
		}
		tFormat := &format.TargetFormat{}
		switch {
		default:
			tFormat, _ = format.SetAs(format.FilmHD)
		case strings.Contains(task.Name(), " SD"):
			tFormat, _ = format.SetAs(format.FilmSD)
		case strings.Contains(task.Name(), " 4K"):
			tFormat, _ = format.SetAs(format.Film4K)
		}
		fmt.Println(" ")
		fmt.Println("Check Input:")
		for _, issue := range Issues(tFormat, task) {
			fmt.Println("WARNING: " + issue)
		}
		rep, _ := probe.NewReport(arg)
		audioStreams := rep.Audio()
		needSelect := false
		if len(audioStreams) > 2 {
			needSelect = true
		}
		for _, as := range audioStreams {
			switch as.ChanLayout() {
			case "stereo", "5.1", "5.1(side)":
			default:
				needSelect = true
			}
		}
		selectedAudio := audioStreams
		if needSelect {
			selectedAudio = probe.SelectAudio(rep.Audio())
		}
		langData := []string{}
		langAdded := make(map[string]int)

		//checkMono(ad []probe.AudioData)

		for a, str := range selectedAudio {
			fmt.Println(" ")
			lang := str.FCmapKey() + "__AUDIO"
			lang += handle.SelectionSingle(fmt.Sprintf("Язык стрима [0:%v]: %v", a, str.String()), []string{"RUS", "ENG", "QQQ"}...)
			switch str.ChanLayout() {
			default:
				lang += "ХЗ"
			case "stereo":
				lang += "20"
			case "5.1", "5.1(side)":
				lang += "51"
			}
			if langAdded[lang] > 0 {
				lang += fmt.Sprintf("_%v", langAdded[lang])
			}
			langAdded[lang]++
			langData = append(langData, lang)
		}

		fmt.Println("||||||||||||||||")
		ffmpegCMD, errff := formAmediaFFmpegComLine(arg, langData, task)
		if errff != nil {
			fmt.Println(errff)
		}

		fullCommand := formFullCommandFilm(arg, ffmpegCMD, targetDir, outBaseName)
		processOS := handle.SelectionSingle("На какой операционной системе будет просчет?", "Linux", "Windows")
		switch processOS {
		default:
			return fmt.Errorf("выбрана неизвестная ОС")
		case "Windows":
		case "Linux":
			fullCommand = handle.ConvertToLinux(fullCommand)
		}

		fmt.Println(" ")
		fmt.Printf("Скопируйте строку в терминал машины, что будет заниматься просчетом:\n")
		fmt.Println(color.HiGreenString(fmt.Sprintf("\n%v\n", fullCommand)))
		fmt.Println(" ")
		//echo %DATE% %TIME% >\\nas\buffer\IN\%v.ready
		//&& mv /home/pemaltynov/IN/_IN_PROGRESS/Потеря_надежды_Hope_Lost_2015_2.mkv /home/pemaltynov/IN/_DONE/
		//уточняем машину на которой будет считаться
		if err := handle.ArchiveDelay(arg, archive); err != nil {
			return err
		}
	}
	return nil
}
*/
func CheckArguments(c *cli.Context) ([]string, error) {
	args := c.Args()
	if len(args) == 0 {
		return nil, fmt.Errorf("no arguments provided")
	}
	fmt.Println(" ")
	for _, arg := range args {
		fmt.Println(arg)
	}
	fmt.Println(" ")
	selected := handle.SelectionSingle("Перечень исходных файлов корректен?", []string{"ДА", "НЕТ"}...)
	if selected != "ДА" {
		return nil, fmt.Errorf("перечень исходных файлов не корректен")
	}
	for _, arg := range args {
		for _, s := range []string{" ", "(", ")"} {
			if strings.Contains(arg, s) {
				fmt.Printf("Внимание: Имя файла содержит недопустимый символ: '%v'\n", s)
				return nil, fmt.Errorf("недопустимый символ '%v'", s)
			}
		}
		com, err := command.New(
			command.CommandLineArguments("fflite", "-i "+arg),
			command.Set(command.TERMINAL_OFF),
			command.Set(command.BUFFER_ON),
		)
		if err != nil {
			return nil, err
		}
		fmt.Println(" ")
		com.Run()
		inputBuffer = append(inputBuffer, strings.Split(com.StdOut(), "\n")...)
	}
	for _, line := range inputBuffer {
		fmt.Println(color.HiYellowString(line))
	}
	return args, nil
}
