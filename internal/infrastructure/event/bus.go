package event

import (
	"context"
	"fmt"

	"github.com/exanubes/url-shortener/internal/domain"
)

type Bus struct {
	events  chan domain.Event
	handler func(domain.Event) error
}

func NewBus(handler func(domain.Event) error) *Bus {
	return &Bus{
		events:  make(chan domain.Event, 5),
		handler: handler,
	}
}

func (bus *Bus) Publish(ctx context.Context, event domain.Event) error {
	select {
	case bus.events <- event:
	default:
		return fmt.Errorf("Channel is full")
	}

	return nil
}

func (bus *Bus) Start(ctx context.Context) {
	go func() {
		defer func() {

			err := recover()
			if err != nil {
				fmt.Println(err)
			}

		}()

		for {
			select {
			case event, ok := <-bus.events:
				if !ok {
					return
				}
				err := bus.handler(event)

				if err != nil {
					fmt.Println(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
