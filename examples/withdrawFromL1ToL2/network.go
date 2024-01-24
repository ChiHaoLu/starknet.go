package main

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
