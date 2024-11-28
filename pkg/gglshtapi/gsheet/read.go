package gsheet

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func ReadInSheet() {
	// Replace 'path/to/your/credentials.json' with the actual path to your downloaded JSON key file.
	credentialsFile := `c:\Users\pemaltynov\client_secret.json`

	// Replace 'your-spreadsheet-id' with the ID of the spreadsheet you want to read.
	//https://docs.google.com/spreadsheets/d/1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg/edit?gid=250314867#gid=250314867
	spreadsheetID := "your-spreadsheet-id"

	// Replace 'Sheet1!A1:B10' with the range of cells you want to read.
	readRange := "Sheet1!A1:B10"

	// Initialize Google Sheets API
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		log.Fatalf("Unable to initialize Sheets API: %v", err)
	}

	// Read data from the specified range
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// Print the values from the response
	if len(resp.Values) > 0 {
		fmt.Println("Data from sheet:")
		for _, row := range resp.Values {
			for _, cell := range row {
				fmt.Printf("%v\t", cell)
			}
			fmt.Println()
		}
	} else {
		fmt.Println("No data found.")
	}
}
