package main

import (
	"context"
	"fmt"
	"os"

	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

const (
	chainName          string = "sepolia"
)

func main() {
	// Load ENV vriables
	fmt.Println("Starting getTokenBalance example")
	godotenv.Load(fmt.Sprintf(".env.%s", chainName))
	base := os.Getenv("INTEGRATION_BASE")

	// Setup L1 Client
	c, err := ethrpc.DialContext(context.Background(), base)
	if err != nil {
		fmt.Println("Failed to connect to the client, did you specify the url in the .env.mainnet?")
		panic(err)
	}
	// Setup L2 Client
	clientv02 := rpc.NewProvider(c)
	fmt.Println("Established connection with the client")

	// Load Account

	// Make sure the balance is enough

	// Call the initiate_withdraw on the L2 StarkGate Contract.

	// Trigger withdraw on the L1 StarkGate Contract.

	// Check the L1 recipient ETH balance

}
