package mapper

import (
	"errors"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/http/routes/dto"
)

func ToCreateLinkCommand(payload dto.CreateUrlRequest) (domain.CreateLinkCommand, error) {
	if payload.Url == "" {
		return domain.CreateLinkCommand{}, errors.New("Url cannot be empty")
	}

	url, err := domain.NewUrl(payload.Url)

	if err != nil {
		return domain.CreateLinkCommand{}, err
	}

	if payload.ExpiresAfter.Value == 0 {
		payload.ExpiresAfter.Value = 30
		payload.ExpiresAfter.Unit = "day"
	}

	if payload.ExpiresAfter.Unit == "" {
		payload.ExpiresAfter.Value = 30
		payload.ExpiresAfter.Unit = "day"
	}

	max_visits := 0
	max_age := time.Duration(payload.ExpiresAfter.Value) * parse_unit(payload.ExpiresAfter.Unit)

	if payload.OneTimeLink {
		max_visits = 1
	}

	policy_settings, err := domain.NewPolicySettings(max_visits, max_age)

	if err != nil {
		return domain.CreateLinkCommand{}, err
	}

	return domain.CreateLinkCommand{
		Url:            url,
		PolicySettings: policy_settings,
	}, nil
}

func parse_unit(unit dto.ExpirationUnit) time.Duration {
	switch unit {
	case "day":
		return 24 * time.Hour
	case "hour":
		return time.Hour
	case "minute":
		return time.Minute
	}

	return 24 * time.Hour
}
