package internal

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
)

const SHARD_COUNT = 10

var buckets = map[string]string{"HOUR": "2006-01-02T15", "DAY": "2006-01-02", "MONTH": "2006-01", "YEAR": "2006"}

func CreateLinkVisitBucketPartitionKeys(id string, visited_at time.Time) []PrimaryKey {
	pk := fmt.Sprintf("AGG#%s", id)
	var keys []PrimaryKey
	seed, _ := rand_hex()
	shard := pick_shard(seed, SHARD_COUNT)

	for bucket, layout := range buckets {
		keys = append(keys, PrimaryKey{
			PK: fmt.Sprintf("%s#%s#%s", pk, bucket, visited_at.Format(layout)),
			SK: fmt.Sprintf("SHARD#%02d", shard),
		})
	}

	return keys
}

func CreateLinkVisitPartitionKey(id string, visited_at time.Time) PrimaryKey {
	now := time.Now()
	pk := fmt.Sprintf("VISIT#%s#%s", id, now.Format("2006-01-02"))
	dedup, _ := rand_hex()
	sk := fmt.Sprintf("TS#%s#%s", visited_at.UTC().Format(time.RFC3339Nano), dedup)
	return PrimaryKey{PK: pk, SK: sk}
}

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

// Pads digits with a zero to keep lexographical order
func pad_num_zero(input int) string {
	return fmt.Sprintf("%02d", input)
}

func rand_hex() (string, error) {
	b := make([]byte, 8) // 64 bits
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// PickShard returns a shard number in range [0, shardCount).
// The same seed will always map to the same shard.
func pick_shard(seed string, shardCount int) int {
	if shardCount <= 1 {
		return 0
	}

	sum := sha1.Sum([]byte(seed))

	// Take first 4 bytes → uint32
	random := binary.BigEndian.Uint32(sum[0:4])

	return int(random % uint32(shardCount))
}
