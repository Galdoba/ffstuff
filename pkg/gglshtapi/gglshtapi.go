package main

import (
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	spreadsheet "gopkg.in/Iwark/spreadsheet.v2"
)

const (
	spreadsheetID = ""
	//readRange     = "Sheet9!A:C"
	credentials = ``
)

func main() {
	data, err := os.ReadFile(credentials)
	checkError(err)
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)
	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	spreadsheet, err := service.FetchSpreadsheet(spreadsheetID)
	checkError(err)
	sheet, err := spreadsheet.SheetByTitle("Амедиатека")

	checkError(err)
	// for r, row := range sheet.Rows {
	// 	for c, cell := range row {
	// 		fmt.Println(r, c, cell.Value)
	// 	}
	// }
	//Комментарий
	// Update cell content

	sheet.Update(223, 0, "29.11.2024")
	sheet.Update(223, 1, "29.11.2024")
	sheet.Update(223, 2, "29.11.2024")
	sheet.Update(223, 3, "29.11.2024")
	sheet.Update(223, 4, "29.11.2024")
	sheet.Update(223, 5, "29.11.2024")
	sheet.Update(223, 6, "29.11.2024")
	sheet.Update(223, 7, "29.11.2024")

	// Make sure call Synchronize to reflect the changes
	err = sheet.Synchronize()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
