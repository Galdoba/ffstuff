package fftask

import (
	"fmt"
	"os"
	"testing"

	"github.com/olekukonko/tablewriter"
)

func TestTask(t *testing.T) {
	testOperations := validOperations()
	for operation, validity := range testOperations {
		fmt.Printf("START testing operation '%v':\n", operation)
		tsk, err := New(operation)
		if err != nil {
			t.Errorf("unexpected error: %v", err.Error())
		}
		fmt.Println(validOperations(), operation)
		if validity == false {
			t.Errorf("operation %v is not valid or not implemented", tsk.operation)
		}
		fmt.Printf("END   testing operation '%v':\n", operation)
	}
}

func TestTable(t *testing.T) {
	data := [][]string{
		{"Staged_s01e01_AUDIORUS20.m4a", "11.30", "0 0", "41<57<2"},
		{"Staged_s02e07_AUDIORUS20.m4a", "6.60", "0 0", "16<84<0"},
		{"The_Magicians_s02e02_AUDIORUS20.m4a", "12.00", "0 0", "39<61<0"},
		{"Scream_AUDIOENG51.m4a", "20.20", "-4 -4 0 -10 -9 -8", "69<25<7"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "RA", "channels", "stats"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table.SetBorder(false) // Set Border to false
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
