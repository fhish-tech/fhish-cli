package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func WaitForBlock(rpcURL string, targetBlock uint64, timeout time.Duration) error {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for block %d", targetBlock)
		case <-ticker.C:
			blockNumber, err := client.BlockNumber(context.Background())
			if err != nil {
				continue
			}
			if blockNumber >= targetBlock {
				return nil
			}
		}
	}
}
