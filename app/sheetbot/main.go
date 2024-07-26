package main

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	SaveDataToSpreadSheet()
	fmt.Println()
}

func SaveDataToSpreadSheet() {

	ctx := context.Background()
	fmt.Println(os.Getenv("KEY_JSON_BASE64"))

	credBytes, err := b64.StdEncoding.DecodeString(os.Getenv("KEY_JSON_BASE64"))
	if err != nil {
		fmt.Println(credBytes)
		fmt.Println(string(credBytes))
		fmt.Println(credBytes)
		log.Error(err)
		return
	}

	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Error(err)
		return
	}

	client := config.Client(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	//https://docs.google.com/spreadsheets/d/1pIrdfVbUy5I9NF70USDMfgr_H8by3CwGJstFfWTTaug/edit#gid=1156983772

	sheetId := 176472191
	spreadSheetId := "1pIrdfVbUy5I9NF70USDMfgr_H8by3CwGJstFfWTTaug"

	resp, err := srv.Spreadsheets.Get(spreadSheetId).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil {
		log.Error(err)
		return
	}

	sheetName := ""
	for _, v := range resp.Sheets {
		props := v.Properties
		if props.SheetId == int64(sheetId) {
			sheetName = props.Title
			break
		}
	}
	fmt.Println("sheetname =", sheetName)
	records := sheets.ValueRange{
		Values: [][]interface{}{{"1", "abd", "ABD@mail.ru"}},
	}

	resp2, err := srv.Spreadsheets.Values.Append(spreadSheetId, sheetName, &records).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	if err != nil || resp2.HTTPStatusCode != 200 {
		log.Error(err)
		return
	}

	resp3, err := srv.Spreadsheets.Values.Get(spreadSheetId, "Sheet5!A1:C1").Do()
	fmt.Println(resp3.Range)
	fmt.Println(resp3.Values)
}
