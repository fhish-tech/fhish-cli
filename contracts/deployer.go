package contracts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fhish/fhish-cli/utils"
)

type ContractAddresses struct {
	Gateway     string `json:"gateway"`
	PrivateVoting string `json:"private_voting"`
}

func DeployContracts(contractsSrc string, rpcURL string, privateKey string) (*ContractAddresses, error) {
	if !utils.CheckTool("forge") {
		return nil, fmt.Errorf("foundry (forge) not found. Please install it first")
	}

	utils.PrintInfo("Deploying FHE Gateway...")
	// forge create contracts/gateway/FhishGateway.sol:FhishGateway --rpc-url <rpcURL> --private-key <pk> --json
	// This is a placeholder for the actual forge command execution and parsing
	gatewayAddr := "0x5FbDB2315678afecb367f032d93F642f64180aa3" 

	utils.PrintInfo("Deploying PrivateVotingV2...")
	// forge create contracts/PrivateVotingV2.sol:PrivateVotingV2 --rpc-url <rpcURL> --private-key <pk> --json --constructor-args <gatewayAddr> ...
	votingAddr := "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"

	addrs := &ContractAddresses{
		Gateway:       gatewayAddr,
		PrivateVoting: votingAddr,
	}

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
