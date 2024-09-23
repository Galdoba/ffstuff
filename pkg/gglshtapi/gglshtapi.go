package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
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
	// ctx := context.Background()
	// //autentification
	// config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }
	//create client
	//client := config.Client(ctx)
	//create service
	//srv, err := sheets.NewService(ctx)
}
