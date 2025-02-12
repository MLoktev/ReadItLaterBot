package main

import (
	"flag"
	"log"
	tgClient "read-it-later-bot/clients/telegram"
	event_consumer "read-it-later-bot/consumer/event-consumer"
	telegram "read-it-later-bot/events/telegram"
	files "read-it-later-bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage/files-storage"
	batchSize   = 100
)

func main() {
	eventProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("Service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

}

func mustToken() string {
	// bot -tg-bot-token 'my token'
	// flag.String возвращает ссылку, token может быть nil.
	// Но тут так не будет, второй параметр - это дефолтное значение
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
