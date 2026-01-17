package event

import (
	"context"
	"fmt"

	"github.com/exanubes/url-shortener/internal/domain"
)

type Bus struct {
	events  chan domain.LinkVisited
	handler func(domain.LinkVisited) error
}

func NewBus(handler func(domain.LinkVisited) error) *Bus {
	return &Bus{
		events:  make(chan domain.LinkVisited, 5),
		handler: handler,
	}
}

func (bus *Bus) Publish(event domain.LinkVisited) bool {
	select {
	case bus.events <- event:
	default:
		fmt.Println("Channel is full")
		return false
	}

	return true
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
