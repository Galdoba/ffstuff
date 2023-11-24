package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
)

type tableData struct {
	data         [][]string
	filters      map[string]string
	cursorRow    int
	cursorCol    int
	selectedRows []int
	selectedCols []int
}

func newTableData(path string) tableData {
	tb := tableData{}
	f, _ := os.Open(path)
	defer f.Close()
	reader := csv.NewReader(f)
	tb.data, _ = reader.ReadAll()

	return tb
}

type content struct {
	columns int
	rows    int
	cells   map[string]*cell
}

func newContent(data [][]string) *content {
	cn := content{}
	cn.cells = make(map[string]*cell)
	cn.update(data)
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
	if i > 2 && len(lText) > i {
		lText = lText[:i-2]
		lText = append(lText, ".")
		lText = append(lText, ".")
	}
	for len(lText) < i {
		lText = append(lText, " ")
	}
	return merge(lText)
}

func (cn *content) update(data [][]string) {
	columnLen := columnSizes(data)
	cn.columns = len(columnLen)
	//rowList:
	for r, line := range data {
		//	cellList:
		for c, rawtext := range line {
			crd := coord(r, c)

			if _, ok := cn.cells[crd.String()]; !ok {
				cll := newCell(r, c, rawtext)

				//maxLen := columnLen[c]

				ltrText := letters(rawtext)
				// for len(ltrText) < maxLen {
				// 	ltrText = append(ltrText, " ")
				// }
				cll.fmtText = merge(ltrText)

				cn.cells[crd.String()] = cll

			}

		}
		cn.rows++
	}
}

func (cn *content) Cell(key string) string {

	return cn.cells[key].fmtText
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
