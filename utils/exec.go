package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunCommandWithOutput(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	
	if err := cmd.Start(); err != nil {
		return err
	}
	
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Printf("   [OUT] %s\n", scanner.Text())
		}
	}()
	
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Printf("   [ERR] %s\n", scanner.Text())
		}
	}()
	
	return cmd.Wait()
}

func CheckTool(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
