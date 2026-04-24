package minievm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fhish/fhish-cli/utils"
)

func Build(targetDir string) error {
	repoURL := "https://github.com/initia-labs/minievm.git"
	srcDir := filepath.Join(targetDir, "src", "minievm")
	
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		utils.PrintInfo("Cloning MiniEVM...")
		_ = os.MkdirAll(filepath.Dir(srcDir), 0755)
		err := utils.RunCommand("git", "clone", repoURL, srcDir, "--depth", 1)
		if err != nil {
			return err
		}
	} else {
		utils.PrintInfo("MiniEVM source already exists, skipping clone.")
	}

	utils.PrintInfo("Building MiniEVM...")
	err := utils.RunCommandWithOutput(srcDir, "make", "install")
	if err != nil {
		return fmt.Errorf("failed to build minievm: %w", err)
	}

	utils.PrintSuccess("MiniEVM built successfully.")
	return nil
}

func InitNode(homeDir string, chainID string) error {
	utils.PrintInfo("Initializing MiniEVM node...")
	_ = os.MkdirAll(homeDir, 0755)
	err := utils.RunCommand("minievm", "init", "fhish-node", "--chain-id", chainID, "--home", homeDir)
	if err != nil {
		return err
	}
	utils.PrintSuccess("Node initialized.")
	return nil
}
