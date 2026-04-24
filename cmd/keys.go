package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fhish/fhish-cli/utils"
	"github.com/spf13/cobra"
)

func KeysCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage FHE keys for the fhish rollup",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "generate-fhe",
			Short: "Generate a new FHE keypair (Client, Server, Compact Public Key)",
			RunE: func(cmd *cobra.Command, args []string) error {
				return generateFHEKeys()
			},
		},
	)

	return cmd
}

func generateFHEKeys() error {
	utils.PrintInfo("Generating new FHE keys natively via fhish-wasm...")

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	wasmSrc := filepath.Join(filepath.Dir(cwd), "packages", "fhish-wasm")
	if _, err := os.Stat(wasmSrc); os.IsNotExist(err) {
		return fmt.Errorf("fhish-wasm package not found at %s", wasmSrc)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	outputDir := filepath.Join(home, ".fhish", "keys-shortint")
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create keys directory: %v", err)
	}

	// Create command to run cargo run --bin shortint_keygen <output_dir>
	cargoCmd := exec.Command("cargo", "run", "--release", "--bin", "shortint_keygen", "--", outputDir)
	cargoCmd.Dir = wasmSrc

	utils.PrintInfo(fmt.Sprintf("Executing: cargo run --release --bin shortint_keygen -- %s", outputDir))
	cargoCmd.Stdout = os.Stdout
	cargoCmd.Stderr = os.Stderr

	if err := cargoCmd.Run(); err != nil {
		return fmt.Errorf("failed to generate FHE keys: %v", err)
	}

	utils.PrintSuccess(fmt.Sprintf("FHE keys successfully generated at %s", outputDir))
	return nil
}
