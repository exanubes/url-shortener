package cloudfront

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type LogProcessor struct {
}

func NewHandler() *LogProcessor {
	return &LogProcessor{}
}

func (handler LogProcessor) Handle(ctx context.Context, event events.KinesisEvent) (events.KinesisEventResponse, error) {
	var failures = make([]events.KinesisBatchItemFailure, 0)

	return events.KinesisEventResponse{
		BatchItemFailures: failures,
	}, nil
}
