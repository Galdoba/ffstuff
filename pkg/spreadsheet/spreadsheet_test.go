package spreadsheet

import (
	"fmt"
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
	fmt.Println(len(sp.csvData), "lines of text in CSV")
}
