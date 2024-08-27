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

func GetMerkleRoot(proof *types.Proof, leafHash [32]byte) (string, error) {
	var parentHash [32]byte
	for _, merkleNode := range proof.Hashes {
		nodeHash, err := hex.DecodeString(merkleNode.Hash)
		if err != nil {
			return "", fmt.Errorf("failed to decode Merkle node hash: %v", err)
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
	return merkleRoot, nil
}

// DataItemInclusionProof proves the validity of a given Data Item.
// Therefore, it regenerates the Merkle root with the provided Merkle nodes in order to
// compare the locally computed hash with the one stored on-chain.
// It returns true if the given Data Item is the actually validated one by KYVE and false if not.
func DataItemInclusionProof(dataItem []byte, proof *types.Proof, endpoint string) error {

	decodedDataItem := map[string]json.RawMessage{}
	err := json.Unmarshal(dataItem, &decodedDataItem)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data item: %v", err)
	}

	constructedDataItem := map[string]interface{}{
		"key":   proof.DataItemKey,
		"value": decodedDataItem[proof.DataItemValueKey],
	}

	dataItemBytes, err := json.Marshal(constructedDataItem)
	if err != nil {
		return fmt.Errorf("failed to marshal constructed data item: %v", err)
	}

	// 1. Compute the Merkle root from the proof
	leafHash := sha256.Sum256(dataItemBytes)
	merkleRoot, err := GetMerkleRoot(proof, leafHash)

	fmt.Println(hex.EncodeToString(leafHash[:]))

	if err != nil {
		return fmt.Errorf("failed to compute Merkle root: %v", err)
	}

	// 2. Compare local Merkle root with Merkle root stored on-chain
	restEndpoint := utils.GetChainRest(proof.ChainId, endpoint)

	bundle, err := utils.GetFinalizedBundle(restEndpoint, proof.PoolId, proof.BundleId)
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
