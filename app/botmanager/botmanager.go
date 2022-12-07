package main

import (
	"flag"
	"log"

	tgClient "github.com/Galdoba/ffstuff/app/botmanager/clients/telegram"
	eventconsumer "github.com/Galdoba/ffstuff/app/botmanager/consumer/event-consumer"
	"github.com/Galdoba/ffstuff/app/botmanager/events/telegram"
	"github.com/Galdoba/ffstuff/app/botmanager/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = `c:\Users\Public\Documents\Link_Storage\`
	batchSize   = 100
)

func main() {

	eventProccessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)
	log.Print("service started")
	consumer := eventconsumer.New(eventProccessor, eventProccessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service stopped", err)
	}

}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)
	flag.Parse()
	if *token == "" {
		log.Fatal("token not specified")
	}
	return *token
}
