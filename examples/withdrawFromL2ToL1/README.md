Note: To run this example, you need a sepolia endpoint.

Steps to run this example on sepolia: 
0. Have deployed your starknet-sepolia account already. ([ref.](https://book.starknet.io/ch02-05-01-start-with-sepolia.html)) 
1. rename ".env.template" to ".env.sepolia"
2. uncomment, and set all env variables
4. make sure you are in the "withdrawFromL1ToL2" directory
5. execute `go mod tidy`
6. execute `go run main.go`


### Workflow

1. Preprocessing
    - Load the [StarkGate address list](https://github.com/starknet-io/starknet-addresses/blob/master/bridged_tokens/sepolia.json).
    - Import the private key to control a deployed StarkNet account contract (L2 sender).
    - Import the private key to control an Ethereum EOA (L1 recipient).
    - Make sure the balance of the two accounts is enough.
2. L2
    - Call the [`initiate_withdraw`](https://docs.starknet.io/documentation/tools/starkgate_function_reference/#initiate_withdraw) on the L2 StarkGate Contract.
    - Wait for the `initiate_withdraw` transaction to be mined.
3. L1
    - Waiting for the matching L2->L1 message (`WITHDRAW_MESSAGE`) to exist and be ready to be consumed.
    - Trigger [`withdraw`](https://docs.starknet.io/documentation/tools/starkgate_function_reference/#withdraw) on the L1 StarkGate Contract.
    - In the end, we check the l1 recipient ETH balance in this step.