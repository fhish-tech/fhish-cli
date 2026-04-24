package contracts

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/fhish/fhish-cli/utils"
)

type ContractAddresses struct {
	Gateway     string `json:"gateway"`
	PrivateVoting string `json:"private_voting"`
}

func DeployContracts(contractsSrc string, rpcURL string, privateKey string) (*ContractAddresses, error) {
	utils.PrintInfo("Deploying FHE Contracts via Hardhat...")

	// Create command to run hardhat deploy script
	cmd := exec.Command("npx", "hardhat", "run", "scripts/deploy.js", "--network", "sepolia")
	cmd.Dir = contractsSrc

	// Set environment variables for Hardhat
	env := os.Environ()
	env = append(env, "SEPOLIA_RPC_URL="+rpcURL)
	env = append(env, "PRIVATE_KEY="+privateKey)
	cmd.Env = env

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("hardhat deploy failed: %v\nOutput: %s", err, string(output))
	}

	outStr := string(output)
	
	// Parse addresses from output
	addrs := &ContractAddresses{}
	
	// Example output: GATEWAY_ADDRESS=0x...
	gatewayRe := regexp.MustCompile(`GATEWAY_ADDRESS=(0x[a-fA-F0-9]{40})`)
	votingRe := regexp.MustCompile(`VOTING_ADDRESS=(0x[a-fA-F0-9]{40})`)

	if match := gatewayRe.FindStringSubmatch(outStr); len(match) > 1 {
		addrs.Gateway = match[1]
	} else {
		return nil, fmt.Errorf("failed to find GATEWAY_ADDRESS in deploy output")
	}

	if match := votingRe.FindStringSubmatch(outStr); len(match) > 1 {
		addrs.PrivateVoting = match[1]
	} else {
		return nil, fmt.Errorf("failed to find VOTING_ADDRESS in deploy output")
	}

	utils.PrintSuccess(fmt.Sprintf("Deployed Gateway: %s", addrs.Gateway))
	utils.PrintSuccess(fmt.Sprintf("Deployed Voting: %s", addrs.PrivateVoting))

	return addrs, nil
}

func SaveAddresses(path string, addrs *ContractAddresses) error {
	data, err := json.MarshalIndent(addrs, "", "  ")
	if err != nil {
		return err
	}
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	return os.WriteFile(path, data, 0644)
}
