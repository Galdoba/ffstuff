package tablemanager

import (
	"fmt"
	"testing"

	"github.com/Galdoba/ffstuff/pkg/spreadsheet"
)

func TestParsing(t *testing.T) {
	return
	sp, _ := spreadsheet.New()
	for _, data := range sp.Data() {
		//data := strings.Split(line, `","`)
		row, err := ParseRow(data)
		if err != nil {
			t.Errorf("%v\nparseRow(line) returned error: %v", data, err.Error())
			fmt.Println(row)
			continue
		}
		if row.rowType == NoData {
			t.Errorf("%v\nrowType not asigned: %v", data, data[2])
		}
		if row.readyTrailerStatus == BadData {
			t.Errorf("%v\nunknown data for readyTrailerStatus: %v", data, data[2])
		}
		if row.trailerStatus == BadData {
			t.Errorf("%v\nunknown data for trailerStatus: %v", data, data[3])
		}
		if row.posterStatus == BadData {
			t.Errorf("%v\nunknown data for posterStatus: %v", data, data[5])
		}
		if row.taskName == "" && row.rowType == RowTypeInfo {
			t.Errorf("%v\ntask is unnamed: %v", data, data[8])
		}
		if row.filmStatus == BadData {
			t.Errorf("%v\nunknown data for filmStatus: %v", data, data[9])
		}
		if row.muxingStatus == BadData {
			t.Errorf("%v\nunknown data for muxingStatus: %v", data, data[10])
		}
		//fmt.Println(row.String())
	}
}

func TestListing(t *testing.T) {
	return
	sheet, _ := spreadsheet.New()
	sheet.Update()
	TaskList := TaskListFrom(sheet)

	fmt.Println("Downloading")
	for _, task := range TaskList.Downloading() {
		fmt.Println(task.String())
	}

	fmt.Println("ReadyForDemux")
	for _, task := range TaskList.ReadyForDemux() {
		fmt.Println(task.String())
	}

	fmt.Println("ReadyForEdit")
	for _, task := range TaskList.ReadyForEdit() {
		fmt.Println(task.String())
	}

	fmt.Println("ReadyTrailer")
	for _, task := range TaskList.ChooseTrailer() {
		fmt.Println(task.String())
	}
}

func TestTargetDirectoryPath(t *testing.T) {
	return
	sheet, _ := spreadsheet.New()
	//sheet.Update()
	taskList := TaskListFrom(sheet)
	for _, task := range taskList.tasks {
		if task.rowType != RowTypeInfo {
			continue
		}
		if task.path != "" {
			continue
		}
		//propose := ProposeTargetDirectorySubFolder(taskList, task)
		//fmt.Printf("%v (%v)\n", task.taskName, task.contragent)
		//fmt.Printf("propose: %v\n", propose)
	}
}
