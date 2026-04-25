package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

func dockerComposeFile() string {
	home, _ := os.UserHomeDir()
	candidates := []string{
		"docker/docker-compose.yml",
		filepath.Join(home, ".fhish", "rollups", "fhish-1", "docker", "docker-compose.yml"),
		filepath.Join(home, ".fhish", "docker-compose.yml"),
	}
	
	// Also check relative to binary if possible
	_, src, _, _ := runtime.Caller(0)
	candidates = append(candidates, filepath.Join(filepath.Dir(src), "../docker/docker-compose.yml"))
	
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	
	// Final fallback to the most likely VPS path
	vpsPath := filepath.Join(home, ".fhish", "rollups", "fhish-1", "docker", "docker-compose.yml")
	return vpsPath
}

func runCompose(args ...string) error {
	composeFile := dockerComposeFile()
	cmdArgs := append([]string{"compose", "-f", composeFile}, args...)
	cmd := exec.Command("docker", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func DockerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker",
		Short: "Manage the fhish Docker stack",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "up",
			Short: "Start the full fhish stack in Docker",
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Println("🚀 Starting fhish private rollup stack...")
				return runCompose("up", "--build", "-d", "--remove-orphans")
			},
		},
		&cobra.Command{
			Use:   "down",
			Short: "Stop and remove the fhish Docker stack",
			RunE: func(cmd *cobra.Command, args []string) error {
				return runCompose("down", "-v")
			},
		},
		&cobra.Command{
			Use:   "logs",
			Short: "Follow logs",
			RunE: func(cmd *cobra.Command, args []string) error {
				return runCompose("logs", "-f")
			},
		},
		&cobra.Command{
			Use:   "verify",
			Short: "Run FHE smoke test",
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Println("🔐 Running FHE end-to-end smoke test...")
				return runCompose("--profile", "verify", "run", "--rm", "verifier")
			},
		},
	)

	return cmd
}
