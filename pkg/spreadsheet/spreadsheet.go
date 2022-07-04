package spreadsheet

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/ffstuff/pkg/config"
)

/*
csv получается из curl запроса. в запросе могут меняться локальные данные от машины к машине, поэтому строку для запроса пока лучше хранить в Конфиге.
TODO: посмотреть строится ли запрос на csv просто из ссылки.
*/

const (
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
	if err := sp.updateCSV(); err != nil {
		return &sp, fmt.Errorf("sp.updateCSV() = %v", err.Error())
	}
	if err := sp.readCSVData(); err != nil {
		return &sp, fmt.Errorf("sp.readCSVData() = %v", err.Error())
	}
	return &sp, nil
}

func (sp *spsht) updateCSV() error {
	comm, err := command.New(
		command.CommandLineArguments("curl "+sp.curl+sp.csvPath),
		command.Set(command.BUFFER_OFF),
		command.Set(command.TERMINAL_ON),
	)
	if err != nil {
		return err
	}
	comm.Run()
	return nil
}

func (sp *spsht) readCSVData() error {

	file, err := os.Open(sp.csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sp.csvData = append(sp.csvData, scanner.Text())
	}
	return scanner.Err()
}
