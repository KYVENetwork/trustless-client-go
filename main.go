package main

import (
	"fmt"

	"github.com/KYVENetwork/trustless-client-go/trustlessapi"
)

func main() {
	_, err := trustlessapi.Get("http://localhost:4242/cronos-zkevm/value?height=1", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = trustlessapi.Get("http://localhost:4242/celestia/GetSharesByNamespace?height=800005&namespace=AAAAAAAAAAAAAAAAAAAAAAAAAIZiad33fbxA7Z0%3D", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = trustlessapi.Get("http://localhost:4242/ethereum/beacon/blob_sidecars?block_height=19426587", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = trustlessapi.Get("http://localhost:4242/lava/block?height=1", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = trustlessapi.Get("http://localhost:4242/lava/block_results?height=1", "")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("verified")
}
