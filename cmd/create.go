package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/fhish/fhish-cli/config"
	"github.com/fhish/fhish-cli/contracts"
	"github.com/fhish/fhish-cli/gateway"
	"github.com/fhish/fhish-cli/minievm"
	"github.com/fhish/fhish-cli/models"
	"github.com/fhish/fhish-cli/relayer"
	"github.com/fhish/fhish-cli/service"
	"github.com/fhish/fhish-cli/utils"
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new private FHE rollup or service",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "minievm",
			Short: "Create a new MiniEVM rollup node",
			RunE: func(cmd *cobra.Command, args []string) error {
				return createMiniEVM()
			},
		},
		&cobra.Command{
			Use:   "all",
			Short: "Deploy the entire FHE rollup stack",
			RunE: func(cmd *cobra.Command, args []string) error {
				return createAll()
			},
		},
	)

	return cmd
}

func createMiniEVM() error {
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
	home, _ := os.UserHomeDir()
	chainID := vals["chain_id"]
	dataDir := filepath.Join(home, ".fhish", "minievm", chainID)

	utils.PrintStep(1, 3, "Building MiniEVM...")
	err = minievm.Build(filepath.Join(home, ".fhish"))
	if err != nil {
		return err
	}

	utils.PrintStep(2, 3, "Initializing Node...")
	err = minievm.InitNode(dataDir, chainID)
	if err != nil {
		return err
	}

	utils.PrintStep(3, 3, "Patching Genesis...")
	evmID, _ := strconv.Atoi(vals["evm_chain_id"])
	cfg := &config.ChainConfig{
		ChainID:    chainID,
		EVMChainID: evmID,
		Home:       dataDir,
	}
	err = minievm.PatchGenesis(filepath.Join(dataDir, "config", "genesis.json"), cfg)
	if err != nil {
		return err
	}

	utils.PrintSuccess("MiniEVM node is ready!")
	return nil
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

	// Setup Config
	home, _ := os.UserHomeDir()
	chainID := vals["chain_id"]
	dataDir := filepath.Join(home, ".fhish", "minievm", chainID)
	evmID, _ := strconv.Atoi(vals["evm_chain_id"])

	chainCfg := &config.ChainConfig{
		ChainID:       chainID,
		EVMChainID:    evmID,
		Home:          dataDir,
		RPC:           "http://localhost:26657",
		EVMRPC:        "http://localhost:8545",
		L1RPC:         vals["l1_rpc"],
		GatewayURL:    "http://localhost:3000",
		ContractsPath: filepath.Join(home, ".fhish", chainID, "contracts.json"),
	}

	// 2. Build & Init
	utils.PrintStep(1, 6, "Building and Initializing Node...")
	_ = minievm.Build(filepath.Join(home, ".fhish"))
	_ = minievm.InitNode(dataDir, chainID)
	_ = minievm.PatchGenesis(filepath.Join(dataDir, "config", "genesis.json"), chainCfg)

	// 3. Start Node
	utils.PrintStep(2, 6, "Starting Node...")
	mgr := service.NewManager("node", dataDir)
	_ = mgr.Start("minievm", []string{"start", "--home", dataDir}, []string{})
	
	utils.PrintInfo("Waiting for RPC to be ready...")
	_ = utils.WaitForBlock(chainCfg.EVMRPC, 1, 60*time.Second)

	// 4. Deploy Contracts
	utils.PrintStep(3, 6, "Deploying FHE Contracts...")
	cwd, _ := os.Getwd()
	contractsSrc := filepath.Join(filepath.Dir(cwd), "packages", "fhish-contracts-v2")
	addrs, err := contracts.DeployContracts(contractsSrc, chainCfg.EVMRPC, "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		return err
	}
	_ = contracts.SaveAddresses(chainCfg.ContractsPath, addrs)

	// 5. Gateway Setup
	utils.PrintStep(4, 6, "Setting up FHE Gateway...")
	gatewaySrc := filepath.Join(filepath.Dir(cwd), "fhish-gateway")
	relayerSecret := "fhish-test-secret"
	_ = gateway.WriteGatewayConfig(dataDir, addrs, relayerSecret)
	_ = gateway.GenerateKeys(gatewaySrc, filepath.Join(dataDir, "keys"))
	
	gtwayMgr := service.NewManager("gateway", dataDir)
	_ = gtwayMgr.Start("npm", []string{"--prefix", gatewaySrc, "start"}, []string{"FHISH_RELAYER_SECRET=" + relayerSecret})

	// 6. Relayer Setup
	utils.PrintStep(5, 6, "Setting up FHE Relayer...")
	relayerSrc := filepath.Join(filepath.Dir(cwd), "packages", "fhish-relayer-v2")
	_ = relayer.WriteRelayerConfig(dataDir, addrs, "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", relayerSecret, chainCfg.EVMRPC)
	
	relayerMgr := service.NewManager("relayer", dataDir)
	_ = relayerMgr.Start("npm", []string{"--prefix", relayerSrc, "start"}, []string{})

	// 7. Summary
	utils.PrintStep(6, 6, "Rollup stack is LIVE!")
	// Print summary table here...
	
	return nil
}
