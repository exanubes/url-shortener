package cloudfront

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type LogProcessor struct {
}

func NewHandler() *LogProcessor {
	return &LogProcessor{}
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

		var log_item LogItem

		log_item.Timestamp = Timestamp(parse_timestamp(fields[0]))
		log_item.IpAddress = fields[1]
		log_item.Status = fields[2]
		log_item.Method = fields[3]
		log_item.URI = fields[4]
		log_item.UserAgent = fields[5]

		fmt.Println("TIMESTAMP: ", log_item.Timestamp)
		fmt.Println("METHOD: ", log_item.Method)
		fmt.Println("URI: ", log_item.URI)
		fmt.Println("STATUS: ", log_item.Status)
		fmt.Println("USER_AGENT: ", log_item.UserAgent)
		fmt.Println("IP_ADDRESS: ", log_item.IpAddress)

	}

	if len(event.Records) == 0 {
		fmt.Println("NO RECORDS")
	}

	var err error

	if len(errors) != 0 {
		err = fmt.Errorf("ERRORS: %s", strings.Join(errors, "\n"))
	}

	return events.KinesisEventResponse{
		BatchItemFailures: failures,
	}, err
}
