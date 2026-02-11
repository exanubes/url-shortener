package cloudfront

import (
	"time"
)

type LogItem struct {
	Timestamp time.Time
	IpAddress string
	Method    string
	URI       string
	Status    int
	UserAgent string
}
