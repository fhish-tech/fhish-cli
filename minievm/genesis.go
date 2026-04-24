package minievm

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fhish/fhish-cli/config"
)

func PatchGenesis(genesisPath string, cfg *config.ChainConfig) error {
	data, err := os.ReadFile(genesisPath)
	if err != nil {
		return fmt.Errorf("failed to read genesis file: %w", err)
	}

	var genesis map[string]interface{}
	if err := json.Unmarshal(data, &genesis); err != nil {
		return fmt.Errorf("failed to unmarshal genesis: %w", err)
	}

	// 1. Set EVM params
	appState := genesis["app_state"].(map[string]interface{})
	evmState := appState["evm"].(map[string]interface{})
	params := evmState["params"].(map[string]interface{})
	
	// Set evm_chain_id
	params["evm_chain_id"] = fmt.Sprintf("%d", cfg.EVMChainID)

	// 2. Set high gas limits for FHE
	consensusParams := genesis["consensus_params"].(map[string]interface{})
	blockParams := consensusParams["block"].(map[string]interface{})
	blockParams["max_gas"] = "100000000000"

	// 3. whitelisted precompiles (optional, but good for FHE stack)
	// Some MiniEVM versions have a whitelist for precompiles
	// params["whitelisted_precompiles"] = append(params["whitelisted_precompiles"].([]interface{}), "0x000000000000000000000000000000000000005d")

	newData, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal patched genesis: %w", err)
	}

	return os.WriteFile(genesisPath, newData, 0644)
}

func PatchAppToml(appTomlPath string, gasDenom string) error {
	data, err := os.ReadFile(appTomlPath)
	if err != nil {
		return err
	}

	content := string(data)
	// Simple string replacement for key config items
	content = replaceKey(content, "minimum-gas-prices", fmt.Sprintf("\"0%s\"", gasDenom))
	content = replaceKey(content, "enable", "true") // For JSON-RPC
	content = replaceKey(content, "address", "\"0.0.0.0:8545\"")
	
	return os.WriteFile(appTomlPath, []byte(content), 0644)
}

func replaceKey(content string, key string, value string) string {
	// Simple line-based replacement
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), key) {
			lines[i] = fmt.Sprintf("%s = %s", key, value)
		}
	}
	return strings.Join(lines, "\n")
}
