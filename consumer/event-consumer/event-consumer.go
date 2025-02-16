package event_consumer

import (
	"context"
	"log"
	"read-it-later-bot/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		// По хорошему тут надо в Fetcher встроить retry (повтор запроса при ошибке, обычно 3 раза)
		gotEvents, err := c.fetcher.Fetch(context.Background(), c.batchSize)
		if err != nil {
			log.Printf("[ERROR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(context.Background(), gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

/*
	Проблемы:
	1. Потеря событий: ретраи, возращение в хранилище, фоллбэк (сохранять в файл или в Оперативку)
		подтверждение Fetcher-а
	2. Обработка всей пачки: останавливаться после первой ошибки, счётчик ошибок, параллельная обработка (sync.WaitGroup())
*/

func (c *Consumer) handleEvents(ctx context.Context, events []events.Event) error {
	for _, event := range events {
		log.Printf("got new events: %s", event.Text)

		// Нужен механизм ретрайя и бэкапа
		if err := c.processor.Process(ctx, event); err != nil {
			log.Printf("can't handle event: %s", err.Error())
			continue
		}
	}

	return nil
}
