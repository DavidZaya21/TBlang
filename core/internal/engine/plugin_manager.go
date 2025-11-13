package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/tblang/core/pkg/plugin"
)

// PluginManager handles plugin discovery, lifecycle, and communication
type PluginManager struct {
	pluginDir string
	plugins   map[string]*Plugin
	mu        sync.RWMutex
}

// Plugin represents a loaded provider plugin
type Plugin struct {
	Name       string
	Version    string
	Path       string
	Client     plugin.ProviderPlugin
	Process    *os.Process
	configured bool
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(pluginDir string) *PluginManager {
	return &PluginManager{
		pluginDir: pluginDir,
		plugins:   make(map[string]*Plugin),
	}
}

// DiscoverPlugins discovers available plugins in the plugin directory
func (m *PluginManager) DiscoverPlugins() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, err := os.Stat(m.pluginDir); os.IsNotExist(err) {
		// Create plugin directory if it doesn't exist
		if err := os.MkdirAll(m.pluginDir, 0755); err != nil {
			return fmt.Errorf("failed to create plugin directory: %w", err)
		}
		return nil // No plugins to discover yet
	}

	entries, err := os.ReadDir(m.pluginDir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, "tblang-provider-") {
			continue
		}

		// Extract provider name from filename
		// e.g., tblang-provider-aws -> aws
		providerName := strings.TrimPrefix(name, "tblang-provider-")
		if idx := strings.Index(providerName, "_"); idx != -1 {
			providerName = providerName[:idx]
		}

		pluginPath := filepath.Join(m.pluginDir, name)
		
		// Check if file is executable
		if info, err := os.Stat(pluginPath); err == nil {
			if info.Mode()&0111 != 0 { // Check execute permission
				m.plugins[providerName] = &Plugin{
					Name:    providerName,
					Path:    pluginPath,
					Version: "1.0.0", // TODO: Extract from plugin metadata
				}
			}
		}
	}

	return nil
}

// LoadPlugin loads and starts a specific plugin
func (m *PluginManager) LoadPlugin(ctx context.Context, providerName string) (*Plugin, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pluginInstance, exists := m.plugins[providerName]
	if !exists {
		return nil, fmt.Errorf("plugin not found: %s", providerName)
	}

	if pluginInstance.Client != nil {
		return pluginInstance, nil // Already loaded
	}

	// Start the plugin process
	cmd := exec.CommandContext(ctx, pluginInstance.Path)
	cmd.Env = append(os.Environ(), "TBLANG_PLUGIN_MODE=1")
	
	// Create pipes for communication
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	fmt.Printf("Starting plugin: %s\n", pluginInstance.Path)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start plugin %s: %w", providerName, err)
	}

	pluginInstance.Process = cmd.Process
	fmt.Printf("Plugin process started with PID: %d\n", cmd.Process.Pid)

	// Read connection info from plugin
	var connectionInfo map[string]interface{}
	decoder := json.NewDecoder(stdout)
	
	// Add timeout for reading connection info
	done := make(chan error, 1)
	go func() {
		done <- decoder.Decode(&connectionInfo)
	}()
	
	select {
	case err := <-done:
		if err != nil {
			cmd.Process.Kill()
			return nil, fmt.Errorf("failed to read connection info: %w", err)
		}
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		return nil, fmt.Errorf("timeout waiting for plugin connection info")
	}

	// Connect to plugin
	address, ok := connectionInfo["address"].(string)
	if !ok {
		cmd.Process.Kill()
		return nil, fmt.Errorf("invalid connection info")
	}

	protocol, _ := connectionInfo["protocol"].(string)
	if protocol != "grpc" {
		cmd.Process.Kill()
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}

	// Wait a bit for the plugin to be ready
	time.Sleep(100 * time.Millisecond)

	client, err := plugin.NewGRPCClient(address)
	if err != nil {
		cmd.Process.Kill()
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	pluginInstance.Client = client

	return pluginInstance, nil
}

// GetPlugin returns a loaded plugin
func (m *PluginManager) GetPlugin(providerName string) (*Plugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[providerName]
	if !exists {
		return nil, fmt.Errorf("plugin not found: %s", providerName)
	}

	if plugin.Client == nil {
		return nil, fmt.Errorf("plugin not loaded: %s", providerName)
	}

	return plugin, nil
}

// ConfigurePlugin configures a plugin with the given configuration
func (m *PluginManager) ConfigurePlugin(ctx context.Context, providerName string, config interface{}) error {
	pluginInstance, err := m.GetPlugin(providerName)
	if err != nil {
		return err
	}

	req := &plugin.ConfigureRequest{
		TerraformVersion: "1.0.0", // TBLang version
		Config:           config,
	}

	resp, err := pluginInstance.Client.Configure(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to configure plugin %s: %w", providerName, err)
	}

	if len(resp.Diagnostics) > 0 {
		for _, diag := range resp.Diagnostics {
			if diag.Severity == "error" {
				return fmt.Errorf("plugin configuration error: %s", diag.Summary)
			}
		}
	}

	pluginInstance.configured = true
	return nil
}

// ShutdownAll shuts down all loaded plugins
func (m *PluginManager) ShutdownAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []string
	for name, plugin := range m.plugins {
		if plugin.Process != nil {
			if err := plugin.Process.Kill(); err != nil {
				errors = append(errors, fmt.Sprintf("failed to kill plugin %s: %v", name, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("plugin shutdown errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// ListPlugins returns a list of discovered plugins
func (m *PluginManager) ListPlugins() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var names []string
	for name := range m.plugins {
		names = append(names, name)
	}
	return names
}

