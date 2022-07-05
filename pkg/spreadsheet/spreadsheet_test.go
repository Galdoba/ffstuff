package spreadsheet

import (
	"fmt"
	"strings"
	"testing"
)

func TestSpreadsheet(t *testing.T) {
	sp, err := New()
	if sp == nil {
		t.Errorf("func New() returned no object")
	}
	if err != nil {
		t.Errorf("func New() returned error: %v", err.Error())
	}
	if sp.curl == "" {
		t.Errorf("no curl request set")
	}
	if sp.csvPath == "" {
		t.Errorf("no path set")
	}
	if len(sp.csvData) == 0 {
		t.Errorf("no csv data found")
	}
	for _, v := range sp.csvData {
		fmt.Println(v)
	}
	//sp.Update()
}

func TestParsing(t *testing.T) {
	sp, _ := New()
	for _, line := range sp.csvData {
		data := strings.Split(line, `","`)
		row, err := parseRow(line)
		if err != nil {
			t.Errorf("%v\nparseRow(line) returned error: %v", line, err.Error())
			fmt.Println(row)
			continue
		}
		if row.rowType == "" {
			t.Errorf("%v\nrowType not asigned: %v", line, data[2])
		}
		if row.readyTrailerStatus == badData {
			t.Errorf("%v\nunknown data for readyTrailerStatus: %v", line, data[2])
		}
		if row.trailerStatus == badData {
			t.Errorf("%v\nunknown data for trailerStatus: %v", line, data[3])
		}
		if row.posterStatus == badData {
			t.Errorf("%v\nunknown data for posterStatus: %v", line, data[5])
		}
		if row.taskName == "" && row.rowType == "INFO" {
			t.Errorf("%v\ntask is unnamed: %v", line, data[8])
		}
		if row.filmStatus == badData {
			t.Errorf("%v\nunknown data for filmStatus: %v", line, data[9])
		}

		//fmt.Println(row)
	}

}
