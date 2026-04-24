package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fhish/fhish-cli/config"
	"github.com/fhish/fhish-cli/service"
	"github.com/fhish/fhish-cli/utils"
)

func NodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage the rollup node",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "start",
			Short: "Start the node service",
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, _ := config.GetActiveChain()
				m := service.NewManager("node", cfg.Home)
				return m.Start("minievm", "start", "--home", cfg.Home)
			},
		},
		&cobra.Command{
			Use:   "stop",
			Short: "Stop the node service",
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, _ := config.GetActiveChain()
				m := service.NewManager("node", cfg.Home)
				return m.Stop()
			},
		},
		&cobra.Command{
			Use:   "status",
			Short: "Check node status",
			Run: func(cmd *cobra.Command, args []string) {
				cfg, _ := config.GetActiveChain()
				
				// 1. Check local OS process status
				m := service.NewManager("node", cfg.Home)
				procStatus := m.Status()
				fmt.Printf("Process status: %s\n", procStatus)

				// 2. Check actual RPC connectivity and chain state
				if cfg.EVMRPC != "" {
					blockNum, err := utils.GetChainStatus(cfg.EVMRPC)
					if err != nil {
						fmt.Printf("RPC Status: Offline (%v)\n", err)
					} else {
						fmt.Printf("RPC Status: Online\n")
						fmt.Printf("Current Block: %d\n", blockNum)
					}
				}
			},
		},
	)

	return cmd
}
