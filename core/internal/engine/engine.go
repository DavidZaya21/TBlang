package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/tblang/core/internal/ast"
	"github.com/tblang/core/internal/compiler"
	"github.com/tblang/core/internal/state"
	"github.com/tblang/core/pkg/plugin"
)

// Engine is the main TBLang core engine
type Engine struct {
	compiler      *compiler.Compiler
	stateManager  *state.Manager
	pluginManager *PluginManager
	workingDir    string
}

var (
	// Color definitions
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	warningColor = color.New(color.FgYellow, color.Bold)
	infoColor    = color.New(color.FgCyan, color.Bold)
	headerColor  = color.New(color.FgMagenta, color.Bold)
	createColor  = color.New(color.FgGreen)
	updateColor  = color.New(color.FgYellow)
	deleteColor  = color.New(color.FgRed)
)

// NewEngine creates a new TBLang core engine
func New() *Engine {
	workingDir, _ := os.Getwd()
	
	// Try global plugin directory first, then local
	pluginDir := "/usr/local/lib/tblang/plugins"
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		pluginDir = filepath.Join(workingDir, ".tblang", "plugins")
	}
	
	return &Engine{
		compiler:      compiler.New(),
		stateManager:  state.NewManager(filepath.Join(workingDir, ".tblang")),
		pluginManager: NewPluginManager(pluginDir),
		workingDir:    workingDir,
	}
}

// Initialize initializes the engine and discovers plugins
func (e *Engine) Initialize(ctx context.Context) error {
	// Discover available plugins
	if err := e.pluginManager.DiscoverPlugins(); err != nil {
		return fmt.Errorf("failed to discover plugins: %w", err)
	}

	fmt.Printf("Discovered plugins: %v\n", e.pluginManager.ListPlugins())
	return nil
}

// Plan shows what changes will be made
func (e *Engine) Plan(ctx context.Context, filename string) error {
	fmt.Println("Planning infrastructure changes...")
	
	// Compile the tblang file
	program, err := e.compiler.CompileFile(filename)
	if err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	// Load and configure required plugins
	if err := e.loadRequiredPlugins(ctx, program); err != nil {
		return fmt.Errorf("failed to load plugins: %w", err)
	}
	
	// Load current state
	currentState, err := e.stateManager.LoadState()
	if err != nil {
		fmt.Println("No existing state found, will create new infrastructure")
		currentState = &state.State{Resources: make(map[string]*state.ResourceState)}
	}
	
	// Compare desired vs current state
	changes := e.calculateChanges(program, currentState)
	
	// Display plan
	e.displayPlan(changes)
	
	return nil
}

// Apply creates/updates the infrastructure
func (e *Engine) Apply(ctx context.Context, filename string) error {
	infoColor.Println("Applying infrastructure changes...")
	
	// Compile the tblang file
	program, err := e.compiler.CompileFile(filename)
	if err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	// Load and configure required plugins
	if err := e.loadAndConfigurePlugins(ctx, program); err != nil {
		return fmt.Errorf("failed to load plugins: %w", err)
	}
	
	// Load current state
	currentState, err := e.stateManager.LoadState()
	if err != nil {
		currentState = &state.State{Resources: make(map[string]*state.ResourceState)}
	}
	
	// Compare desired vs current state
	changes := e.calculateChanges(program, currentState)
	
	// Display plan
	e.displayPlan(changes)
	
	// Ask for confirmation
	fmt.Print("\nDo you want to perform these actions? (yes/no): ")
	var response string
	fmt.Scanln(&response)
	
	if response != "yes" && response != "y" {
		warningColor.Println("Apply cancelled.")
		return nil
	}
	
	// Apply changes using plugins
	if err := e.applyChanges(ctx, changes, currentState); err != nil {
		return fmt.Errorf("apply failed: %w", err)
	}
	
	successColor.Println("\nApply complete!")
	return nil
}

// Destroy removes all infrastructure
func (e *Engine) Destroy(ctx context.Context, filename string) error {
	fmt.Println("Destroying infrastructure...")
	
	// Compile the tblang file to get cloud vendor configuration (including profile)
	program, err := e.compiler.CompileFile(filename)
	if err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	// Load and configure required plugins (sets AWS profile and loads plugin)
	if err := e.loadAndConfigurePlugins(ctx, program); err != nil {
		return fmt.Errorf("failed to load plugins: %w", err)
	}
	
	// Load current state
	currentState, err := e.stateManager.LoadState()
	if err != nil {
		fmt.Println("No state found, nothing to destroy")
		return nil
	}

	// Show what will be destroyed
	fmt.Println("\nThe following resources will be destroyed:")
	for name, resource := range currentState.Resources {
		fmt.Printf("  - %s (%s)\n", name, resource.Type)
	}
	
	// Ask for confirmation
	fmt.Print("\nDo you really want to destroy all resources? (yes/no): ")
	var response string
	fmt.Scanln(&response)
	
	if response != "yes" && response != "y" {
		fmt.Println("Destroy cancelled.")
		return nil
	}

	// Destroy resources using plugins
	if err := e.destroyResources(ctx, currentState); err != nil {
		return fmt.Errorf("failed to destroy resources: %w", err)
	}
	
	fmt.Println("Destroy complete!")
	return nil
}

// Show displays current infrastructure state
func (e *Engine) Show() error {
	fmt.Println("Current infrastructure state:")
	
	currentState, err := e.stateManager.LoadState()
	if err != nil {
		fmt.Println("No state found")
		return nil
	}
	
	if len(currentState.Resources) == 0 {
		fmt.Println("No resources found")
		return nil
	}
	
	for name, resource := range currentState.Resources {
		fmt.Printf("\nResource: %s\n", name)
		fmt.Printf("   Type: %s\n", resource.Type)
		fmt.Printf("   Status: %s\n", resource.Status)
		if len(resource.Attributes) > 0 {
			fmt.Println("   Attributes:")
			for key, value := range resource.Attributes {
				fmt.Printf("     %s: %v\n", key, value)
			}
		}
	}
	
	return nil
}

// Graph shows the dependency graph
func (e *Engine) Graph(ctx context.Context, filename string) error {
	program, err := e.compiler.CompileFile(filename)
	if err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	e.displayVisualGraph(program)
	return nil
}

func (e *Engine) displayVisualGraph(program *compiler.Program) {
	headerColor.Println("\nDependency Graph & Deployment Order:")
	headerColor.Println(strings.Repeat("=", 50))
	
	// Create a simplified dependency map
	dependencies := make(map[string][]string)
	
	// Analyze dependencies from resource properties
	for _, resource := range program.Resources {
		var deps []string
		for _, value := range resource.Properties {
			if refs := e.findResourceReferences(value, program.Resources); len(refs) > 0 {
				for _, ref := range refs {
					if ref != resource.Name { // Avoid self-references
						deps = append(deps, ref)
					}
				}
			}
		}
		dependencies[resource.Name] = e.removeDuplicates(deps)
	}
	
	// Display visual graph
	fmt.Println()
	for i, resource := range program.Resources {
		// Resource node with color
		resourceColor := e.getResourceColor(resource.Type)
		resourceColor.Printf("[%d] %s", i+1, resource.Name)
		fmt.Printf(" (%s)\n", resource.Type)
		
		// Show dependencies
		if deps := dependencies[resource.Name]; len(deps) > 0 {
			infoColor.Print("    Dependencies: ")
			for j, dep := range deps {
				if j > 0 {
					fmt.Print(", ")
				}
				warningColor.Print(dep)
			}
			fmt.Println()
		} else {
			infoColor.Println("    No dependencies")
		}
		
		if i < len(program.Resources)-1 {
			fmt.Println("    |")
			fmt.Println("    v")
		}
	}
	
	// Display deployment flow
	fmt.Println()
	headerColor.Println("Deployment Flow:")
	headerColor.Println(strings.Repeat("-", 30))
	
	for i, resource := range program.Resources {
		resourceColor := e.getResourceColor(resource.Type)
		if i > 0 {
			infoColor.Print(" --> ")
		}
		resourceColor.Print(resource.Name)
	}
	fmt.Println("\n")
}

func (e *Engine) removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}

func (e *Engine) getResourceColor(resourceType string) *color.Color {
	switch resourceType {
	case "vpc":
		return color.New(color.FgBlue, color.Bold)
	case "subnet":
		return color.New(color.FgGreen, color.Bold)
	case "security_group":
		return color.New(color.FgYellow, color.Bold)
	case "ec2":
		return color.New(color.FgRed, color.Bold)
	default:
		return color.New(color.FgWhite, color.Bold)
	}
}

func (e *Engine) findResourceReferences(value interface{}, resources []*ast.Resource) []string {
	var refs []string
	resourceNames := make(map[string]bool)
	
	// Build resource name map
	for _, res := range resources {
		resourceNames[res.Name] = true
	}
	
	switch v := value.(type) {
	case string:
		// Only add if it's a different resource (not self-reference)
		if resourceNames[v] {
			refs = append(refs, v)
		}
	case map[string]interface{}:
		for _, val := range v {
			refs = append(refs, e.findResourceReferences(val, resources)...)
		}
	case []interface{}:
		for _, val := range v {
			refs = append(refs, e.findResourceReferences(val, resources)...)
		}
	}
	
	return refs
}

// Shutdown gracefully shuts down the engine
func (e *Engine) Shutdown() error {
	return e.pluginManager.ShutdownAll()
}

// ListPlugins returns available plugins
func (e *Engine) ListPlugins() []string {
	return e.pluginManager.ListPlugins()
}

// Helper methods

func (e *Engine) loadRequiredPlugins(ctx context.Context, program *compiler.Program) error {
	// Extract required providers from cloud_vendor blocks
	for providerName, config := range program.CloudVendors {
		infoColor.Printf("Found provider: %s\n", providerName)
		fmt.Printf("  Region: %v\n", config.Properties["region"])
		
		// Set AWS profile if specified in configuration
		if profile, exists := config.Properties["profile"]; exists {
			if profileStr, ok := profile.(string); ok {
				os.Setenv("AWS_PROFILE", profileStr)
				infoColor.Printf("  Profile: %s\n", profileStr)
			}
		}
		
		if accountID, exists := config.Properties["account_id"]; exists {
			fmt.Printf("  Account ID: %v\n", accountID)
		}
		
		// For now, skip plugin loading and use mock mode for testing
		successColor.Printf("Provider %s configured (mock mode)\n", providerName)
	}

	return nil
}

func (e *Engine) loadAndConfigurePlugins(ctx context.Context, program *compiler.Program) error {
	// Extract required providers from cloud_vendor blocks
	for providerName, config := range program.CloudVendors {
		infoColor.Printf("Found provider: %s\n", providerName)
		fmt.Printf("  Region: %v\n", config.Properties["region"])
		
		// Set AWS profile if specified in configuration
		if profile, exists := config.Properties["profile"]; exists {
			if profileStr, ok := profile.(string); ok {
				os.Setenv("AWS_PROFILE", profileStr)
				infoColor.Printf("  Profile: %s\n", profileStr)
			}
		}
		
		if accountID, exists := config.Properties["account_id"]; exists {
			fmt.Printf("  Account ID: %v\n", accountID)
		}
		
		// Load the plugin
		_, err := e.pluginManager.LoadPlugin(ctx, providerName)
		if err != nil {
			return fmt.Errorf("failed to load plugin %s: %w", providerName, err)
		}
		
		// Configure the plugin
		if err := e.pluginManager.ConfigurePlugin(ctx, providerName, config.Properties); err != nil {
			return fmt.Errorf("failed to configure plugin %s: %w", providerName, err)
		}
		
		successColor.Printf("Provider %s loaded and configured\n", providerName)
	}

	return nil
}

func (e *Engine) calculateChanges(program *compiler.Program, currentState *state.State) *PlanChanges {
	changes := &PlanChanges{
		Create: make([]*state.ResourceState, 0),
		Update: make([]*state.ResourceState, 0),
		Delete: make([]*state.ResourceState, 0),
	}
	
	// Check for new resources to create
	for _, resource := range program.Resources {
		if _, exists := currentState.Resources[resource.Name]; !exists {
			changes.Create = append(changes.Create, &state.ResourceState{
				Name:       resource.Name,
				Type:       resource.Type,
				Status:     "planned",
				Attributes: resource.Properties,
			})
		}
	}
	
	// Check for resources to delete (exist in state but not in config)
	programResources := make(map[string]bool)
	for _, resource := range program.Resources {
		programResources[resource.Name] = true
	}

	for name, resource := range currentState.Resources {
		if !programResources[name] {
			changes.Delete = append(changes.Delete, resource)
		}
	}
	
	return changes
}

func (e *Engine) applyChanges(ctx context.Context, changes *PlanChanges, currentState *state.State) error {
	// Apply creates
	for _, resource := range changes.Create {
		resourceColor := e.getResourceColor(resource.Type)
		resourceColor.Printf("\nCreating %s (%s)...\n", resource.Name, resource.Type)
		
		// Use plugin to create resource
		newState, err := e.createResourceWithPlugin(ctx, resource)
		if err != nil {
			errorColor.Printf("  ✗ Failed to create %s: %v\n", resource.Name, err)
			return fmt.Errorf("failed to create %s: %w", resource.Name, err)
		}
		
		// Update resource with new state from plugin
		if newState != nil {
			if stateMap, ok := newState.(map[string]interface{}); ok {
				resource.Attributes = stateMap
			}
		}
		
		// Update state
		resource.Status = "created"
		currentState.Resources[resource.Name] = resource
		if err := e.stateManager.SaveState(currentState); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
		
		successColor.Printf("  ✓ Created %s (%s)\n", resource.Name, resource.Type)
	}

	// Apply deletes
	for _, resource := range changes.Delete {
		warningColor.Printf("\nDeleting %s (%s)...\n", resource.Name, resource.Type)
		
		// Use plugin to delete resource
		if err := e.destroyResourceWithPlugin(ctx, resource); err != nil {
			errorColor.Printf("  ✗ Failed to delete %s: %v\n", resource.Name, err)
		} else {
			successColor.Printf("  ✓ Deleted %s (%s)\n", resource.Name, resource.Type)
		}
		
		// Remove from state
		delete(currentState.Resources, resource.Name)
		if err := e.stateManager.SaveState(currentState); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
	}

	return nil
}

func (e *Engine) createResourceWithPlugin(ctx context.Context, resource *state.ResourceState) (interface{}, error) {
	// Get the AWS plugin
	pluginInstance, err := e.pluginManager.GetPlugin("aws")
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS plugin: %w", err)
	}

	// Resolve resource references in attributes
	resolvedAttrs := e.resolveResourceReferences(resource.Attributes)

	// Call ApplyResourceChange to create the resource
	req := &plugin.ApplyResourceChangeRequest{
		TypeName:     resource.Type,
		PriorState:   nil, // nil for new resources
		PlannedState: resolvedAttrs,
		Config:       resolvedAttrs,
	}

	resp, err := pluginInstance.Client.ApplyResourceChange(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("plugin error: %w", err)
	}

	// Check for diagnostics
	if len(resp.Diagnostics) > 0 {
		for _, diag := range resp.Diagnostics {
			if diag.Severity == "error" {
				return nil, fmt.Errorf("%s: %s", diag.Summary, diag.Detail)
			} else if diag.Severity == "warning" {
				warningColor.Printf("  ⚠ Warning: %s\n", diag.Summary)
			}
		}
	}

	return resp.NewState, nil
}

func (e *Engine) resolveResourceReferences(attrs map[string]interface{}) map[string]interface{} {
	resolved := make(map[string]interface{})
	
	// Load current state to resolve references
	currentState, err := e.stateManager.LoadState()
	if err != nil {
		return attrs // Return original if can't load state
	}
	
	for key, value := range attrs {
		// Check if value is a string that might be a resource reference
		if strValue, ok := value.(string); ok {
			// Check if this string matches a resource name in state
			if resource, exists := currentState.Resources[strValue]; exists {
				// Try to resolve to the actual resource ID
				switch key {
				case "vpc_id":
					if vpcID, ok := resource.Attributes["vpc_id"].(string); ok {
						resolved[key] = vpcID
						continue
					}
				case "subnet_id":
					if subnetID, ok := resource.Attributes["subnet_id"].(string); ok {
						resolved[key] = subnetID
						continue
					}
				case "group_id":
					if groupID, ok := resource.Attributes["group_id"].(string); ok {
						resolved[key] = groupID
						continue
					}
				}
			}
		}
		// Keep original value if not resolved
		resolved[key] = value
	}
	
	return resolved
}

func (e *Engine) destroyResources(ctx context.Context, currentState *state.State) error {
	// Sort resources by type to ensure proper destruction order
	// Order: security_group -> subnet -> vpc
	var securityGroups []*state.ResourceState
	var subnets []*state.ResourceState
	var vpcs []*state.ResourceState
	var others []*state.ResourceState
	
	for _, resource := range currentState.Resources {
		switch resource.Type {
		case "security_group":
			securityGroups = append(securityGroups, resource)
		case "subnet":
			subnets = append(subnets, resource)
		case "vpc":
			vpcs = append(vpcs, resource)
		default:
			others = append(others, resource)
		}
	}
	
	// Destroy in proper order: security groups first, then subnets, then VPCs
	orderedResources := append([]*state.ResourceState{}, securityGroups...)
	orderedResources = append(orderedResources, subnets...)
	orderedResources = append(orderedResources, others...)
	orderedResources = append(orderedResources, vpcs...)
	
	// Destroy each resource
	for _, resource := range orderedResources {
		warningColor.Printf("Destroying %s (%s)...\n", resource.Name, resource.Type)
		
		// Use plugin to destroy resource
		if err := e.destroyResourceWithPlugin(ctx, resource); err != nil {
			errorColor.Printf("  Error: failed to destroy %s: %v\n", resource.Name, err)
			// Continue with other resources even if one fails
		} else {
			successColor.Printf("Deleted %s (%s)\n", resource.Name, resource.Type)
		}
		
		// Remove from state
		delete(currentState.Resources, resource.Name)
		if err := e.stateManager.SaveState(currentState); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
	}

	return nil
}

func (e *Engine) destroyResourceWithPlugin(ctx context.Context, resource *state.ResourceState) error {
	// Get the AWS plugin
	pluginInstance, err := e.pluginManager.GetPlugin("aws")
	if err != nil {
		return fmt.Errorf("failed to get AWS plugin: %w", err)
	}

	// Import the plugin package types
	req := &plugin.ApplyResourceChangeRequest{
		TypeName:     resource.Type,
		PriorState:   resource.Attributes,
		PlannedState: nil, // nil indicates destroy
		Config:       resource.Attributes,
	}

	resp, err := pluginInstance.Client.ApplyResourceChange(ctx, req)
	if err != nil {
		return fmt.Errorf("plugin error: %w", err)
	}

	// Check for diagnostics
	if len(resp.Diagnostics) > 0 {
		for _, diag := range resp.Diagnostics {
			if diag.Severity == "error" {
				return fmt.Errorf("%s: %s", diag.Summary, diag.Detail)
			}
		}
	}

	return nil
}

func (e *Engine) displayPlan(changes *PlanChanges) {
	headerColor.Println("\nPlan Summary:")
	
	if len(changes.Create) > 0 {
		createColor.Printf("\nResources to create (%d):\n", len(changes.Create))
		for _, resource := range changes.Create {
			createColor.Printf("  + %s ", resource.Name)
			fmt.Printf("(%s)\n", resource.Type)
		}
	}
	
	if len(changes.Update) > 0 {
		updateColor.Printf("\nResources to update (%d):\n", len(changes.Update))
		for _, resource := range changes.Update {
			updateColor.Printf("  ~ %s ", resource.Name)
			fmt.Printf("(%s)\n", resource.Type)
		}
	}
	
	if len(changes.Delete) > 0 {
		deleteColor.Printf("\nResources to delete (%d):\n", len(changes.Delete))
		for _, resource := range changes.Delete {
			deleteColor.Printf("  - %s ", resource.Name)
			fmt.Printf("(%s)\n", resource.Type)
		}
	}
	
	if len(changes.Create) == 0 && len(changes.Update) == 0 && len(changes.Delete) == 0 {
		infoColor.Println("\nNo changes. Infrastructure is up-to-date.")
	}
}

type PlanChanges struct {
	Create []*state.ResourceState
	Update []*state.ResourceState
	Delete []*state.ResourceState
}

// AWS CLI integration methods for testing

func (e *Engine) createVPCWithAWSCLI(resource *state.ResourceState) error {
	cidrBlock, ok := resource.Attributes["cidr_block"].(string)
	if !ok {
		return fmt.Errorf("cidr_block not found in VPC configuration")
	}
	
	infoColor.Printf("  Creating VPC with CIDR: %s\n", cidrBlock)
	
	// Create VPC using AWS CLI
	cmd := exec.Command("aws", "ec2", "create-vpc", "--cidr-block", cidrBlock, "--output", "json")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("aws cli error: %w", err)
	}
	
	// Parse the output to get VPC ID
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return fmt.Errorf("failed to parse AWS CLI output: %w", err)
	}
	
	vpc, ok := result["Vpc"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid AWS CLI response format")
	}
	
	vpcID, ok := vpc["VpcId"].(string)
	if !ok {
		return fmt.Errorf("VPC ID not found in response")
	}
	
	// Update resource attributes with actual VPC ID
	resource.Attributes["vpc_id"] = vpcID
	resource.Attributes["state"] = vpc["State"]
	
	successColor.Printf("  VPC created with ID: %s\n", vpcID)
	
	// Add tags if specified
	if tags, exists := resource.Attributes["tags"]; exists {
		if err := e.tagVPCWithAWSCLI(vpcID, tags); err != nil {
			fmt.Printf("  Warning: failed to tag VPC: %v\n", err)
		}
	}
	
	return nil
}

func (e *Engine) tagVPCWithAWSCLI(vpcID string, tags interface{}) error {
	tagsMap, ok := tags.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid tags format")
	}
	
	for key, value := range tagsMap {
		tagSpec := fmt.Sprintf("Key=%s,Value=%s", key, value)
		cmd := exec.Command("aws", "ec2", "create-tags", "--resources", vpcID, "--tags", tagSpec)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to tag VPC with %s=%s: %w", key, value, err)
		}
	}
	
	successColor.Printf("  VPC tagged successfully\n")
	return nil
}

func (e *Engine) deleteVPCWithAWSCLI(resource *state.ResourceState) error {
	vpcID, ok := resource.Attributes["vpc_id"].(string)
	if !ok {
		return fmt.Errorf("vpc_id not found in resource state")
	}
	
	fmt.Printf("  Deleting VPC: %s\n", vpcID)
	
	// Delete VPC using AWS CLI
	cmd := exec.Command("aws", "ec2", "delete-vpc", "--vpc-id", vpcID)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("aws cli error: %w", err)
	}
	
	fmt.Printf("  ✓ VPC deleted: %s\n", vpcID)
	return nil
}

func (e *Engine) createSubnetWithAWSCLI(resource *state.ResourceState, currentState *state.State) error {
	// Get required parameters
	cidrBlock, ok := resource.Attributes["cidr_block"].(string)
	if !ok {
		return fmt.Errorf("cidr_block not found in Subnet configuration")
	}
	
	availabilityZone, ok := resource.Attributes["availability_zone"].(string)
	if !ok {
		return fmt.Errorf("availability_zone not found in Subnet configuration")
	}
	
	// Resolve VPC ID from dependency
	vpcRef, ok := resource.Attributes["vpc_id"].(string)
	if !ok {
		return fmt.Errorf("vpc_id not found in Subnet configuration")
	}
	
	// Find the VPC resource in current state
	var vpcID string
	for _, res := range currentState.Resources {
		if res.Name == vpcRef && res.Type == "vpc" {
			if id, exists := res.Attributes["vpc_id"].(string); exists {
				vpcID = id
				break
			}
		}
	}
	
	if vpcID == "" {
		return fmt.Errorf("could not resolve VPC ID for reference: %s", vpcRef)
	}
	
	infoColor.Printf("  Creating Subnet with CIDR: %s in VPC: %s\n", cidrBlock, vpcID)
	
	// Create Subnet using AWS CLI
	cmd := exec.Command("aws", "ec2", "create-subnet", 
		"--vpc-id", vpcID,
		"--cidr-block", cidrBlock,
		"--availability-zone", availabilityZone,
		"--output", "json")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("aws cli error: %w, output: %s", err, string(output))
	}
	
	// Parse the output to get Subnet ID
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return fmt.Errorf("failed to parse AWS CLI output: %w", err)
	}
	
	subnet, ok := result["Subnet"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid AWS CLI response format")
	}
	
	subnetID, ok := subnet["SubnetId"].(string)
	if !ok {
		return fmt.Errorf("Subnet ID not found in response")
	}
	
	// Update resource attributes with actual Subnet ID
	resource.Attributes["subnet_id"] = subnetID
	resource.Attributes["vpc_id"] = vpcID // Store resolved VPC ID
	resource.Attributes["state"] = subnet["State"]
	
	successColor.Printf("  Subnet created with ID: %s\n", subnetID)
	
	// Configure public IP mapping if specified
	if mapPublicIP, exists := resource.Attributes["map_public_ip"]; exists {
		if mapPublic, ok := mapPublicIP.(bool); ok && mapPublic {
			if err := e.configureSubnetPublicIP(subnetID, true); err != nil {
				fmt.Printf("  Warning: failed to configure public IP mapping: %v\n", err)
			}
		}
	}
	
	// Add tags if specified
	if tags, exists := resource.Attributes["tags"]; exists {
		if err := e.tagSubnetWithAWSCLI(subnetID, tags); err != nil {
			fmt.Printf("  Warning: failed to tag Subnet: %v\n", err)
		}
	}
	
	return nil
}

func (e *Engine) configureSubnetPublicIP(subnetID string, mapPublicIP bool) error {
	cmd := exec.Command("aws", "ec2", "modify-subnet-attribute",
		"--subnet-id", subnetID,
		"--map-public-ip-on-launch")
	
	if !mapPublicIP {
		cmd = exec.Command("aws", "ec2", "modify-subnet-attribute",
			"--subnet-id", subnetID,
			"--no-map-public-ip-on-launch")
	}
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to configure public IP mapping: %w", err)
	}
	
	successColor.Printf("  Subnet public IP mapping configured\n")
	return nil
}

func (e *Engine) tagSubnetWithAWSCLI(subnetID string, tags interface{}) error {
	tagsMap, ok := tags.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid tags format")
	}
	
	for key, value := range tagsMap {
		tagSpec := fmt.Sprintf("Key=%s,Value=%s", key, value)
		cmd := exec.Command("aws", "ec2", "create-tags", "--resources", subnetID, "--tags", tagSpec)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to tag Subnet with %s=%s: %w", key, value, err)
		}
	}
	
	successColor.Printf("  Subnet tagged successfully\n")
	return nil
}

func (e *Engine) deleteSubnetWithAWSCLI(resource *state.ResourceState) error {
	subnetID, ok := resource.Attributes["subnet_id"].(string)
	if !ok {
		return fmt.Errorf("subnet_id not found in resource state")
	}
	
	fmt.Printf("  Deleting Subnet: %s\n", subnetID)
	
	// Delete Subnet using AWS CLI
	cmd := exec.Command("aws", "ec2", "delete-subnet", "--subnet-id", subnetID)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("aws cli error: %w", err)
	}
	
	fmt.Printf("  ✓ Subnet deleted: %s\n", subnetID)
	return nil
}