package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/fhish/fhish-cli/config"
	"github.com/fhish/fhish-cli/service"
)

func RelayerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relayer",
		Short: "Manage the FHE Relayer",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "start",
			Short: "Start the relayer service",
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, _ := config.GetActiveChain()
				m := service.NewManager("relayer", cfg.Home)
				
				cwd, _ := os.Getwd()
				relayerDir := filepath.Join(filepath.Dir(cwd), "packages", "fhish-relayer-v2")
				
				return m.Start("npm", []string{"--prefix", relayerDir, "start"}, []string{})
			},
		},
		&cobra.Command{
			Use:   "stop",
			Short: "Stop the relayer service",
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, _ := config.GetActiveChain()
				m := service.NewManager("relayer", cfg.Home)
				return m.Stop()
			},
		},
		&cobra.Command{
			Use:   "status",
			Short: "Check relayer status",
			Run: func(cmd *cobra.Command, args []string) {
				cfg, _ := config.GetActiveChain()
				m := service.NewManager("relayer", cfg.Home)
				running, pid, _ := m.Status()
				if running {
					fmt.Printf("Relayer status: running (pid: %d)\n", pid)
				} else {
					fmt.Printf("Relayer status: stopped\n")
				}
			},
		},
	)

	return cmd
}
