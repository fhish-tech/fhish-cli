package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/fhish/fhish-cli/config"
	"github.com/fhish/fhish-cli/service"
)

func GatewayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "Manage the FHE Gateway",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "start",
			Short: "Start the gateway service",
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, _ := config.GetActiveChain()
				m := service.NewManager("gateway", cfg.Home)
				
				cwd, _ := os.Getwd()
				gatewayDir := filepath.Join(filepath.Dir(cwd), "fhish-gateway")
				
				return m.Start("npm", "--prefix", gatewayDir, "start")
			},
		},
		&cobra.Command{
			Use:   "stop",
			Short: "Stop the gateway service",
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, _ := config.GetActiveChain()
				m := service.NewManager("gateway", cfg.Home)
				return m.Stop()
			},
		},
		&cobra.Command{
			Use:   "status",
			Short: "Check gateway status",
			Run: func(cmd *cobra.Command, args []string) {
				cfg, _ := config.GetActiveChain()
				m := service.NewManager("gateway", cfg.Home)
				fmt.Printf("Gateway status: %s\n", m.Status())
			},
		},
	)

	return cmd
}
