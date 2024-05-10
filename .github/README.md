# Trustless Client (Go)

🤝 **Enables the seamless integration of KYVE's [Trustless API](https://github.com/KYVENetwork/trustless-api)** 🤝

## Overview

The Trustless Client enables simple and efficient use of KYVE's Trustless API. It implements the proof of inclusion of a received data item from the Trustless API. This ensures that the local data matches the data validated by KYVE.

Therefore, the Merkle nodes of the received proof for a specific data item are used to generate a Merkle root. This locally computed Merkle root is then compared to the root stored on the chain for comparison. The client returns either true or false, depending on whether the Merkle roots match.

## Usage

The `DataItemInclusionProof` serves all the required logic. It expects a TrustlessDataItem, which is similar to the response type of the Trustless API. Besides this, it expects an endpoint that is used to compare the local computed Merkle root with the one stored on-chain. Although the endpoint is not required, 
it's recommended to specify either your own KYVE node endpoint or a KYVE node you're trusting. By default, KYVE's official node endpoints are used.