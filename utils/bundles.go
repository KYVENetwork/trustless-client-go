package utils

import (
	"encoding/json"
	"fmt"
	"github.com/KYVENetwork/trustless-client-go/types"
)

func GetFinalizedBundle(restEndpoint string, poolId int64, bundleId int64) (*types.FinalizedBundle, error) {
	raw, err := GetFromUrlWithBackoff(fmt.Sprintf(
		"%s/kyve/v1/bundles/%d/%d",
		restEndpoint,
		poolId,
		bundleId,
	))
	if err != nil {
		return nil, err
	}

	var finalizedBundle types.FinalizedBundle

	if err := json.Unmarshal(raw, &finalizedBundle); err != nil {
		return nil, err
	}

	return &finalizedBundle, nil
}
