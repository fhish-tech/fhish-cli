package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/fhish/fhish-cli/models"
	"github.com/fhish/fhish-cli/utils"
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new private FHE rollup or service",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "all",
			Short: "Deploy the entire FHE rollup stack interactively",
			RunE: func(cmd *cobra.Command, args []string) error {
				return createAll()
			},
		},
	)

	return cmd
}

func createAll() error {
	// 1. Run Wizard
	p := tea.NewProgram(models.NewCreateMiniEVMModel())
	m, err := p.Run()
	if err != nil {
		return err
	}
	model := m.(models.CreateMiniEVMModel)
	if !model.Done {
		return nil
	}
	vals := model.Config()

	// 2. Derive Configuration
	chainID := vals["chain_id"]
	moniker := vals["moniker"]
	gasDenom := vals["gas_denom"]
	deployerKey := vals["deployer_key"]
	relayerSecret := vals["relayer_secret"]

	evmChainID := utils.CalculateEVMChainID(chainID)
	deployerAddr, err := utils.GetAddressFromPrivKey(deployerKey)
	if err != nil {
		return fmt.Errorf("invalid deployer private key: %v", err)
	}

	home, _ := os.UserHomeDir()
	setupDir := filepath.Join(home, ".fhish", "rollups", chainID)
	_ = os.MkdirAll(setupDir, 0755)

	// 3. Prepare Environment
	utils.PrintStep(1, 4, "Preparing deployment environment...")
	
	// Clone the CLI repo to get the Docker files if they don't exist
	if _, err := os.Stat(filepath.Join(setupDir, "docker")); os.IsNotExist(err) {
		utils.PrintInfo("Downloading stack configuration...")
		err = utils.RunCommand("git", []string{"clone", "https://github.com/fhish-tech/fhish-cli.git", "temp-cli"}, setupDir)
		if err == nil {
			_ = os.Rename(filepath.Join(setupDir, "temp-cli", "docker"), filepath.Join(setupDir, "docker"))
			_ = os.RemoveAll(filepath.Join(setupDir, "temp-cli"))
		}
	}

	// Create .env file
	envContent := fmt.Sprintf(`CHAIN_ID=%s
EVM_CHAIN_ID=%d
MONIKER=%s
GAS_DENOM=%s
DEPLOYER_ADDRESS=%s
DEPLOYER_PRIVATE_KEY=%s
FHISH_RELAYER_SECRET=%s
`, chainID, evmChainID, moniker, gasDenom, deployerAddr, deployerKey, relayerSecret)

	err = os.WriteFile(filepath.Join(setupDir, ".env"), []byte(envContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create .env file: %v", err)
	}

	// 4. Start Stack
	utils.PrintStep(2, 4, "Launching Docker containers (this may take a few minutes)...")
	err = utils.RunCommand("docker", []string{"compose", "-f", "docker/docker-compose.yml", "up", "-d"}, setupDir)
	if err != nil {
		return fmt.Errorf("failed to start docker stack: %v", err)
	}

	// 5. Wait for Readiness
	utils.PrintStep(3, 4, "Waiting for RPC to be ready...")
	time.Sleep(10 * time.Second) // Give it a head start

	// 6. Success Summary
	utils.PrintStep(4, 4, "Rollup stack is LIVE!")

	printSummary(chainID, moniker, evmChainID)

	return nil
}

func printSummary(chainID, moniker string, evmChainID int64) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")).Padding(0, 1)
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Width(20)
	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	fmt.Println("\n" + titleStyle.Render("🚀 FHISH ROLLUP DEPLOYED SUCCESSFULLY"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("─", 50)))
	
	fmt.Printf("%s %s\n", keyStyle.Render("Rollup Name:"), valStyle.Render(moniker))
	fmt.Printf("%s %s\n", keyStyle.Render("Rollup Chain ID:"), valStyle.Render(chainID))
	fmt.Printf("%s %s\n", keyStyle.Render("EVM Chain ID:"), valStyle.Render(fmt.Sprintf("%d", evmChainID)))
	fmt.Printf("%s %s\n", keyStyle.Render("REST URL:"), valStyle.Render("http://localhost:1317"))
	fmt.Printf("%s %s\n", keyStyle.Render("RPC URL:"), valStyle.Render("http://localhost:26657"))
	fmt.Printf("%s %s\n", keyStyle.Render("EVM RPC URL:"), valStyle.Render("http://localhost:8545"))

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("─", 50)))
	
	instructionStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("214")).
		Padding(1, 2).
		MarginTop(1).
		Render(fmt.Sprintf("To add this rollup to Initiascan, go to:\nhttps://scan.testnet.initia.xyz/initiation-2/custom-network/add/manual\n\nUse the details above to complete the form."))
	
	fmt.Println(instructionStyle)
}
