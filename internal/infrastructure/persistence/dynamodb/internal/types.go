package internal

import (
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
)

type LinkRow struct {
	PK          string              `dynamodbav:"PK"`
	SK          string              `dynamodbav:"SK"`
	Url         string              `dynamodbav:"Url"`
	Shortcode   string              `dynamodbav:"Shortcode"`
	PolicySpecs []domain.PolicySpec `dynamodbav:"PolicySpecs"`
	CreatedAt   time.Time           `dynamodbav:"CreatedAt"`
	ConsumedAt  time.Time           `dynamodbav:"ConsumedAt,omitempty"`
	Version     int                 `dynamodbav:"Version"`
}

type ConsumeSingleUseLinkParams struct {
	Shortcode  string    `dynamodbav:"Shortcode"`
	ConsumedAt time.Time `dynamodbav:"ConsumedAt,omitempty"`
}

type PrimaryKey struct {
	PK string `dynamodbav:"PK"`
	SK string `dynamodbav:"PK"`
}
