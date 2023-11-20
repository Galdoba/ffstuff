package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
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
	columns []columnData
	cells   map[string]*cell
}

func newContent(data [][]string, preset string) *content {
	cn := content{}
	cn.cells = make(map[string]*cell)
	cn.update(data, preset)
	return &cn
}

func letters(s string) []string {
	return strings.Split(s, "")
}

func merge(sl []string) string {
	return strings.Join(sl, "")
}

func widen(text string, i int) string {
	lText := letters(text)
	if i > 2 && len(lText) > i-2 {
		lText = lText[:i-2]
		lText = append(lText, ".")
		lText = append(lText, ".")
	}
	for len(lText) < i {
		lText = append(lText, " ")
	}
	return merge(lText)
}

func (cn *content) update(data [][]string, preset string) {
	columnLen := columnSizes(data)
	for r, line := range data {
		for c, rawtext := range line {
			crd := coord(r, c)
			if _, ok := cn.cells[crd.String()]; !ok {
				cll := newCell(r, c, rawtext)
				ltrText := letters(rawtext)
				for len(ltrText) < columnLen[r] {
					ltrText = append(ltrText, " ")
				}
				cll.fmtText = merge(ltrText)
				cn.cells[crd.String()] = cll
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
	canBrake bool
	hidden   bool
	maxWidth int
	rawText  string
	fmtText  string
	colStyle *color.Style256
	row      int
	col      int
}
