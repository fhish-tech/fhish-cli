package minievm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fhish/fhish-cli/utils"
)

func GetLatestMinievmTag() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/initia-labs/minievm/releases/latest")
	if err != nil {
		return "v0.3.0", nil // Fallback
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "v0.3.0", nil
	}
	return release.TagName, nil
}

func Build(targetDir string) error {
	tag, err := GetLatestMinievmTag()
	if err != nil {
		tag = "main" // fallback
	}

	repoURL := "https://github.com/initia-labs/minievm.git"
	srcDir := filepath.Join(targetDir, "src", "minievm")
	
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		utils.PrintInfo(fmt.Sprintf("Cloning MiniEVM (%s)...", tag))
		_ = os.MkdirAll(filepath.Dir(srcDir), 0755)
		err := utils.RunCommand("git", "clone", repoURL, srcDir, "--depth", 1, "--branch", tag)
		if err != nil {
			return err
		}
	}

	utils.PrintInfo("Building MiniEVM...")
	err = utils.RunCommandWithOutput(srcDir, "make", "install")
	if err != nil {
		return fmt.Errorf("failed to build minievm: %w", err)
	}

	utils.PrintSuccess("MiniEVM built and installed.")
	return nil
}
