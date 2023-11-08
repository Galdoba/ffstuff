package main

type tableData struct {
	filepath     string
	tableName    string
	data         []byte
	filters      map[string]string
	cursorRow    int
	cursorCol    int
	selectedRows []int
	selectedCols []int
	hiddenRows   []int
	hiddenCols   []int
}

func newTableData(path string) tableData {
	tb := tableData{}
	tb.data = []byte("R1C1,R1C2,R1C3\nR2C1,R2C2,R2C3\nR3C1,R3C2,R3C3")
	return tableData{}
}
