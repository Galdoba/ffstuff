package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type tableData struct {
	filepath     string
	tableName    string
	data         [][]string
	filters      map[string]string
	cursorRow    int
	cursorCol    int
	selectedRows []int
	selectedCols []int
	hiddenRows   map[int]bool
	hiddenCols   map[int]bool
}

func newTableData(path string) tableData {
	tb := tableData{}
	tb.filepath = path
	file, _ := os.Open(tb.filepath)

	reader := csv.NewReader(file)
	tb.data, _ = reader.ReadAll()
	tb.hiddenRows = make(map[int]bool)
	tb.hiddenCols = make(map[int]bool)
	//fmt.Println(tb.data[4][13])

	return tb
}

type content struct {
	cells map[coordinates]*cell
}

func newContent(data [][]string) *content {
	cn := content{}
	cn.cells = make(map[string]*cell)
	cn.update(data)
	return &cn
}

func (cn *content) update(data [][]string) {
	columnLen := columnSizes(data)
	for r, line := range data {
		for c, rawtext := range line {
			crd := coord(r, c)
			if _, ok := cn.cells[crd]; !ok {
				cn.cells[crd] = newCell(r, c, rawtext)
			}

		}
	}
}

type coordinates struct {
	row int
	col int
}

func coord(r, c int) coordinates {
	return coordinates{r, c}
}

func (c *coordinates) String() string {
	return fmt.Sprintf("R%vC%v", c.row, c.col)
}

func sameCell(a, b coordinates) bool {
	if a.col != b.col {
		return false
	}
	if a.row != b.row {
		return false
	}
	return true
}

func sameRow(a, b coordinates) bool {
	if a.row != b.row {
		return false
	}
	return true
}

func sameCol(a, b coordinates) bool {
	if a.col != b.col {
		return false
	}
	return true
}

func newCell(r, c int, rawText string) *cell {
	cl := cell{}
	cl.row = r
	cl.col = c
	cl.rawText = rawText
	cl.maxWidth = -1
	return &cl
}

type cell struct {
	letters  []string
	canTrim  bool
	hidden   bool
	maxWidth int
	rawText  string
	fmtText  string
	fgCol    int
	bgCol    int
	row      int
	col      int
}
