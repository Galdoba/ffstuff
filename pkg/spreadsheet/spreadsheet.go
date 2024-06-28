package spreadsheet

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
)

/*
csv получается из curl запроса. в запросе могут меняться локальные данные от машины к машине, поэтому строку для запроса пока лучше хранить в Конфиге.
TODO: посмотреть строится ли запрос на csv просто из ссылки.
*/
const (
	SpreadsheetDataPath    = "SpreadsheetDataPath"
	SpreadsheetCurlRequest = "SpreadsheetCurlRequest"
	SheetComment           = iota
	SheetPath
	SheetReadyTrailer
	SheetTrailerStatus
	SheetTrailerMaker
	SheetPosterStatus
	SheetPosterMaker
	SheetLineData0
	SheetTaskName
	SheetTaskStatus
	SheetLineData1
	SheetLineData2
	SheetAgent
	SheetPublicationDate
)

var configFields map[string]string

// func init() {
// 	conf, err := config.ReadConfig()
// 	if err != nil {
// 		fmt.Println("init error :=", err.Error())
// 		return
// 	}
// 	configFields = conf.Field
// }

type spsht struct {
	url      string
	curl     string
	csvPath  string
	csvDataR [][]string
}

//New - объект хранящий данные, расположение и метод обновления
func New() (*spsht, error) {
	sp := spsht{}
	sp.csvPath = configFields[SpreadsheetDataPath]
	sp.curl = configFields[SpreadsheetCurlRequest]
	// if err := sp.updateCSV(); err != nil {
	// 	return &sp, fmt.Errorf("sp.updateCSV() = %v", err.Error())
	// }
	if err := sp.fillCSVData(); err != nil {
		return &sp, fmt.Errorf("sp.fillCSVData() = %v", err.Error())
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
	if err := sp.fillCSVData(); err != nil {
		return fmt.Errorf("sp Update(): sp.fillCSVData() = %v", err.Error())
	}
	fmt.Println("Update Status: ok")
	return nil
}

func (sp *spsht) fillCSVData() error {
	file, err := os.Open(sp.csvPath)
	if err != nil {
		return err
	}
	defer file.Close()
	sp.csvDataR = nil
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), `","`)                   //разделяем строку на ячейки
		data[0] = strings.TrimPrefix(data[0], `"`)                     //чистим первую ячейку
		data[len(data)-1] = strings.TrimSuffix(data[len(data)-1], `"`) //и последнюю
		sp.csvDataR = append(sp.csvDataR, data)
	}
	return scanner.Err()
}

func (sp *spsht) Cell(row, col int) string {
	return sp.csvDataR[row][col]
}

func (sp *spsht) Data() [][]string {
	return sp.csvDataR
}
