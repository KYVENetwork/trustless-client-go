package proof

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/KYVENetwork/trustless-client-go/types"
	"github.com/KYVENetwork/trustless-client-go/utils"
)

var (
	logger = utils.TrustlessClientLogger("proof")
)

// DataItemInclusionProof proves the validity of a given Data Item.
// Therefore, it regenerates the Merkle root with the provided Merkle nodes in order to
// compare the locally computed hash with the one stored on-chain.
// It returns true if the given Data Item is the actually validated one by KYVE and false if not.
func DataItemInclusionProof(trustlessDataItem types.TrustlessDataItem, endpoint string) error {
	// 1. Check if Merkle leaf in proof matches hash of local data item
	leafHash := sha256.Sum256(trustlessDataItem.Value)

	var parentHash [32]byte
	for _, merkleNode := range trustlessDataItem.Proof {
		nodeHash, err := hex.DecodeString(merkleNode.Hash)
		if err != nil {
			return fmt.Errorf("failed to decode Merkle node hash: %v", err)
		}

		if merkleNode.Left {
			combined := append(leafHash[:], nodeHash[:]...)
			parentHash = sha256.Sum256(combined)
		} else {
			combined := append(nodeHash[:], leafHash[:]...)
			parentHash = sha256.Sum256(combined)
		}

		leafHash = parentHash
	}
	merkleRoot := hex.EncodeToString(parentHash[:])

	// 2. Compare local Merkle root with Merkle root stored on-chain
	restEndpoint := utils.GetChainRest(trustlessDataItem.ChainId, endpoint)

	bundle, err := utils.GetFinalizedBundle(restEndpoint, trustlessDataItem.PoolId, trustlessDataItem.BundleId)
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("Failed to get finalized bundle")
		return err
	}

	var bundleSummary *types.BundleSummary
	if err = json.Unmarshal([]byte(bundle.BundleSummary), &bundleSummary); err != nil {
		return fmt.Errorf("failed to unmarshal bundle summary: %v", err)
	}

	if bundleSummary.MerkleRoot != merkleRoot {
		logger.Fatal().Msg("Mismatch: Local Merkle root != Chain Merkle root")
		return types.MerkleRootNotValidError{
			OnChain:     bundleSummary.MerkleRoot,
			Constructed: merkleRoot,
		}
	}

	return nil
}
