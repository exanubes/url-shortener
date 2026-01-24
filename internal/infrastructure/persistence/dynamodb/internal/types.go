package internal

import (
	"encoding/json"
	"time"
)

type LinkRow struct {
	PK          string          `dynamodbav:"PK"`
	SK          string          `dynamodbav:"SK"`
	Url         string          `dynamodbav:"url"`
	Shortcode   string          `dynamodbav:"shortcode"`
	PolicySpecs []PolicySpecDto `dynamodbav:"policy_specs"`
	CreatedAt   time.Time       `dynamodbav:"created_at"`
	ConsumedAt  time.Time       `dynamodbav:"consumed_at,omitempty"`
	Version     string          `dynamodbav:"version"`
}

type ConsumeSingleUseLinkParams struct {
	Shortcode  string    `dynamodbav:"shortcode"`
	ConsumedAt time.Time `dynamodbav:"consumed_at,omitempty"`
}

type PrimaryKey struct {
	PK string `dynamodbav:"PK"`
	SK string `dynamodbav:"SK"`
}

type PolicySpecDto struct {
	Kind   string          `dynamodbav:"kind"`
	Config json.RawMessage `dynamodbav:"config"`
}

type MaxAgeParamsDto struct {
	DurationNanoseconds time.Duration `json:"duration_nanoseconds"`
}
