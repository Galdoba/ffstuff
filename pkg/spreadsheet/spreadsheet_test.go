package spreadsheet

import (
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
	if len(sp.csvDataR) == 0 {
		t.Errorf("no csv data found")
	}
	//for _, v := range sp.csvData {
	//fmt.Println(v)
	//}
	sp.Update()
}
