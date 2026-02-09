package cloudfront

import (
	"strconv"
	"strings"
	"time"
)

func parse_timestamp(timestamp string) time.Time {
	parts := strings.Split(timestamp, ".")
	secStr, nsecStr := parts[0], parts[1]

	sec, _ := strconv.ParseInt(secStr, 10, 64)
	nsec, _ := strconv.ParseInt(nsecStr+"000000", 10, 64) // Pad to nanoseconds

	return time.Unix(sec, nsec)
}
