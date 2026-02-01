package link_visited

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	visiturl "github.com/exanubes/url-shortener/internal/app/usecases/visit_url"
)

type LinkVisitedHandler struct {
	event_store visiturl.LinkEventStore
}

func NewHandler(event_store visiturl.LinkEventStore) *LinkVisitedHandler {
	return &LinkVisitedHandler{
		event_store: event_store,
	}
}

func (handler LinkVisitedHandler) Handle(ctx context.Context, event events.SQSEvent) events.SQSEventResponse {
	var failed_records []events.SQSBatchItemFailure
	for _, record := range event.Records {
		msg, err := parse_message(record.Body)

		if err != nil {
			failed_records = append(failed_records, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
			continue
		}

		handler.event_store.Visit(ctx, map_to_domain_event(msg))
	}

	return events.SQSEventResponse{
		BatchItemFailures: failed_records,
	}
}
