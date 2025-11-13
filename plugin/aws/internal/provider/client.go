package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// AWSClient wraps AWS SDK clients for the plugin
type AWSClient struct {
	EC2    *ec2.Client
	Config aws.Config
	Region string
}

// VPCResult represents the result of VPC creation
type VPCResult struct {
	VpcID string
	State string
}

// NewAWSClient creates a new AWS client for the plugin
func NewAWSClient(region string) (*AWSClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &AWSClient{
		EC2:    ec2.NewFromConfig(cfg),
		Config: cfg,
		Region: region,
	}, nil
}

// CreateVPC creates a VPC and returns the result
func (c *AWSClient) CreateVPC(ctx context.Context, cidrBlock string, tags map[string]string) (*VPCResult, error) {
	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(cidrBlock),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeVpc,
				Tags:         c.buildTags("vpc", tags),
			},
		},
	}

	result, err := c.EC2.CreateVpc(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create VPC: %w", err)
	}

	if result.Vpc == nil {
		return nil, fmt.Errorf("VPC creation returned nil")
	}

	// Enable DNS hostnames and support
	vpcID := *result.Vpc.VpcId
	if err := c.enableVPCDNS(ctx, vpcID); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: failed to enable DNS for VPC %s: %v\n", vpcID, err)
	}

	return &VPCResult{
		VpcID: vpcID,
		State: string(result.Vpc.State),
	}, nil
}

// enableVPCDNS enables DNS hostnames and support for a VPC
func (c *AWSClient) enableVPCDNS(ctx context.Context, vpcID string) error {
	// Enable DNS hostnames
	_, err := c.EC2.ModifyVpcAttribute(ctx, &ec2.ModifyVpcAttributeInput{
		VpcId:              aws.String(vpcID),
		EnableDnsHostnames: &types.AttributeBooleanValue{Value: aws.Bool(true)},
	})
	if err != nil {
		return fmt.Errorf("failed to enable DNS hostnames: %w", err)
	}

	// Enable DNS support
	_, err = c.EC2.ModifyVpcAttribute(ctx, &ec2.ModifyVpcAttributeInput{
		VpcId:           aws.String(vpcID),
		EnableDnsSupport: &types.AttributeBooleanValue{Value: aws.Bool(true)},
	})
	if err != nil {
		return fmt.Errorf("failed to enable DNS support: %w", err)
	}

	return nil
}

// SubnetResult represents the result of Subnet creation
type SubnetResult struct {
	SubnetID string
	State    string
}

// CreateSubnet creates a subnet and returns the result
func (c *AWSClient) CreateSubnet(ctx context.Context, vpcID, cidrBlock, availabilityZone string, tags map[string]string) (*SubnetResult, error) {
	input := &ec2.CreateSubnetInput{
		VpcId:            aws.String(vpcID),
		CidrBlock:        aws.String(cidrBlock),
		AvailabilityZone: aws.String(availabilityZone),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeSubnet,
				Tags:         c.buildTags("subnet", tags),
			},
		},
	}

	result, err := c.EC2.CreateSubnet(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create Subnet: %w", err)
	}

	if result.Subnet == nil {
		return nil, fmt.Errorf("Subnet creation returned nil")
	}

	return &SubnetResult{
		SubnetID: *result.Subnet.SubnetId,
		State:    string(result.Subnet.State),
	}, nil
}

// ConfigureSubnetPublicIP configures public IP mapping for a subnet
func (c *AWSClient) ConfigureSubnetPublicIP(ctx context.Context, subnetID string, mapPublicIP bool) error {
	input := &ec2.ModifySubnetAttributeInput{
		SubnetId: aws.String(subnetID),
		MapPublicIpOnLaunch: &types.AttributeBooleanValue{
			Value: aws.Bool(mapPublicIP),
		},
	}

	_, err := c.EC2.ModifySubnetAttribute(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to configure public IP mapping: %w", err)
	}

	return nil
}

// DeleteVPC deletes a VPC
func (c *AWSClient) DeleteVPC(ctx context.Context, vpcID string) error {
	input := &ec2.DeleteVpcInput{
		VpcId: aws.String(vpcID),
	}

	_, err := c.EC2.DeleteVpc(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete VPC %s: %w", vpcID, err)
	}

	return nil
}

// SecurityGroupResult represents the result of Security Group creation
type SecurityGroupResult struct {
	GroupID string
}

// CreateSecurityGroup creates a security group and returns the result
func (c *AWSClient) CreateSecurityGroup(ctx context.Context, vpcID, name, description string, tags map[string]string) (*SecurityGroupResult, error) {
	input := &ec2.CreateSecurityGroupInput{
		VpcId:       aws.String(vpcID),
		GroupName:   aws.String(name),
		Description: aws.String(description),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeSecurityGroup,
				Tags:         c.buildTags("security-group", tags),
			},
		},
	}

	result, err := c.EC2.CreateSecurityGroup(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create Security Group: %w", err)
	}

	if result.GroupId == nil {
		return nil, fmt.Errorf("Security Group creation returned nil")
	}

	return &SecurityGroupResult{
		GroupID: *result.GroupId,
	}, nil
}

// AuthorizeSecurityGroupIngress adds ingress rules to a security group
func (c *AWSClient) AuthorizeSecurityGroupIngress(ctx context.Context, groupID string, rules []types.IpPermission) error {
	if len(rules) == 0 {
		return nil
	}

	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:       aws.String(groupID),
		IpPermissions: rules,
	}

	_, err := c.EC2.AuthorizeSecurityGroupIngress(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to authorize ingress rules: %w", err)
	}

	return nil
}

// AuthorizeSecurityGroupEgress adds egress rules to a security group
func (c *AWSClient) AuthorizeSecurityGroupEgress(ctx context.Context, groupID string, rules []types.IpPermission) error {
	if len(rules) == 0 {
		return nil
	}

	input := &ec2.AuthorizeSecurityGroupEgressInput{
		GroupId:       aws.String(groupID),
		IpPermissions: rules,
	}

	_, err := c.EC2.AuthorizeSecurityGroupEgress(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to authorize egress rules: %w", err)
	}

	return nil
}

// DeleteSecurityGroup deletes a security group
func (c *AWSClient) DeleteSecurityGroup(ctx context.Context, groupID string) error {
	input := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(groupID),
	}

	_, err := c.EC2.DeleteSecurityGroup(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete Security Group %s: %w", groupID, err)
	}

	return nil
}

// DeleteSubnet deletes a subnet
func (c *AWSClient) DeleteSubnet(ctx context.Context, subnetID string) error {
	input := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(subnetID),
	}

	_, err := c.EC2.DeleteSubnet(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete Subnet %s: %w", subnetID, err)
	}

	return nil
}

// buildTags creates AWS tags from a map
func (c *AWSClient) buildTags(resourceName string, additionalTags map[string]string) []types.Tag {
	// Start with additional tags
	tagMap := make(map[string]string)
	for key, value := range additionalTags {
		tagMap[key] = value
	}
	
	// Add default tags only if not already present
	if _, exists := tagMap["ManagedBy"]; !exists {
		tagMap["ManagedBy"] = "TBLang"
	}
	
	// Convert to AWS tags
	var tags []types.Tag
	for key, value := range tagMap {
		tags = append(tags, types.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}

	return tags
}