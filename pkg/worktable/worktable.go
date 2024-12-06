package worktable

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	spreadsheet "gopkg.in/Iwark/spreadsheet.v2"
)

const (
	Col_Коментарий      = 0
	Col_Путь            = 1
	Col_ГТ              = 2
	Col_Т               = 3
	Col_Трейлер         = 4
	Col_П               = 5
	Col_Постеры         = 6
	Col_М               = 7
	Col_Наименование    = 8
	Col_С               = 9
	Col_З               = 10
	Col_О               = 11
	Col_I               = 12
	Col_Контрагент      = 13
	Col_Дата_публикации = 14

	spreadsheetID = "1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg"
	credentials   = `c:\Users\pemaltynov\.credentials\tablegrabber.json`
	sheetTitle    = "Амедиатека"
)

type WorkTable struct {
	service   *spreadsheet.Service
	sheet     *spreadsheet.Sheet
	Text      map[CellCoords]string
	Comment   map[CellCoords]string
	StartTime time.Time
}

func initWorkTable() (*WorkTable, error) {
	wt := WorkTable{}
	wt.StartTime = time.Now()
	wt.Text = make(map[CellCoords]string)
	wt.Comment = make(map[CellCoords]string)
	service, err := initService()
	if err != nil {
		return nil, err
	}
	wt.service = service
	return &wt, nil
}

func initService() (*spreadsheet.Service, error) {
	data, err := os.ReadFile(credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to read credetinals: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	if err != nil {
		return nil, fmt.Errorf("failed to create service configuration: %v", err)
	}
	client := conf.Client(context.TODO())
	service := spreadsheet.NewServiceWithClient(client)
	return service, nil
}

// func initSheet() (*spreadsheet.Sheet, error) {
// 	data, err := os.ReadFile(credentials)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read credetinals: %v", err)
// 	}
// 	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create service configuration: %v", err)
// 	}
// 	client := conf.Client(context.TODO())
// 	service := spreadsheet.NewServiceWithClient(client)
// 	spreadsheet, err := service.FetchSpreadsheet(spreadsheetID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch spreadsheet: %v", err)
// 	}
// 	sheet, err := spreadsheet.SheetByTitle(sheetTitle)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get sheet by title: %v", err)
// 	}
// 	return sheet, nil
// }

func (wt *WorkTable) Update() error {
	spreadsheet, err := wt.service.FetchSpreadsheet(spreadsheetID)
	if err != nil {
		return fmt.Errorf("failed to fetch spreadsheet: %v", err)
	}
	sheet, err := spreadsheet.SheetByTitle(sheetTitle)
	if err != nil {
		return fmt.Errorf("failed to get sheet by title: %v", err)
	}
	wt.Text = make(map[CellCoords]string)
	wt.Comment = make(map[CellCoords]string)
	for r, row := range sheet.Rows {
		for c, cell := range row {
			coords := Cell(r, c)
			wt.Text[coords] = cell.Value
			wt.Comment[coords] = cell.Note
		}
	}
	return nil
}

func (wt *WorkTable) ReadText(coords CellCoords) (string, error) {

	return "nil", nil
}

func (wt *WorkTable) EditCell(coords CellCoords) error {
	return nil
}

func (wt *WorkTable) EditComment(coords CellCoords) error {
	return nil
}

type CellCoords struct {
	RowNum    int
	ColumnNum int
}

func Cell(rowNum, colNum int) CellCoords {
	return CellCoords{RowNum: rowNum, ColumnNum: colNum}
}

/*
VIEW
==HEADER==========================
-Key map
-Table Filters
-Table stats
==TABLE==========================
| cellA1 | cellB1 | ... | cellO1 |
| cellA2 | cellB2 | ... | cellO2 |
| cellA3 | cellB3 | ... | cellO3 |
==SUMMARY========================

*/
