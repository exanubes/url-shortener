package cloudfront

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	visiturl "github.com/exanubes/url-shortener/internal/app/usecases/visit_url"
	"github.com/exanubes/url-shortener/internal/domain"
)

type LogProcessor struct {
	event_store visiturl.LinkEventStore
}

func NewHandler(event_store visiturl.LinkEventStore) *LogProcessor {
	return &LogProcessor{event_store}
}

func (handler LogProcessor) Handle(ctx context.Context, event events.KinesisEvent) (events.KinesisEventResponse, error) {
	var failures = make([]events.KinesisBatchItemFailure, 0)
	var errors = make([]string, 0)
	for _, record := range event.Records {
		fields := strings.Split(
			strings.TrimSpace(string(record.Kinesis.Data)),
			"\t",
		)

		if len(fields) != 6 {
			failures = append(failures, events.KinesisBatchItemFailure{
				ItemIdentifier: record.EventID,
			})
			errors = append(errors, fmt.Sprintf("Unexpected field count: %d", len(fields)))
			continue
		}

		status, err := strconv.Atoi(fields[2])

		if err != nil {
			failures = append(failures, events.KinesisBatchItemFailure{
				ItemIdentifier: record.EventID,
			})
			errors = append(errors, err.Error())
			continue
		}

		log_item := LogItem{
			Timestamp: parse_timestamp(fields[0]),
			IpAddress: fields[1],
			Status:    status,
			Method:    fields[3],
			URI:       fields[4],
			UserAgent: fields[5],
		}

		if log_item.Method == "GET" && log_item.Status == http.StatusTemporaryRedirect {
			err := handler.event_store.Visit(ctx, map_to_domain_event(log_item))

			if err != nil {
				failures = append(failures, events.KinesisBatchItemFailure{
					ItemIdentifier: record.EventID,
				})
				errors = append(errors, err.Error())
				continue
			}
		}

	}

	var err error

	if len(errors) != 0 {
		err = fmt.Errorf("ERRORS: %s", strings.Join(errors, "\n"))
	}

	return events.KinesisEventResponse{
		BatchItemFailures: failures,
	}, err
}

func map_to_domain_event(msg LogItem) domain.LinkVisited {

	short_code := strings.Split(strings.Split(msg.URI, "/")[1], "?")[0]

	return domain.LinkVisited{
		ShortCode: short_code,
		VisitedAt: msg.Timestamp,
		IpAddress: msg.IpAddress,
		UserAgent: msg.UserAgent,
	}
}
