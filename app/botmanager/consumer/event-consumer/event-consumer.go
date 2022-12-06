package eventconsumer

import (
	"log"
	"time"

	"github.com/Galdoba/ffstuff/app/botmanager/events"
)

type ConsumerUnit struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(f events.Fetcher, p events.Processor, bs int) ConsumerUnit {
	return ConsumerUnit{
		fetcher:   f,
		processor: p,
		batchSize: bs,
	}
}

func (c *ConsumerUnit) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERROR] consumer: %s", err.Error())
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := c.handleEvents(gotEvents); err != nil {
			log.Printf("[ERROR] consumer: %s", err.Error())
			continue
		}
	}
}

func (c *ConsumerUnit) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("new event: %s", event.Text)
		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", event.Text)
			continue
		}
	}
	return nil
}
