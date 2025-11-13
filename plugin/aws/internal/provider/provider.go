package provider

import (
	"context"
	"fmt"

	"github.com/tblang/core/pkg/plugin"
)

// AWSProvider implements the TBLang gRPC provider plugin interface for AWS
type AWSProvider struct {
	region    string
	accountID string
	client    *AWSClient
}

// NewAWSProvider creates a new AWS provider
func NewAWSProvider() *AWSProvider {
	return &AWSProvider{}
}

// GetSchema returns the provider and resource schemas
func (p *AWSProvider) GetSchema(ctx context.Context, req *plugin.GetSchemaRequest) (*plugin.GetSchemaResponse, error) {
	return &plugin.GetSchemaResponse{
		Provider: &plugin.Schema{
			Version: 1,
			Block: &plugin.SchemaBlock{
				Attributes: map[string]*plugin.Attribute{
					"region": {
						Type:        "string",
						Description: "AWS region",
						Required:    true,
					},
					"account_id": {
						Type:        "string", 
						Description: "AWS account ID",
						Optional:    true,
					},
				},
			},
		},
		ResourceSchemas: map[string]*plugin.Schema{
			"vpc": {
				Version: 1,
				Block: &plugin.SchemaBlock{
					Attributes: map[string]*plugin.Attribute{
						"cidr_block": {
							Type:        "string",
							Description: "CIDR block for VPC",
							Required:    true,
						},
						"enable_dns_hostnames": {
							Type:        "bool",
							Description: "Enable DNS hostnames",
							Optional:    true,
						},
						"enable_dns_support": {
							Type:        "bool", 
							Description: "Enable DNS support",
							Optional:    true,
						},
						"tags": {
							Type:        "map",
							Description: "Resource tags",
							Optional:    true,
						},
						"vpc_id": {
							Type:        "string",
							Description: "VPC ID",
							Computed:    true,
						},
					},
				},
			},
			"subnet": {
				Version: 1,
				Block: &plugin.SchemaBlock{
					Attributes: map[string]*plugin.Attribute{
						"vpc_id": {
							Type:        "string",
							Description: "VPC ID",
							Required:    true,
						},
						"cidr_block": {
							Type:        "string",
							Description: "CIDR block for subnet",
							Required:    true,
						},
						"availability_zone": {
							Type:        "string",
							Description: "Availability zone",
							Required:    true,
						},
						"map_public_ip": {
							Type:        "bool",
							Description: "Map public IP on launch",
							Optional:    true,
						},
						"tags": {
							Type:        "map",
							Description: "Resource tags",
							Optional:    true,
						},
						"subnet_id": {
							Type:        "string",
							Description: "Subnet ID",
							Computed:    true,
						},
					},
				},
			},
			"security_group": {
				Version: 1,
				Block: &plugin.SchemaBlock{
					Attributes: map[string]*plugin.Attribute{
						"vpc_id": {
							Type:        "string",
							Description: "VPC ID",
							Required:    true,
						},
						"name": {
							Type:        "string",
							Description: "Security group name",
							Required:    true,
						},
						"description": {
							Type:        "string",
							Description: "Security group description",
							Optional:    true,
						},
						"ingress_rules": {
							Type:        "list",
							Description: "Ingress rules",
							Optional:    true,
						},
						"egress_rules": {
							Type:        "list",
							Description: "Egress rules",
							Optional:    true,
						},
						"tags": {
							Type:        "map",
							Description: "Resource tags",
							Optional:    true,
						},
						"group_id": {
							Type:        "string",
							Description: "Security group ID",
							Computed:    true,
						},
					},
				},
			},
			"ec2": {
				Version: 1,
				Block: &plugin.SchemaBlock{
					Attributes: map[string]*plugin.Attribute{
						"ami": {
							Type:        "string",
							Description: "AMI ID",
							Required:    true,
						},
						"instance_type": {
							Type:        "string",
							Description: "Instance type",
							Required:    true,
						},
						"subnet_id": {
							Type:        "string",
							Description: "Subnet ID",
							Required:    true,
						},
						"security_groups": {
							Type:        "list",
							Description: "Security group IDs",
							Optional:    true,
						},
						"key_name": {
							Type:        "string",
							Description: "Key pair name",
							Optional:    true,
						},
						"user_data": {
							Type:        "string",
							Description: "User data script",
							Optional:    true,
						},
						"tags": {
							Type:        "map",
							Description: "Resource tags",
							Optional:    true,
						},
						"instance_id": {
							Type:        "string",
							Description: "Instance ID",
							Computed:    true,
						},
						"public_ip": {
							Type:        "string",
							Description: "Public IP address",
							Computed:    true,
						},
						"private_ip": {
							Type:        "string",
							Description: "Private IP address",
							Computed:    true,
						},
					},
				},
			},
		},
	}, nil
}

// Configure configures the provider with the given configuration
func (p *AWSProvider) Configure(ctx context.Context, req *plugin.ConfigureRequest) (*plugin.ConfigureResponse, error) {
	config, ok := req.Config.(map[string]interface{})
	if !ok {
		return &plugin.ConfigureResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Invalid configuration",
					Detail:   "Configuration must be a map",
				},
			},
		}, nil
	}

	// Extract region
	if region, exists := config["region"]; exists {
		if regionStr, ok := region.(string); ok {
			p.region = regionStr
		}
	}

	// Extract account ID
	if accountID, exists := config["account_id"]; exists {
		if accountIDStr, ok := accountID.(string); ok {
			p.accountID = accountIDStr
		}
	}

	// Initialize AWS client
	client, err := NewAWSClient(p.region)
	if err != nil {
		return &plugin.ConfigureResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Failed to initialize AWS client",
					Detail:   err.Error(),
				},
			},
		}, nil
	}

	p.client = client

	return &plugin.ConfigureResponse{}, nil
}

// PlanResourceChange plans changes for a resource
func (p *AWSProvider) PlanResourceChange(ctx context.Context, req *plugin.PlanResourceChangeRequest) (*plugin.PlanResourceChangeResponse, error) {
	// For now, just return the proposed state
	// In a real implementation, this would validate the configuration
	// and determine what changes are needed
	return &plugin.PlanResourceChangeResponse{
		PlannedState: req.ProposedNewState,
	}, nil
}

// ApplyResourceChange applies changes to a resource
func (p *AWSProvider) ApplyResourceChange(ctx context.Context, req *plugin.ApplyResourceChangeRequest) (*plugin.ApplyResourceChangeResponse, error) {
	if p.client == nil {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Provider not configured",
					Detail:   "AWS provider must be configured before use",
				},
			},
		}, nil
	}

	// Check if this is a destroy operation (planned state is nil/empty but prior state exists)
	isDestroy := req.PlannedState == nil && req.PriorState != nil

	if isDestroy {
		// Handle destroy operations
		switch req.TypeName {
		case "vpc":
			return p.destroyVPC(ctx, req)
		case "subnet":
			return p.destroySubnet(ctx, req)
		case "security_group":
			return p.destroySecurityGroup(ctx, req)
		case "ec2":
			return p.destroyEC2(ctx, req)
		default:
			return &plugin.ApplyResourceChangeResponse{
				Diagnostics: []*plugin.Diagnostic{
					{
						Severity: "error",
						Summary:  "Unsupported resource type",
						Detail:   fmt.Sprintf("Resource type %s is not supported", req.TypeName),
					},
				},
			}, nil
		}
	}

	// Handle create/update operations
	switch req.TypeName {
	case "vpc":
		return p.applyVPC(ctx, req)
	case "subnet":
		return p.applySubnet(ctx, req)
	case "security_group":
		return p.applySecurityGroup(ctx, req)
	case "ec2":
		return p.applyEC2(ctx, req)
	default:
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Unsupported resource type",
					Detail:   fmt.Sprintf("Resource type %s is not supported", req.TypeName),
				},
			},
		}, nil
	}
}

// ReadResource reads the current state of a resource
func (p *AWSProvider) ReadResource(ctx context.Context, req *plugin.ReadResourceRequest) (*plugin.ReadResourceResponse, error) {
	// TODO: Implement resource reading from AWS
	return &plugin.ReadResourceResponse{
		NewState: req.CurrentState,
	}, nil
}

// ImportResource imports an existing resource
func (p *AWSProvider) ImportResource(ctx context.Context, req *plugin.ImportResourceRequest) (*plugin.ImportResourceResponse, error) {
	// TODO: Implement resource import
	return &plugin.ImportResourceResponse{}, nil
}

// ValidateResourceConfig validates a resource configuration
func (p *AWSProvider) ValidateResourceConfig(ctx context.Context, req *plugin.ValidateResourceConfigRequest) (*plugin.ValidateResourceConfigResponse, error) {
	// TODO: Implement configuration validation
	return &plugin.ValidateResourceConfigResponse{}, nil
}

// Resource-specific apply methods

func (p *AWSProvider) applyVPC(ctx context.Context, req *plugin.ApplyResourceChangeRequest) (*plugin.ApplyResourceChangeResponse, error) {
	config, ok := req.Config.(map[string]interface{})
	if !ok {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Invalid VPC configuration",
					Detail:   "Configuration must be a map",
				},
			},
		}, nil
	}

	cidrBlock, _ := config["cidr_block"].(string)
	if cidrBlock == "" {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Missing required field",
					Detail:   "cidr_block is required for VPC",
				},
			},
		}, nil
	}

	// Create VPC using AWS client
	vpc, err := p.client.CreateVPC(ctx, cidrBlock, extractTags(config))
	if err != nil {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Failed to create VPC",
					Detail:   err.Error(),
				},
			},
		}, nil
	}

	// Return new state
	newState := make(map[string]interface{})
	for k, v := range config {
		newState[k] = v
	}
	newState["vpc_id"] = vpc.VpcID
	newState["state"] = vpc.State

	return &plugin.ApplyResourceChangeResponse{
		NewState: newState,
	}, nil
}

func (p *AWSProvider) applySubnet(ctx context.Context, req *plugin.ApplyResourceChangeRequest) (*plugin.ApplyResourceChangeResponse, error) {
	config, ok := req.Config.(map[string]interface{})
	if !ok {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Invalid Subnet configuration",
					Detail:   "Configuration must be a map",
				},
			},
		}, nil
	}

	vpcID, _ := config["vpc_id"].(string)
	cidrBlock, _ := config["cidr_block"].(string)
	availabilityZone, _ := config["availability_zone"].(string)
	
	if vpcID == "" || cidrBlock == "" || availabilityZone == "" {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Missing required fields",
					Detail:   "vpc_id, cidr_block, and availability_zone are required for Subnet",
				},
			},
		}, nil
	}

	// Create Subnet using AWS client
	subnet, err := p.client.CreateSubnet(ctx, vpcID, cidrBlock, availabilityZone, extractTags(config))
	if err != nil {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Failed to create Subnet",
					Detail:   err.Error(),
				},
			},
		}, nil
	}

	// Configure public IP mapping if specified
	if mapPublicIP, exists := config["map_public_ip"]; exists {
		if mapPublic, ok := mapPublicIP.(bool); ok && mapPublic {
			if err := p.client.ConfigureSubnetPublicIP(ctx, subnet.SubnetID, true); err != nil {
				// Log warning but don't fail
				fmt.Printf("Warning: failed to configure public IP mapping: %v\n", err)
			}
		}
	}

	// Return new state
	newState := make(map[string]interface{})
	for k, v := range config {
		newState[k] = v
	}
	newState["subnet_id"] = subnet.SubnetID
	newState["state"] = subnet.State

	return &plugin.ApplyResourceChangeResponse{
		NewState: newState,
	}, nil
}

func (p *AWSProvider) applySecurityGroup(ctx context.Context, req *plugin.ApplyResourceChangeRequest) (*plugin.ApplyResourceChangeResponse, error) {
	config, ok := req.Config.(map[string]interface{})
	if !ok {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Invalid Security Group configuration",
					Detail:   "Configuration must be a map",
				},
			},
		}, nil
	}

	vpcID, _ := config["vpc_id"].(string)
	name, _ := config["name"].(string)
	description, _ := config["description"].(string)
	
	if vpcID == "" || name == "" {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Missing required fields",
					Detail:   "vpc_id and name are required for Security Group",
				},
			},
		}, nil
	}

	if description == "" {
		description = "Managed by TBLang"
	}

	// Create Security Group using AWS client
	sg, err := p.client.CreateSecurityGroup(ctx, vpcID, name, description, extractTags(config))
	if err != nil {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Failed to create Security Group",
					Detail:   err.Error(),
				},
			},
		}, nil
	}

	// Add ingress rules if specified
	if ingressRules, exists := config["ingress_rules"]; exists {
		if rules, ok := ingressRules.([]interface{}); ok && len(rules) > 0 {
			// Convert rules to AWS format
			// For now, skip rule creation - would need proper conversion
			fmt.Printf("  Note: Ingress rules configuration found but not yet implemented\n")
		}
	}

	// Add egress rules if specified
	if egressRules, exists := config["egress_rules"]; exists {
		if rules, ok := egressRules.([]interface{}); ok && len(rules) > 0 {
			// Convert rules to AWS format
			// For now, skip rule creation - would need proper conversion
			fmt.Printf("  Note: Egress rules configuration found but not yet implemented\n")
		}
	}

	// Return new state
	newState := make(map[string]interface{})
	for k, v := range config {
		newState[k] = v
	}
	newState["group_id"] = sg.GroupID

	return &plugin.ApplyResourceChangeResponse{
		NewState: newState,
	}, nil
}

func (p *AWSProvider) applyEC2(ctx context.Context, req *plugin.ApplyResourceChangeRequest) (*plugin.ApplyResourceChangeResponse, error) {
	// TODO: Implement EC2 instance creation
	return &plugin.ApplyResourceChangeResponse{
		NewState: req.PlannedState,
	}, nil
}

// Resource-specific destroy methods

func (p *AWSProvider) destroyVPC(ctx context.Context, req *plugin.ApplyResourceChangeRequest) (*plugin.ApplyResourceChangeResponse, error) {
	priorState, ok := req.PriorState.(map[string]interface{})
	if !ok {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Invalid prior state",
					Detail:   "Prior state must be a map",
				},
			},
		}, nil
	}

	vpcID, ok := priorState["vpc_id"].(string)
	if !ok || vpcID == "" {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Missing VPC ID",
					Detail:   "vpc_id is required to destroy VPC",
				},
			},
		}, nil
	}

	// Delete VPC using AWS client
	if err := p.client.DeleteVPC(ctx, vpcID); err != nil {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Failed to delete VPC",
					Detail:   err.Error(),
				},
			},
		}, nil
	}

	// Return nil state to indicate resource is destroyed
	return &plugin.ApplyResourceChangeResponse{
		NewState: nil,
	}, nil
}

func (p *AWSProvider) destroySubnet(ctx context.Context, req *plugin.ApplyResourceChangeRequest) (*plugin.ApplyResourceChangeResponse, error) {
	priorState, ok := req.PriorState.(map[string]interface{})
	if !ok {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Invalid prior state",
					Detail:   "Prior state must be a map",
				},
			},
		}, nil
	}

	subnetID, ok := priorState["subnet_id"].(string)
	if !ok || subnetID == "" {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Missing Subnet ID",
					Detail:   "subnet_id is required to destroy Subnet",
				},
			},
		}, nil
	}

	// Delete Subnet using AWS client
	if err := p.client.DeleteSubnet(ctx, subnetID); err != nil {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Failed to delete Subnet",
					Detail:   err.Error(),
				},
			},
		}, nil
	}

	// Return nil state to indicate resource is destroyed
	return &plugin.ApplyResourceChangeResponse{
		NewState: nil,
	}, nil
}

func (p *AWSProvider) destroySecurityGroup(ctx context.Context, req *plugin.ApplyResourceChangeRequest) (*plugin.ApplyResourceChangeResponse, error) {
	priorState, ok := req.PriorState.(map[string]interface{})
	if !ok {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Invalid prior state",
					Detail:   "Prior state must be a map",
				},
			},
		}, nil
	}

	groupID, ok := priorState["group_id"].(string)
	if !ok || groupID == "" {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Missing Security Group ID",
					Detail:   "group_id is required to destroy Security Group",
				},
			},
		}, nil
	}

	// Delete Security Group using AWS client
	if err := p.client.DeleteSecurityGroup(ctx, groupID); err != nil {
		return &plugin.ApplyResourceChangeResponse{
			Diagnostics: []*plugin.Diagnostic{
				{
					Severity: "error",
					Summary:  "Failed to delete Security Group",
					Detail:   err.Error(),
				},
			},
		}, nil
	}

	// Return nil state to indicate resource is destroyed
	return &plugin.ApplyResourceChangeResponse{
		NewState: nil,
	}, nil
}

func (p *AWSProvider) destroyEC2(ctx context.Context, req *plugin.ApplyResourceChangeRequest) (*plugin.ApplyResourceChangeResponse, error) {
	// TODO: Implement EC2 instance termination
	return &plugin.ApplyResourceChangeResponse{
		NewState: nil,
	}, nil
}

// Helper functions

func extractTags(config map[string]interface{}) map[string]string {
	tags := make(map[string]string)
	
	if tagsInterface, exists := config["tags"]; exists {
		if tagsMap, ok := tagsInterface.(map[string]interface{}); ok {
			for k, v := range tagsMap {
				if str, ok := v.(string); ok {
					tags[k] = str
				}
			}
		}
	}
	
	return tags
}