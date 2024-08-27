package main

import (
	"fmt"

	trustlesspai "github.com/KYVENetwork/trustless-client-go/trustlessapi"
)

func main() {
	_, err := trustlesspai.Get("http://localhost:4242/lava/block?height=1", "https://api.korellia.kyve.network")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("verified!")
}
