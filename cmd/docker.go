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
	candidates := []string{
		"docker/docker-compose.yml",
		filepath.Join(os.Getenv("HOME"), ".fhish", "docker-compose.yml"),
	}
	_, src, _, _ := runtime.Caller(0)
	candidates = append(candidates, filepath.Join(filepath.Dir(src), "../docker/docker-compose.yml"))
	
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return "docker/docker-compose.yml"
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
