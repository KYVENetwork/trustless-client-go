package trustlesspai

import (
	"github.com/KYVENetwork/trustless-client-go/proof"
	"github.com/KYVENetwork/trustless-client-go/utils"
)

// returns the data item value from trustless api response, error if something fails
//
// 1. fetches the data item from the url
// 2. constructs a local merkle root from the data item
// 3. compares the local merkle root against the chains merkle root
//
// uses the chainRest provided for the request
// NOTE: 	This function will only use the chainRest for mainnet pools. If the pool is on a testnet, it will use the official testnet endpoints
func Get(url string, chainRest string) ([]byte, error) {

	rawResponse, proofHex, err := utils.GetFromUrl(url)

	if err != nil {
		return []byte{}, err
	}

	p, err := utils.DecodeProof(proofHex)
	if err != nil {
		return []byte{}, err
	}

	switch p.ChainId {
	case "kyve-1":
		chainRest = utils.RestEndpointMainnet
	case "kaon-1":
		chainRest = utils.RestEndpointKaon
	case "korellia-2":
		chainRest = utils.RestEndpointKorellia
	}

	err = proof.DataItemInclusionProof(rawResponse, p, chainRest)

	if err != nil {
		return []byte{}, err
	}

	return rawResponse, nil
}
