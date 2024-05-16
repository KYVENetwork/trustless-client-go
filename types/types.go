package types

import (
	"encoding/json"
	"fmt"
)

type BundleSummary struct {
	FromSlot   int64  `json:"from_slot,omitempty"`
	MerkleRoot string `json:"merkle_root"`
	ToSlot     int64  `json:"to_slot,omitempty"`
}

// FinalizedBundle is the bundle that is stored on the KYVE chain.
// It includes the storage ID to retrieve the archived data and
// a summary. This summary includes the Merkle root the Trustless
// Client uses to validate that the data item was included in the
// bundle stored on-chain.
type FinalizedBundle struct {
	BundleSummary     string `json:"bundle_summary,omitempty"`
	CompressionId     string `json:"compression_id,omitempty"`
	DataHash          string `json:"data_hash,omitempty"`
	FromKey           string `json:"from_key,omitempty"`
	Id                string `json:"id,omitempty"`
	StorageId         string `json:"storage_id,omitempty"`
	StorageProviderId string `json:"storage_provider_id,omitempty"`
	ToKey             string `json:"to_key,omitempty"`
}

type FinalizedBundlesResponse = struct {
	FinalizedBundles []FinalizedBundle `json:"finalized_bundles"`
	Pagination       Pagination        `json:"pagination"`
}

type MerkleNode struct {
	Hash string `json:"hash"`
	Left bool   `json:"left"`
}

type Pagination struct {
	NextKey []byte `json:"next_key"`
}

// TrustlessDataItem is the response type of the Trustless API.
type TrustlessDataItem struct {
	ChainId  string          `json:"chainId"`
	BundleId int64           `json:"bundleId"`
	PoolId   int64           `json:"poolId"`
	Proof    []MerkleNode    `json:"proof"`
	Value    json.RawMessage `json:"value"`
}

type MerkleRootNotValidError struct {
	Constructed string
	OnChain     string
}

func (mrnv MerkleRootNotValidError) Error() string {
	return fmt.Sprintf("mismatch: local Merkle root (%v) != chain Merkle root (%v)", mrnv.Constructed, mrnv.OnChain)
}
