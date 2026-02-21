package clock

import (
	"github.com/exanubes/url-shortener/internal/domain"
	"time"
)

type Clock struct{}

var _ domain.Clock = (*Clock)(nil)

func NewClock() *Clock {
	return &Clock{}
}

func (Clock) Now() time.Time {
	return time.Now()
}
