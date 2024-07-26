package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/option"
	"google.golang.org/api/sheets/v4"
)

func main() {

	msg, err := os.ReadFile(`c:\Users\pemaltynov\client_secret.json`)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(msg))
	key := base64.StdEncoding.EncodeToString(msg)
	fmt.Println("encoded")
	fmt.Println(key)

	data, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Fatal("error:", err)
	}
	fmt.Println("decoded")
	fmt.Printf("%q\n", data)
}

func SaveDataToSheet(credBytes []byte) {
	ctx := context.Background()
	//autentification
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("%v", err)
	}
	//create client
	client := config.Client(ctx)
	//create service
	srv, err := sheets.NewService(ctx, option.WithClient(client))
}
