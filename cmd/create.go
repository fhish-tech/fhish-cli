package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fhish/fhish-cli/utils"
)

func CreateCommand() *cobra.Command {
	var createFlag string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new private FHE rollup or service",
		RunE: func(cmd *cobra.Command, args []string) error {
			target := ""
			if len(args) > 0 {
				target = args[0]
			} else if createFlag != "" {
				target = createFlag
			}

			if target == "" {
				return fmt.Errorf("please specify what to create (minievm, gateway, relayer, or all)")
			}

			switch target {
			case "minievm":
				return createMiniEVM()
			case "gateway":
				return createGateway()
			case "relayer":
				return createRelayer()
			case "all":
				return createAll()
			default:
				return fmt.Errorf("unsupported target: %s", target)
			}
		},
	}

	cmd.Flags().StringVar(&createFlag, "create", "", "Target to create")
	
	return cmd
}

func createMiniEVM() error {
	utils.PrintStep(1, 1, "Starting MiniEVM Creation Wizard...")
	// TODO: Call bubbletea model
	return nil
}

func createGateway() error {
	utils.PrintStep(1, 1, "Starting Gateway Setup...")
	return nil
}

func createRelayer() error {
	utils.PrintStep(1, 1, "Starting Relayer Setup...")
	return nil
}

func createAll() error {
	utils.PrintStep(1, 6, "Orchestrating Full FHE Rollup Stack...")
	// 1. Wizard
	// 2. Clone/Build
	// 3. Init/Patch
	// 4. Register
	// 5. Deploy Contracts
	// 6. Start Services
	return nil
}
