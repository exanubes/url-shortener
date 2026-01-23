package postgresql

import "github.com/exanubes/url-shortener/internal/domain"

func convert_params_to_dto(params domain.PolicyParams) any {
	switch params := params.(type) {
	case domain.MaxAgeParams:
		return max_age_params_dto{DurationNanoseconds: params.TTL}
	}

	return nil
}
