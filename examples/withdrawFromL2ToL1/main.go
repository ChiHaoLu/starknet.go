package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

const (
	chainName              string = "sepolia"
	accountContractVersion int    = 0 //Replace with the cairo version of your account contract
)

/*
 * @author ChiHaoLu (chihaolu.me)
 * @title WithdrawFromL2ToL1 example
 */
func main() {
	// Load ENV vriables
	fmt.Println("Starting withdrawFromL2ToL1 example")
	godotenv.Load(fmt.Sprintf(".env.%s", chainName))

	// Setup L1 Client
	l1Client, err := ethrpc.DialContext(context.Background(), os.Getenv("L1_RPC_URL"))
	if err != nil {
		fmt.Println("Failed to connect to the client, did you specify the url in the .env.sepolia?")
		panic(err)
	}

	// Setup L2 Client
	l2ClientBase, err := ethrpc.DialContext(context.Background(), os.Getenv("L2_RPC_URL"))
	if err != nil {
		fmt.Println("Failed to connect to the client, did you specify the url in the .env.sepolia?")
		panic(err)
	}
	l2Client := rpc.NewProvider(l2ClientBase)

	// Load L1 Account
	l1PrivKey := os.Getenv("L1_PRIV_KEY")
	l1Address := GetEthereumAddressFromHexPrivKey(l1PrivKey)
	fmt.Println("Load L1 account: ", l1Address)

	// Load L2 Account
	l2Account := ImportAccount(l2Client, os.Getenv("L2_PUB_KEY"), os.Getenv("L2_PRIV_KEY"), accountContractVersion)
	fmt.Println("Load L2 account: ", l2Account.AccountAddress)

	// Make sure the balance is enough
	L1Balance, err := GetL1ETHBalance(l1Client, l1Address)
	if err != nil {
		panic(err)
	}
	if L1Balance.Cmp(big.NewFloat(0)) <= 0 {
		panic("You don't have enough ether in l1 accounts")
	}
	fmt.Printf("L1 Balance: %f ETH\n", L1Balance)
	L2Balance, err := GetL2ETHBalance(l2Client, l2Account.AccountAddress, Sepolia.tokenList[0].l2TokenAddr, Sepolia.tokenList[0].decimals)
	if err != nil {
		panic(err)
	}
	if L2Balance.Cmp(big.NewFloat(0)) <= 0 {
		panic("You don't have enough ether in l2 accounts")
	}
	fmt.Printf("L2 Balance: %f ETH\n", L2Balance)

	// Call the initiate_withdraw on the L2 StarkGate Contract.
	var calldata []*felt.Felt

	l1_recipient, _ := utils.HexToFelt(l1Address.String())
	amount := utils.BigIntToFelt(big.NewInt(100000000000000))
	calldata = append(calldata, l1_recipient)
	calldata = append(calldata, amount)
	InvokeL2Contract(l2Client, l2Account, Sepolia.tokenList[0].l2BridgeAddr, "initiate_withdraw", "0x9184e72a000", calldata)

	// Trigger withdraw on the L1 StarkGate Contract.

	// Check the L1 recipient ETH balance
}

/* ======================
 * Chain Information
 * ====================== */

type ChainInfo struct {
	name      string
	rpcUrl    string
	tokenList []TokenInfo
}

type TokenInfo struct {
	name         string
	symbol       string
	decimals     int
	l1TokenAddr  string
	l2TokenAddr  string
	l1BridgeAddr string
	l2BridgeAddr string
}

// Ref: https://github.com/starknet-io/starknet-addresses/blob/master/bridged_tokens/sepolia.json
var Sepolia = ChainInfo{
	name:   "Sepolia",
	rpcUrl: "",
	tokenList: []TokenInfo{
		{
			name:         "Ether",
			symbol:       "ETH",
			decimals:     18,
			l1TokenAddr:  "0x0000000000000000000000000000000000000000",
			l2TokenAddr:  "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
			l1BridgeAddr: "0x8453FC6Cd1bCfE8D4dFC069C400B433054d47bDc",
			l2BridgeAddr: "0x04c5772d1914fe6ce891b64eb35bf3522aeae1315647314aac58b01137607f3f",
		},
		{
			name:         "StarkNet Token",
			symbol:       "STRK",
			decimals:     18,
			l1TokenAddr:  "0xCa14007Eff0dB1f8135f4C25B34De49AB0d42766",
			l2TokenAddr:  "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d",
			l1BridgeAddr: "0xcE5485Cfb26914C5dcE00B9BAF0580364daFC7a4",
			l2BridgeAddr: "0x0594c1582459ea03f77deaf9eb7e3917d6994a03c13405ba42867f83d85f085d",
		},
	},
}

/* ======================
 * L2 Helper Function
 * ====================== */

func ImportAccount(l2Client *rpc.Provider, pubKey string, l2PrivKey string, accountContractVersion int) *account.Account {
	privKeyInBigInt := new(big.Int)
	privKeyInBigInt.SetString(l2PrivKey, 16)
	ks := account.SetNewMemKeystore(pubKey, privKeyInBigInt)

	accountAddressFelt, err := new(felt.Felt).SetString(os.Getenv("L2_ACCOUNT_ADDR"))
	if err != nil {
		panic("Error casting accountAddress to felt")
	}
	l2Account, err := account.NewAccount(l2Client, accountAddressFelt, os.Getenv("L2_PUB_KEY"), ks, accountContractVersion)
	if err != nil {
		panic(err)
	}
	return l2Account
}

func GetL2ETHBalance(rpcclient *rpc.Provider, accountAddressInFelt *felt.Felt, tokenAddress string, decimals int) (*big.Float, error) {
	tokenAddressInFelt, err := utils.HexToFelt(tokenAddress)
	if err != nil {
		fmt.Println("Failed to transform the token contract address, did you give the hex address?")
		panic(err)
	}
	tx := rpc.FunctionCall{
		ContractAddress:    tokenAddressInFelt,
		EntryPointSelector: utils.GetSelectorFromNameFelt("balanceOf"),
		Calldata:           []*felt.Felt{accountAddressInFelt},
	}
	callResp, err := rpcclient.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
	if err != nil {
		panic(err.Error())
	}

	decimalsInBig := big.NewInt(int64(decimals))
	floatValue := new(big.Float).SetInt(utils.FeltToBigInt(callResp[0]))
	floatValue.Quo(floatValue, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), decimalsInBig, nil)))

	return floatValue, nil
}

func InvokeL2Contract(l2Client *rpc.Provider, accnt *account.Account, contractAddress string, contractMethod string, maxFee string, calldata []*felt.Felt) {
	maxfee, err := utils.HexToFelt(maxFee)
	if err != nil {
		panic(err.Error())
	}
	// Getting the nonce from the account
	nonce, err := accnt.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, accnt.AccountAddress)
	if err != nil {
		panic(err.Error())
	}

	// Building the InvokeTx struct
	InvokeTx := rpc.InvokeTxnV1{
		MaxFee:        maxfee,
		Version:       rpc.TransactionV1,
		Nonce:         nonce,
		Type:          rpc.TransactionType_Invoke,
		SenderAddress: accnt.AccountAddress,
	}

	// Converting the contractAddress from hex to felt
	contractAddressInFelt, err := utils.HexToFelt(contractAddress)
	if err != nil {
		panic(err.Error())
	}

	// Building the functionCall struct, where :
	FnCall := rpc.FunctionCall{
		ContractAddress:    contractAddressInFelt,                         //contractAddress is the contract that we want to call
		EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethod), //this is the function that we want to call
		Calldata:           calldata,
	}

	// Building the Calldata with the help of FmtCalldata where we pass in the FnCall struct along with the Cairo version
	InvokeTx.Calldata, err = accnt.FmtCalldata([]rpc.FunctionCall{FnCall})
	if err != nil {
		panic(err.Error())
	}

	// Signing of the transaction that is done by the account
	err = accnt.SignInvokeTransaction(context.Background(), &InvokeTx)
	if err != nil {
		panic(err.Error())
	}

	// Check the signature is valid or not, this is in StarkCurve!!!
	txHash, err := accnt.TransactionHashInvoke(InvokeTx)
	if err != nil {
		panic(err)
	}
	fmt.Println(txHash)
	checkCalldata := append([]*felt.Felt{
		txHash},
		InvokeTx.Signature...)
	checkTx := rpc.FunctionCall{
		ContractAddress:    accnt.AccountAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt("isValidSignature"),
		Calldata:           checkCalldata,
	}
	fmt.Println(checkCalldata)
	isValid, err := l2Client.Call(context.Background(), checkTx, rpc.BlockID{Tag: "latest"})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(isValid)

	// After the signing we finally call the AddInvokeTransaction in order to invoke the contract function
	resp, err := accnt.AddInvokeTransaction(context.Background(), InvokeTx)
	if err != nil {
		panic(err.Error())
	}
	// This returns us with the transaction hash
	fmt.Println("Transaction hash response : ", resp.TransactionHash)
}

/* ======================
 * L1 Helper Function
 * ====================== */

func GetL1ETHBalance(rpcclient *ethrpc.Client, account common.Address) (*big.Float, error) {
	var result string
	if err := rpcclient.Call(&result, "eth_getBalance", account, "latest"); err != nil {
		return nil, err
	}
	fbalance := new(big.Float)
	fbalance.SetString(result)
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	return ethValue, nil
}

func GetEthereumAddressFromPrivKey(privKey *ecdsa.PrivateKey) common.Address {
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address
}

func GetEthereumAddressFromHexPrivKey(privKeyInHexString string) common.Address {
	privKey, err := crypto.HexToECDSA(privKeyInHexString)
	if err != nil {
		log.Fatal(err)
	}
	return GetEthereumAddressFromPrivKey(privKey)
}

func ConstructTxData(functionSelector string, params []string) (string, []byte) {
	methodID := crypto.Keccak256([]byte(functionSelector + "()"))[:4]
	// fmt.Println(hexutil.Encode(methodID))

	var data []byte
	data = append(data, methodID...)
	// If we need to give input params to contract function, we can append here
	return hexutil.Encode(methodID), data
}

func SignTransaction(rpcUrl string, privKey *ecdsa.PrivateKey, _to string, _value int64, _data string) (*types.Transaction, string) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}
	fromAddress := GetEthereumAddressFromPrivKey(privKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(_value) // in wei (1 eth)
	gasLimit := uint64(21000)   // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress(_to)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)
	if err != nil {
		log.Fatal(err)
	}

	ts := types.Transactions{signedTx}
	rawTxBytes, _ := rlp.EncodeToBytes(ts[0])
	rawTxHex := hex.EncodeToString(rawTxBytes)
	return signedTx, rawTxHex
}

type TxResult struct {
	TBD string
}

type request struct {
	From *common.Address `json:"from"`
	To   *common.Address `json:"to"`
	Data []byte          `json:"data"`
}

func SendTransaction(rpcUrl string, from *common.Address, signedTx *types.Transaction) (string, error) {
	// Setup client
	client, err := ethrpc.DialHTTP(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var txResult TxResult
	req := request{
		From: from,
		To:   signedTx.To(),
		Data: signedTx.Data(),
	}
	if err := client.Call(&txResult, "debug_traceCall", req, "latest"); err != nil {
		log.Fatal(err)
	}
	fmt.Println(txResult)
	return "TBD", nil
}

func WaitForL2L1MsgConsumed(accnt *account.Account, Txnhash *felt.Felt) {
	for {
		resp, err := accnt.GetTransactionStatus(context.Background(), Txnhash)
		if err != nil {
			panic(err)
		}
		if resp.FinalityStatus == rpc.TxnStatus_Accepted_On_L2 {
			break
		}
		time.Sleep(1 * time.Second)
		fmt.Println(resp.ExecutionStatus.MarshalJSON())
	}
}
