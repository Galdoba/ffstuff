package tablemanager

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/translit"
)

const (
	NoData = iota
	BadData
	ReadyTrailerExpected
	ReadyTrailerUploadedAhead
	TrailerMaterialUploaded
	TrailerInWork
	TrailerReady
	TrailerUploaded
	PosterInWork
	PosterReady
	PosterUploded
	FilmProblem
	FilmInBuffer
	FilmUploaded
	FilmDownloading
	MuxingInwork
	MuxingReady
	MuxingUploaded
	RowTypeHeader
	RowTypeSeparator
	RowTypeInfo

	SpreadsheetDataPath    = "SpreadsheetDataPath"
	SpreadsheetCurlRequest = "SpreadsheetCurlRequest"
)

type TableData interface {
	Cell(int, int) string
	Data() [][]string
}

//tasks := spreadsheet.Parse(info)
type TaskList struct {
	tasks       []TaskData
	parseErrors []error
}

func TaskListFrom(sp TableData) *TaskList {
	tl := TaskList{}
	for _, data := range sp.Data() {
		r, err := ParseRow(data)
		if err != nil {
			tl.parseErrors = append(tl.parseErrors, err)
			continue
		}
		tl.tasks = append(tl.tasks, r)
	}
	return &tl
}

func (tl *TaskList) ParseErrors() []error {
	return tl.parseErrors
}

func (tl *TaskList) Downloading() []TaskData {
	list := []TaskData{}
	for _, task := range tl.tasks {
		if task.filmStatus != FilmDownloading {
			continue
		}
		if task.path != "" {
			continue
		}
		list = append(list, task)
	}
	return list
}

func (tl *TaskList) ALL() []TaskData {
	list := []TaskData{}
	for _, task := range tl.tasks {

		list = append(list, task)
	}
	return list
}

func (tl *TaskList) ReadyForDemux() []TaskData {
	list := []TaskData{}
	for _, task := range tl.tasks {
		if task.filmStatus != FilmInBuffer {
			continue
		}
		if task.path != "" {
			continue
		}
		list = append(list, task)
	}
	return list
}

func (tl *TaskList) ReadyForEdit() []TaskData {
	list := []TaskData{}
	for _, task := range tl.tasks {
		if task.rowType != RowTypeInfo {
			continue
		}
		if task.filmStatus != FilmInBuffer {
			continue
		}
		if task.path == "" {
			continue
		}

		if strings.Contains(task.comment, "готовит") || strings.Contains(task.comment, "отмен") {
			continue
		}
		if strings.Contains(task.comment, "уже готов") {
			continue
		}
		list = append(list, task)
	}
	return list
}

func (tl *TaskList) ChooseTrailer() []TaskData {
	list := []TaskData{}
	for _, task := range tl.tasks {
		if task.rowType != RowTypeInfo {
			continue
		}
		if task.readyTrailerStatus != ReadyTrailerExpected {
			continue
		}
		if task.muxingStatus == MuxingUploaded {
			continue
		}
		list = append(list, task)
	}
	return list
}

func (tl *TaskList) ChooseFilm() []TaskData {
	list := []TaskData{}
	for _, task := range tl.tasks {
		if task.rowType != RowTypeInfo {
			continue
		}
		nameTransl := translit.Transliterate(task.taskName)
		if strings.Contains(nameTransl, "sezon") {
			continue
		}
		if task.path != "" {
			continue
		}
		if task.filmStatus != FilmInBuffer {
			continue
		}
		list = append(list, task)
	}
	return list
}

func (tl *TaskList) ChooseSeason() []TaskData {
	list := []TaskData{}
	seasonList := []string{}
	for _, task := range tl.tasks {
		if task.rowType != RowTypeInfo {
			continue
		}
		nameTransl := translit.Transliterate(task.taskName)
		if !strings.Contains(nameTransl, "sezon") {
			continue
		}
		if task.path != "" {
			continue
		}
		if task.filmStatus <= FilmProblem {
			continue
		}
		alreadyHave := false

		for _, inList := range seasonList {
			if task.outputName.outBase+"_s"+task.outputName.season == inList {
				alreadyHave = true
				break
			}
		}
		if !alreadyHave {
			list = append(list, task)
			seasonList = append(seasonList, task.outputName.outBase+"_s"+task.outputName.season)
		}

	}
	return list
}

/*
clear  && mkdir -p /mnt/aakkulov/ROOT/IN/_MEGO_DISTRIBUTION/_DONE/Skvoz_ogon  && mkdir -p /mnt/aakkulov/ROOT/EDIT/_mego_distribushn/
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

func (tl *TaskList) ByName(name string) TaskData {
	for _, v := range tl.tasks {
		if v.Name() == name {
			return v
		}
	}
	return TaskData{}
}

func ProposeTargetDirectoryTrailer() string {
	path := `\\nas\ROOT\EDIT\@trailers_temp\`
	return path
}

func ProposeTargetDirectory(tl *TaskList, task TaskData) string {
	path := ``
	if task.rowType != RowTypeInfo {
		return path
	}
	separator := TaskData{}
	folder1 := ""
	for _, r := range tl.tasks {
		if r.taskName == task.taskName {
			break
		}
		if r.rowType == RowTypeSeparator {
			separator = r
			date, err := newDate(separator.taskName)
			folder1 = date.pathFolder()
			if err != nil {
				agent := translit.Transliterate(task.contragent)
				agent = strings.ToLower(agent)
				folder1 = "_" + agent
			}
		}
	}
	add := addSeasonSubFolder(task)
	path = path + folder1
	if add != "" {
		path += add
	}
	return path
}

func ProposeArchiveDirectory(task TaskData) string {
	agentFolderName := translit.Transliterate(task.contragent)
	nameFolderName := strings.Title(task.outputName.outBase)
	if task.outputName.season != "" {
		nameFolderName += "_s" + task.outputName.season
	}
	//отделить сериалы от фильмов
	path := "_" + strings.ToUpper(agentFolderName) // + `/_DONE/` + nameFolderName
	return path
}

func ProposeArchiveDirectoryTrailer(task TaskData) string {
	//`\\nas\ROOT\IN\@TRAILERS\_DONE\Bolshoe_nebo_s01_TRL\`
	nameFolderName := strings.Title(task.outputName.outBase)
	if task.outputName.season != "" {
		nameFolderName += "_s" + task.outputName.season
	}
	//отделить сериалы от фильмов
	path := `\\nas\ROOT\IN\@TRAILERS\_DONE\` + nameFolderName + `\`
	return path
}

func (t *TaskData) OutputBaseName() string {
	name := t.outputName.outBase
	if t.outputName.season != "" {
		name += "_s" + t.outputName.season
	}
	if t.outputName.episode != "" {
		name += "_" + t.outputName.episode
	}
	name = strings.Title(name)
	return name
}

func (t *TaskData) StringAsTrailer() string {
	name := t.outputName.outBase
	return name
}

func (t *TaskData) StringAsFilm() string {
	name := t.outputName.outBase + " [" + t.contragent + "]"
	return name
}

func (t *TaskData) StringAsSeason() string {
	name := t.outputName.outBase + "_s" + t.outputName.season + " [" + t.contragent + "]"
	return name
}

func addSeasonSubFolder(t TaskData) string {
	if t.outputName.season != "" {
		return `/` + t.outputName.outBase + "_s" + t.outputName.season
	}
	return ""
}

func (d *date) pathFolder() string {
	yr := d.year % 100
	mn := strconv.Itoa(d.month)
	if d.month < 10 {
		mn = "0" + mn
	}
	dy := strconv.Itoa(d.day)
	if d.day < 10 {
		dy = "0" + dy
	}
	return fmt.Sprintf("%v_%v_%v", yr, mn, dy)
}

type TaskData struct {
	comment            string
	path               string //1
	readyTrailerStatus int
	trailerStatus      int
	trailerMaker       string //4
	posterStatus       int
	posterMaker        string //6
	dataRow            bool
	taskName           string //8
	filmStatus         int    //
	muxingStatus       int    //10
	urgent             bool
	veryUrgent         bool   //12
	contragent         string //может быть int
	publicationDate    date
	rowType            int
	outputName         outName
}

type outName struct {
	outBase  string
	season   string
	episode  string
	comments []string
}

func ParseRow(data []string) (TaskData, error) {
	r := TaskData{}
	if strings.Join(data, `","`) == `Комментарий","Путь","ГТ","Т","Трейлер","П","Постеры","М","Наименование","С","З","О","!","Контрагент","Дата публикации` {
		r.rowType = RowTypeHeader
		return r, nil
	}
	if len(data) != 15 {
		return r, fmt.Errorf("row format incorect (expect len(data) = 15 have) %v", len(data))
	}
	sep := strings.Join(data[2:7], "")
	if sep == "" {
		r.rowType = RowTypeSeparator
	}

	for i, val := range data {
		switch i {
		case 0:
			r.comment = val
		case 1:
			r.path = val
		case 2:
			switch val {
			default:
				r.readyTrailerStatus = BadData
			case "":
				r.readyTrailerStatus = NoData
			case "r", "к":
				r.readyTrailerStatus = ReadyTrailerExpected
			case "y", "н":
				r.readyTrailerStatus = ReadyTrailerUploadedAhead
			case "g", "п":
				r.readyTrailerStatus = BadData
			}
		case 3:
			switch val {
			default:
				r.trailerStatus = BadData
			case "":
				r.trailerStatus = NoData
			case "r", "к":
				r.trailerStatus = TrailerInWork
			case "y", "н":
				r.trailerStatus = TrailerReady
			case "g", "п":
				r.trailerStatus = TrailerUploaded
			}
		case 4:
			r.trailerMaker = val
		case 5:
			switch val {
			default:
				r.posterStatus = BadData
			case "":
				r.posterStatus = NoData
			case "r", "к":
				r.posterStatus = PosterInWork
			case "y", "н":
				r.posterStatus = PosterReady
			case "g", "п":
				r.posterStatus = PosterUploded
			}
		case 6:
			r.posterMaker = val
		case 7:
			if val == "O" {
				r.dataRow = true
				r.rowType = RowTypeInfo
			}
		case 8:
			r.taskName = val
		case 9:
			switch val {
			default:
				r.filmStatus = BadData
			case "":
				r.filmStatus = NoData
			case "r", "к":
				r.filmStatus = FilmProblem
			case "y", "н":
				r.filmStatus = FilmInBuffer
			case "g", "п":
				r.filmStatus = FilmUploaded
			case "b", "и", "v", "м":
				r.filmStatus = FilmDownloading
			}
		case 10:
			switch val {
			default:
				r.muxingStatus = BadData
			case "":
				r.muxingStatus = NoData
			case "r", "к":
				r.muxingStatus = MuxingInwork
			case "y", "н":
				r.muxingStatus = MuxingReady
			case "g", "п", "G", "П":
				r.muxingStatus = MuxingUploaded
			}
		case 11:
			switch val {
			default:
				r.urgent = false
			case "v", "м":
				r.urgent = true
			}
		case 12:
			switch val {
			default:
				r.veryUrgent = false
			case "v", "м":
				r.veryUrgent = true
			}
		case 13:
			r.contragent = val
		case 14:
			d, err := newDate(val)
			if err != nil {
				return r, err
			}
			r.publicationDate = d
		}
	}
	on, err := defineOutputName(r.taskName)
	if err != nil {
		return r, err
	}
	r.outputName = on
	return r, nil
}

func defineOutputName(name string) (outName, error) {
	on := outName{}
	cmnts := strings.Split(name, "(")
	trName := translit.Transliterate(cmnts[0])
	on.outBase = trName
	for i, com := range cmnts {
		if i == 0 {
			continue
		}
		on.comments = append(on.comments, com)
	}
	if len(on.comments) > 1 {
		last := strings.Split(name, ")")[0]
		on.comments[len(on.comments)-1] = last
	}
	re, e := regexp.Compile(`.*_\d{2}_sezon_\d{2,3}_seri.*`)
	if e != nil {
		return on, e
	}
	if re.Match([]byte(trName)) {
		parts := strings.Split(trName, "_sezon_")
		rejoined := strings.Join(parts[0:len(parts)-1], "_sezon_")
		name := strings.Split(rejoined, "_")
		on.outBase = strings.Join(name[0:len(name)-1], "_")
		on.season = name[len(name)-1]
		on.episode = strings.Split(parts[len(parts)-1], "_")[0]
	}
	return on, nil
}

func (r *TaskData) Name() string {
	return r.taskName
}

func (r *TaskData) PublicationDate() string {
	return r.publicationDate.String()
}

func (r *TaskData) Agent() string {
	return r.contragent
}

func (r *TaskData) String() string {
	str := fmt.Sprintf("%v [%v]", r.taskName, r.contragent)
	if r.comment != "" {
		str += " (" + r.comment + ")"
	}
	return str
}

func (r *TaskData) Match(str string) bool {
	if r.String() == str {
		return true
	}
	return false
}

type date struct {
	day   int
	month int
	year  int
}

func newDate(s string) (date, error) {
	d := date{}
	if s == "" || s == "Дата публикации" {
		return d, nil
	}
	data := strings.Split(s, ".")
	if len(data) != 3 {
		return d, fmt.Errorf("date format incorect (%v)", data[0])
	}
	day, err := strconv.Atoi(data[0])
	if err != nil {
		return d, fmt.Errorf("date format incorect (day) - %v", data[0])
	}
	d.day = day
	month, err := strconv.Atoi(data[1])
	if err != nil {
		return d, fmt.Errorf("date format incorect (month) - %v", data[1])
	}
	d.month = month
	year, err := strconv.Atoi(data[2])
	if err != nil {
		return d, fmt.Errorf("date format incorect (year) - %v", data[2])
	}
	d.year = year
	return d, nil
}

func (d *date) String() string {
	str := ""
	if d.day+d.month+d.year == 0 {
		return "          "
	}
	if d.day <= 9 && d.day >= 0 {
		str += "0"
	}
	str += fmt.Sprintf("%v.", d.day)
	if d.month <= 9 && d.month >= 0 {
		str += "0"
	}
	str += fmt.Sprintf("%v.", d.month)
	str += fmt.Sprintf("%v", d.year)
	if d.year == 0 {
		str += "000"
	}
	return str
}
