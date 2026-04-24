package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

type Manager struct {
	Name    string
	HomeDir string
}

func NewManager(name string, homeDir string) *Manager {
	return &Manager{
		Name:    name,
		HomeDir: homeDir,
	}
}

func (m *Manager) PIDFile() string {
	return filepath.Join(m.HomeDir, "run", m.Name+".pid")
}

func (m *Manager) LogFile() string {
	return filepath.Join(m.HomeDir, "logs", m.Name+".log")
}

func (m *Manager) Start(command string, args []string, env []string) error {
	_ = os.MkdirAll(filepath.Join(m.HomeDir, "run"), 0755)
	_ = os.MkdirAll(filepath.Join(m.HomeDir, "logs"), 0755)

	if running, _, _ := m.Status(); running {
		return fmt.Errorf("%s is already running", m.Name)
	}

	logFile, err := os.OpenFile(m.LogFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Env = append(os.Environ(), env...)
	
	if err := cmd.Start(); err != nil {
		return err
	}

	err = os.WriteFile(m.PIDFile(), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
	return err
}

func (m *Manager) Stop() error {
	running, pid, err := m.Status()
	if !running {
		return fmt.Errorf("%s is not running", m.Name)
	}
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = process.Signal(syscall.SIGTERM)
	if err == nil {
		_ = os.Remove(m.PIDFile())
	}
	return err
}

func (m *Manager) Status() (bool, int, error) {
	data, err := os.ReadFile(m.PIDFile())
	if err != nil {
		return false, 0, nil
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return false, 0, err
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, 0, nil
	}
	// Check if process is actually running
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false, 0, nil
	}
	return true, pid, nil
}

func (m *Manager) Log(follow bool, lines int) error {
	args := []string{"-n", strconv.Itoa(lines)}
	if follow {
		args = append(args, "-f")
	}
	args = append(args, m.LogFile())
	
	cmd := exec.Command("tail", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
