package main

import (
	"encoding/csv"
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

type coord struct {
	r int
	c int
}

func coordinates(r, c int) coord {
	return coord{r, c}
}

func RowOf(a coord) int {
	return a.r
}

func ColOf(a coord) int {
	return a.c
}

func coordMatch(a, b coord) bool {
	if a.r != b.r {
		return false
	}
	if a.c != b.c {
		return false
	}
	return true
}

func sameRow(a, b coord) bool {
	if a.r != b.r {
		return false
	}
	return true
}

func sameCol(a, b coord) bool {
	if a.c != b.c {
		return false
	}
	return true
}
