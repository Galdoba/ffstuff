package spreadsheet

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/ffstuff/pkg/config"
)

/*
csv получается из curl запроса. в запросе могут меняться локальные данные от машины к машине, поэтому строку для запроса пока лучше хранить в Конфиге.
TODO: посмотреть строится ли запрос на csv просто из ссылки.
*/

const (
	noData = iota
	badData
	readyTrailerExpected
	readyTrailerUploadedAhead
	trailerMaterialUploaded
	trailerInWork
	trailerReady
	trailerUploaded
	posterInWork
	posterReady
	posterUploded
	filmProblem
	filmInBuffer
	filmUploaded
	filmDownloading
	muxingInwork
	muxingReady
	muxingUploaded
	SpreadsheetDataPath    = "SpreadsheetDataPath"
	SpreadsheetCurlRequest = "SpreadsheetCurlRequest"
)

var configFields map[string]string

func init() {
	conf, err := config.ReadConfig()
	if err != nil {
		fmt.Println("init error :=", err.Error())
		return
	}
	configFields = conf.Field
}

type spsht struct {
	url     string
	curl    string
	csvPath string
	csvData []string
}

func New() (*spsht, error) {
	sp := spsht{}
	sp.csvPath = configFields[SpreadsheetDataPath]
	sp.curl = configFields[SpreadsheetCurlRequest]
	// if err := sp.updateCSV(); err != nil {
	// 	return &sp, fmt.Errorf("sp.updateCSV() = %v", err.Error())
	// }
	if err := sp.readCSVData(); err != nil {
		return &sp, fmt.Errorf("sp.readCSVData() = %v", err.Error())
	}
	return &sp, nil
}

func (sp *spsht) Update() error {
	comm, err := command.New(
		command.CommandLineArguments("curl "+sp.curl+sp.csvPath),
		command.Set(command.BUFFER_OFF),
		command.Set(command.TERMINAL_ON),
	)
	if err != nil {
		return err
	}
	fmt.Println("Updating Spreadsheet:")
	comm.Run()
	if err := sp.readCSVData(); err != nil {
		return fmt.Errorf("sp Update(): sp.readCSVData() = %v", err.Error())
	}
	fmt.Println("Update Status: ok")
	return nil
}

func (sp *spsht) readCSVData() error {
	file, err := os.Open(sp.csvPath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sp.csvData = append(sp.csvData, scanner.Text())
	}
	return scanner.Err()
}

type row struct {
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
	rowType            string
}

func parseRow(str string) (row, error) {
	r := row{}
	if str == `"Комментарий","Путь","ГТ","Т","Трейлер","П","Постеры","М","Наименование","С","З","О","!","Контрагент","Дата публикации"` {
		r.rowType = "HEADER"
		return r, nil
	}
	data := strings.Split(str, `","`)
	if len(data) != 15 {
		return r, fmt.Errorf("row format incorect")
	}
	sep := strings.Join(data[2:7], "")
	if sep == "" {
		r.rowType = "SEPARATOR"
	}

	data[0] = strings.TrimPrefix(data[0], `"`)
	data[14] = strings.TrimSuffix(data[14], `"`)
	for i, val := range data {
		switch i {
		case 0:
			r.comment = val
		case 1:
			r.path = val
		case 2:
			switch val {
			default:
				r.readyTrailerStatus = badData
			case "":
				r.readyTrailerStatus = noData
			case "r", "к":
				r.readyTrailerStatus = readyTrailerExpected
			case "y", "н":
				r.readyTrailerStatus = readyTrailerUploadedAhead
			case "g", "п":
				r.readyTrailerStatus = badData
			}
		case 3:
			switch val {
			default:
				r.trailerStatus = badData
			case "":
				r.trailerStatus = noData
			case "r", "к":
				r.trailerStatus = trailerInWork
			case "y", "н":
				r.trailerStatus = trailerReady
			case "g", "п":
				r.trailerStatus = trailerUploaded
			}
		case 4:
			r.posterMaker = val
		case 5:
			switch val {
			default:
				r.posterStatus = badData
			case "":
				r.posterStatus = noData
			case "r", "к":
				r.posterStatus = posterInWork
			case "y", "н":
				r.posterStatus = posterReady
			case "g", "п":
				r.posterStatus = posterUploded
			}
		case 6:
			r.posterMaker = val
		case 7:
			if val == "O" {
				r.dataRow = true
				r.rowType = "INFO"
			}
		case 8:
			r.taskName = val
		case 9:
			switch val {
			default:
				r.filmStatus = badData
			case "":
				r.filmStatus = noData
			case "r", "к":
				r.filmStatus = filmProblem
			case "y", "н":
				r.filmStatus = filmInBuffer
			case "g", "п":
				r.filmStatus = filmUploaded
			case "b", "и":
				r.filmStatus = filmDownloading

			}
		case 10:
			switch val {
			default:
				r.muxingStatus = badData
			case "":
				r.muxingStatus = noData
			case "r", "к":
				r.muxingStatus = muxingInwork
			case "y", "н":
				r.muxingStatus = muxingReady
			case "g", "п":
				r.muxingStatus = muxingUploaded
			}
		case 11:

		case 14:
			d, err := newDate(val)
			if err != nil {
				return r, err
			}
			r.publicationDate = d
		}
	}
	return r, nil
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
	if d.day+d.month+d.year == 0 {
		return ""
	}
	return fmt.Sprintf("%v.%v.%v", d.day, d.month, d.year)
}
