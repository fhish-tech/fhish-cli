package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func SaveDefaultConfig() {
	home, _ := os.UserHomeDir()
	fhishDir := filepath.Join(home, ".fhish")
	configPath := filepath.Join(fhishDir, "config.yaml")

	if _, err := os.Stat(configPath); err == nil {
		return
	}

	defaultConfig := FhishConfig{
		ActiveChain: "fhish-1",
		Chains: map[string]ChainConfig{
			"fhish-1": {
				ChainID:    "fhish-1",
				EVMChainID: 1234,
				Moniker:    "fhish-node",
				Home:       filepath.Join(fhishDir, "minievm", "fhish-1"),
				L1RPC:      "https://rpc.testnet.initia.xyz",
				L1ChainID:  "initiation-2",
				RPC:        "http://localhost:26657",
				EVMRPC:     "http://localhost:8545",
				EVMWS:      "ws://localhost:8546",
				GatewayURL: "http://localhost:3000",
				ContractsPath: filepath.Join(fhishDir, "fhish-1", "contracts.json"),
				KeysPath:      filepath.Join(fhishDir, "fhish-1", "keys"),
			},
		},
	}

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		fmt.Printf("Error marshaling default config: %v\n", err)
		return
	}

	_ = os.MkdirAll(fhishDir, 0755)
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		fmt.Printf("Error writing default config: %v\n", err)
	}
}

func GetConfig() (*FhishConfig, error) {
	var cfg FhishConfig
	err := viper.Unmarshal(&cfg)
	return &cfg, err
}

func GetActiveChain() (*ChainConfig, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}
	chain, ok := cfg.Chains[cfg.ActiveChain]
	if !ok {
		return nil, fmt.Errorf("active chain %s not found in config", cfg.ActiveChain)
	}
	return &chain, nil
}
