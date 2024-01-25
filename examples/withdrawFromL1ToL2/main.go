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

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"

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
	accountContractVersion int    = 1 //Replace with the cairo version of your account contract
)

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

var Mainnet = ChainInfo{
	name:   "Mainnet",
	rpcUrl: "",
	tokenList: []TokenInfo{
		{
			name:         "Ether",
			symbol:       "ETH",
			decimals:     18,
			l1TokenAddr:  "0x0000000000000000000000000000000000455448",
			l2TokenAddr:  "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
			l1BridgeAddr: "0xae0Ee0A63A2cE6BaeEFFE56e7714FB4EFE48D419",
			l2BridgeAddr: "0x073314940630fd6dcda0d772d4c972c4e0a9946bef9dabf4ef84eda8ef542b82",
		},
		{
			name:         "USD Coin",
			symbol:       "USDC",
			decimals:     6,
			l1TokenAddr:  "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			l2TokenAddr:  "0x053c91253bc9682c04929ca02ed00b3e423f6710d2ee7e0d5ebb06f3ecf368a8",
			l1BridgeAddr: "0xF6080D9fbEEbcd44D89aFfBFd42F098cbFf92816",
			l2BridgeAddr: "0x05cd48fccbfd8aa2773fe22c217e808319ffcc1c5a6a463f7d8fa2da48218196",
		},
	},
}

func main() {
	// Load ENV vriables
	fmt.Println("Starting getTokenBalance example")
	godotenv.Load(fmt.Sprintf(".env.%s", chainName))
	base := os.Getenv("INTEGRATION_BASE")

	// Setup L2 Client
	c, err := ethrpc.DialContext(context.Background(), base)
	if err != nil {
		fmt.Println("Failed to connect to the client, did you specify the url in the .env.mainnet?")
		panic(err)
	}
	l2Client := rpc.NewProvider(c)
	fmt.Println("Established connection with the client")

	// Load L1 Account
	l1PrivKey := os.Getenv("L1_PRIV_KEY")
	l1Address := GetEthereumAddressFromHexPrivKey(l1PrivKey)
	fmt.Println("Load L1 account: ", l1Address)

	// Load L2 Account
	pubKey := os.Getenv("L2_PUB_KEY")
	l2PrivKey := os.Getenv("L2_PRIV_KEY")
	l2Account := ImportAccount(l2Client, pubKey, l2PrivKey, accountContractVersion)
	fmt.Println("Load L2 account: ", l2Account.AccountAddress)

	// Make sure the balance is enough

	// Call the initiate_withdraw on the L2 StarkGate Contract.

	// Trigger withdraw on the L1 StarkGate Contract.

	// Check the L1 recipient ETH balance
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
