package internal

import (
	"encoding/json"
	"fmt"

	"github.com/exanubes/url-shortener/internal/domain"
)

func CreateLinkMetaPartitionKey(id string) PrimaryKey {
	pk := fmt.Sprintf("LINK#%s", id)
	sk := "META"

	return PrimaryKey{PK: pk, SK: sk}
}

func ConvertParamsToDto(params domain.PolicyParams) any {
	switch params := params.(type) {
	case domain.MaxAgeParams:
		return MaxAgeParamsDto{DurationNanoseconds: params.TTL}
	}

	return nil
}

func SerializePolicies(policy_specs []domain.PolicySpec) ([]PolicySpecDto, error) {
	policy_spec_dtos := make([]PolicySpecDto, len(policy_specs))
	for index, policy := range policy_specs {
		var err error
		var config json.RawMessage
		params := ConvertParamsToDto(policy.Params)
		if params != nil {
			config, err = json.Marshal(params)
			if err != nil {
				return nil, err
			}
		} else {
			config = []byte("{}")
		}

		policy_spec_dtos[index] = PolicySpecDto{
			Kind:   string(policy.Kind),
			Config: config,
		}

	}

	return policy_spec_dtos, nil
}

func DeserializePolicies(policies []PolicySpecDto) ([]domain.PolicySpec, error) {
	specs := make([]domain.PolicySpec, len(policies))

	for index, p := range policies {
		switch domain.PolicyKind(p.Kind) {
		case domain.PolicyKind_SingleUse:
			specs[index] = domain.PolicySpec{Kind: domain.PolicyKind_SingleUse, Params: domain.SingleUseParams{}}
		case domain.PolicyKind_MaxAge:
			var config MaxAgeParamsDto
			if err := json.Unmarshal(p.Config, &config); err != nil {
				return nil, err
			}

			specs[index] = domain.PolicySpec{Kind: domain.PolicyKind_MaxAge, Params: domain.MaxAgeParams{TTL: config.DurationNanoseconds}}
		}
	}

	return specs, nil
}
