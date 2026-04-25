package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fhish/fhish-cli/config"
)

var (
	cfgFile string
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "fhish",
		Short: "Fhish CLI - Private FHE Rollup Manager for Initia",
		Long: `Fhish is a private fork of the Weave CLI specialized for FHE (Fully Homomorphic Encryption) private rollups on Initia.
It automates the deployment of rollups, FHE contracts, gateways, and relayers.`,
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fhish/config.yaml)")

	rootCmd.AddCommand(
		CreateCommand(),
		DeployCommand(),
		RelayerCommand(),
		GatewayCommand(),
		NodeCommand(),
		DockerCommand(),
		VersionCommand(),
		KeysCommand(),
	)

	return rootCmd.Execute()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fhishDir := filepath.Join(home, ".fhish")
		if _, err := os.Stat(fhishDir); os.IsNotExist(err) {
			_ = os.MkdirAll(fhishDir, 0755)
		}

		viper.AddConfigPath(fhishDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// config loaded
	} else {
		// create default config if not exists
		config.SaveDefaultConfig()
	}
}

func VersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of fhish",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("fhish version v0.1.8")
		},
	}
}
