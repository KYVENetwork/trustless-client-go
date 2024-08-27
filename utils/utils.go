package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/KYVENetwork/trustless-client-go/types"
)

var (
	logger = TrustlessClientLogger("utils")
)

// GetFromUrl tries to fetch data from url with a custom User-Agent header
// returns the data, the proof and an error
func GetFromUrl(url string) ([]byte, string, error) {
	// Create a custom http.Client with the desired User-Agent header
	client := &http.Client{}

	// Create a new GET request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	// Perform the request
	response, err := client.Do(request)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, "", fmt.Errorf("got status code %d != 200", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	return data, response.Header.Get("x-kyve-proof"), nil
}

// GetFromUrlWithBackoff tries to fetch data from url with exponential backoff
func GetFromUrlWithBackoff(url string) (data []byte, err error) {
	for i := 0; i < BackoffMaxRetries; i++ {
		data, _, err = GetFromUrl(url)
		if err != nil {
			delaySec := math.Pow(2, float64(i))
			delay := time.Duration(delaySec) * time.Second

			logger.Error().Msg(fmt.Sprintf("failed to fetch from url %s, retrying in %d seconds", url, int(delaySec)))
			time.Sleep(delay)

			continue
		}

		// only log success message if there were errors previously
		if i > 0 {
			logger.Info().Msg(fmt.Sprintf("successfully fetch data from url %s", url))
		}
		return
	}

	logger.Error().Msg(fmt.Sprintf("failed to fetch data from url within maximum retry limit of %d", BackoffMaxRetries))
	return
}

func GetChainRest(chainId, chainRest string) string {
	if chainRest != "" {
		// trim trailing slash
		return strings.TrimSuffix(chainRest, "/")
	}

	// if no custom rest endpoint was given we take it from the chainId
	if chainRest == "" {
		switch chainId {
		case ChainIdMainnet:
			return RestEndpointMainnet
		case ChainIdKaon:
			return RestEndpointKaon
		case ChainIdKorellia:
			return RestEndpointKorellia
		default:
			panic(fmt.Sprintf("flag --chain-id has to be either \"%s\", \"%s\" or \"%s\"", ChainIdMainnet, ChainIdKaon, ChainIdKorellia))
		}
	}

	return ""
}

// DecodeProof decodes the proof of a data item from a byte array
// encodedProofString is the hex string of the proof
// encoded in big endian
// Structure:
// - 2  bytes: poolId (uint16)
// - 8  bytes: bundleId (uint64)
// - 16 bytes: chainId
// - 16 bytes: dataItemKey
// - 16 bytes: dataItemValueKey
// - Array of merkle nodes:
//   - 1 byte:  left (true/false)
//   - 32 bytes: hash (sha256)
//
// returns the proof as a struct
func DecodeProof(encodedProofString string) (*types.Proof, error) {

	encodedProof, err := hex.DecodeString(encodedProofString)
	if err != nil {
		return nil, err
	}

	if len(encodedProof) < 58 {
		return nil, fmt.Errorf("encoded proof is too short")
	}

	proof := &types.Proof{}

	proof.PoolId = int64(binary.BigEndian.Uint16(encodedProof[0:2]))
	proof.BundleId = int64(binary.BigEndian.Uint64(encodedProof[2:10]))
	// Convert the byte slice to a null-terminated string
	proof.ChainId = string(bytes.TrimRight(encodedProof[10:26], "\x00"))
	proof.DataItemKey = string(bytes.TrimRight(encodedProof[26:42], "\x00"))
	proof.DataItemValueKey = string(bytes.TrimRight(encodedProof[42:58], "\x00"))

	proofBytes := encodedProof[58:]

	for len(proofBytes) >= 33 {
		merkleNode := types.MerkleNode{}
		merkleNode.Left = proofBytes[0] == 1
		merkleNode.Hash = hex.EncodeToString(proofBytes[1:33])
		proof.Hashes = append(proof.Hashes, merkleNode)
		proofBytes = proofBytes[33:]
	}

	if len(proofBytes) != 0 {
		return nil, fmt.Errorf("invalid proof encoding")
	}

	return proof, nil
}
