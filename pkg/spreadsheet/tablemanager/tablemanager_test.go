package tablemanager

import (
	"fmt"
	"testing"

	"github.com/Galdoba/ffstuff/pkg/spreadsheet"
)

func TestParsing(t *testing.T) {
	sp, _ := spreadsheet.New()
	for _, data := range sp.Data() {
		//data := strings.Split(line, `","`)
		row, err := parseRow(data)
		if err != nil {
			t.Errorf("%v\nparseRow(line) returned error: %v", data, err.Error())
			fmt.Println(row)
			continue
		}
		if row.rowType == noData {
			t.Errorf("%v\nrowType not asigned: %v", data, data[2])
		}
		if row.readyTrailerStatus == badData {
			t.Errorf("%v\nunknown data for readyTrailerStatus: %v", data, data[2])
		}
		if row.trailerStatus == badData {
			t.Errorf("%v\nunknown data for trailerStatus: %v", data, data[3])
		}
		if row.posterStatus == badData {
			t.Errorf("%v\nunknown data for posterStatus: %v", data, data[5])
		}
		if row.taskName == "" && row.rowType == rowTypeInfo {
			t.Errorf("%v\ntask is unnamed: %v", data, data[8])
		}
		if row.filmStatus == badData {
			t.Errorf("%v\nunknown data for filmStatus: %v", data, data[9])
		}
		if row.muxingStatus == badData {
			t.Errorf("%v\nunknown data for muxingStatus: %v", data, data[10])
		}
		//fmt.Println(row.String())
	}
}

func TestListing(t *testing.T) {
	sheet, _ := spreadsheet.New()
	sheet.Update()
	TaskList := TaskListFrom(sheet)

	fmt.Println("Downloading")
	for _, task := range TaskList.Downloading() {
		fmt.Println(task.String())
		fmt.Println(ProposeTargetDirectory(TaskList, task))
	}

	fmt.Println("ReadyForDemux")
	for _, task := range TaskList.ReadyForDemux() {
		fmt.Println(task.String())
		fmt.Println(ProposeTargetDirectory(TaskList, task))
	}

	fmt.Println("ReadyForEdit")
	for _, task := range TaskList.ReadyForEdit() {
		fmt.Println(task.String())
		fmt.Println(ProposeTargetDirectory(TaskList, task))
	}
}
