package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// State represents the current state of infrastructure
type State struct {
	Version   string                     `json:"version"`
	Resources map[string]*ResourceState  `json:"resources"`
}

// ResourceState represents the state of a single resource
type ResourceState struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Status     string                 `json:"status"` // created, updating, deleting, etc.
	Attributes map[string]interface{} `json:"attributes"`
}

// Manager handles state persistence
type Manager struct {
	stateDir  string
	stateFile string
}

func NewManager(stateDir string) *Manager {
	return &Manager{
		stateDir:  stateDir,
		stateFile: filepath.Join(stateDir, "tblang.tfstate"),
	}
}

// LoadState loads the current state from disk
func (m *Manager) LoadState() (*State, error) {
	if _, err := os.Stat(m.stateFile); os.IsNotExist(err) {
		return &State{
			Version:   "1.0",
			Resources: make(map[string]*ResourceState),
		}, nil
	}
	
	data, err := os.ReadFile(m.stateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}
	
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}
	
	return &state, nil
}

// SaveState saves the current state to disk
func (m *Manager) SaveState(state *State) error {
	// Ensure state directory exists
	if err := os.MkdirAll(m.stateDir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}
	
	state.Version = "1.0"
	
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}
	
	if err := os.WriteFile(m.stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}
	
	return nil
}

// ClearState removes the state file
func (m *Manager) ClearState() error {
	if _, err := os.Stat(m.stateFile); os.IsNotExist(err) {
		return nil // Already cleared
	}
	
	return os.Remove(m.stateFile)
}

// BackupState creates a backup of the current state
func (m *Manager) BackupState() error {
	if _, err := os.Stat(m.stateFile); os.IsNotExist(err) {
		return nil // No state to backup
	}
	
	backupFile := m.stateFile + ".backup"
	
	data, err := os.ReadFile(m.stateFile)
	if err != nil {
		return fmt.Errorf("failed to read state file for backup: %w", err)
	}
	
	if err := os.WriteFile(backupFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}
	
	return nil
}