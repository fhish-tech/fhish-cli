package config

type FhishConfig struct {
	ActiveChain string                 `yaml:"active_chain"`
	Chains      map[string]ChainConfig `yaml:"chains"`
}

type ChainConfig struct {
	ChainID       string `yaml:"chain_id"`
	EVMChainID    int    `yaml:"evm_chain_id"`
	Moniker       string `yaml:"moniker"`
	Home          string `yaml:"home"`
	L1RPC         string `yaml:"l1_rpc"`
	L1ChainID     string `yaml:"l1_chain_id"`
	RPC           string `yaml:"rpc"`
	EVMRPC        string `yaml:"evm_rpc"`
	EVMWS         string `yaml:"evm_ws"`
	GatewayURL    string `yaml:"gateway_url"`
	ContractsPath string `yaml:"contracts_path"`
	KeysPath      string `yaml:"keys_path"`
}
