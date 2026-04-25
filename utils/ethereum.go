package utils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"hash/crc32"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetChainStatus connects to the EVM RPC and returns the current block number
func GetChainStatus(rpcURL string) (uint64, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to Ethereum RPC: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get block number: %v", err)
	}

	return blockNumber, nil
}

// CalculateEVMChainID generates a deterministic EVM chain ID from a string chain ID
func CalculateEVMChainID(chainID string) int64 {
	// For Initia/MiniEVM, it's often a hash of the chain-id string
	// We'll use CRC32 to get a 32-bit integer that fits in EVM chain ID
	h := crc32.NewIEEE()
	h.Write([]byte(chainID))
	return int64(h.Sum32())
}

// GetAddressFromPrivKey derives the EVM address from a hex private key
func GetAddressFromPrivKey(privKeyHex string) (string, error) {
	if strings.HasPrefix(privKeyHex, "0x") {
		privKeyHex = privKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("failed to cast public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, nil
}
