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

func (m *Manager) Start(command string, args ...string) error {
	_ = os.MkdirAll(filepath.Join(m.HomeDir, "run"), 0755)
	_ = os.MkdirAll(filepath.Join(m.HomeDir, "logs"), 0755)

	logFile, err := os.OpenFile(m.LogFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	
	if err := cmd.Start(); err != nil {
		return err
	}

	err = os.WriteFile(m.PIDFile(), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
	return err
}

func (m *Manager) Stop() error {
	data, err := os.ReadFile(m.PIDFile())
	if err != nil {
		return err
	}
	pid, _ := strconv.Atoi(string(data))
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

func (m *Manager) Status() string {
	if _, err := os.Stat(m.PIDFile()); os.IsNotExist(err) {
		return "Stopped"
	}
	data, _ := os.ReadFile(m.PIDFile())
	pid, _ := strconv.Atoi(string(data))
	process, err := os.FindProcess(pid)
	if err != nil {
		return "Stopped"
	}
	// Check if process is actually running
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return "Stopped"
	}
	return "Running (PID: " + string(data) + ")"
}
