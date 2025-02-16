package main

import (
	"context"
	"flag"
	"log"
	tgClient "read-it-later-bot/clients/telegram"
	event_consumer "read-it-later-bot/consumer/event-consumer"
	telegram "read-it-later-bot/events/telegram"
	"read-it-later-bot/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	storagePath       = "storage/files-storage"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {

	// s := files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatalf("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}
	eventProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
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
