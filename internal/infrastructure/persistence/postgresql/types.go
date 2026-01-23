package postgresql

import "time"

type max_age_params_dto struct {
	DurationNanoseconds time.Duration `json:"duration_nanoseconds"`
}
