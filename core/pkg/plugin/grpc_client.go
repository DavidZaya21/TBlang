package plugin

import (
	"context"
	"fmt"

	"github.com/tblang/core/pkg/plugin/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClient implements ProviderPlugin interface over gRPC
type GRPCClient struct {
	client proto.ProviderClient
	conn   *grpc.ClientConn
}

// NewGRPCClient creates a new gRPC client
func NewGRPCClient(address string) (*GRPCClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to plugin: %w", err)
	}

	client := proto.NewProviderClient(conn)
	
	return &GRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

// GetSchema calls the remote GetSchema method
func (c *GRPCClient) GetSchema(ctx context.Context, req *GetSchemaRequest) (*GetSchemaResponse, error) {
	protoReq := &proto.GetSchemaRequest{}
	
	protoResp, err := c.client.GetSchema(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	
	return ProtoToGetSchemaResponse(protoResp), nil
}

// Configure calls the remote Configure method
func (c *GRPCClient) Configure(ctx context.Context, req *ConfigureRequest) (*ConfigureResponse, error) {
	protoReq := ConfigureRequestToProto(req)
	
	protoResp, err := c.client.Configure(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	
	return ProtoToConfigureResponse(protoResp), nil
}

// PlanResourceChange calls the remote PlanResourceChange method
func (c *GRPCClient) PlanResourceChange(ctx context.Context, req *PlanResourceChangeRequest) (*PlanResourceChangeResponse, error) {
	protoReq := &proto.PlanResourceChangeRequest{
		TypeName: req.TypeName,
	}
	
	protoResp, err := c.client.PlanResourceChange(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	
	return &PlanResourceChangeResponse{
		PlannedState:    req.ProposedNewState, // For now, just return the proposed state
		RequiresReplace: protoResp.RequiresReplace,
		PlannedPrivate:  protoResp.PlannedPrivate,
	}, nil
}

// ApplyResourceChange calls the remote ApplyResourceChange method
func (c *GRPCClient) ApplyResourceChange(ctx context.Context, req *ApplyResourceChangeRequest) (*ApplyResourceChangeResponse, error) {
	protoReq := ApplyResourceChangeRequestToProto(req)
	
	protoResp, err := c.client.ApplyResourceChange(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	
	return ProtoToApplyResourceChangeResponse(protoResp), nil
}

// ReadResource calls the remote ReadResource method
func (c *GRPCClient) ReadResource(ctx context.Context, req *ReadResourceRequest) (*ReadResourceResponse, error) {
	protoReq := &proto.ReadResourceRequest{
		TypeName: req.TypeName,
		Private:  req.Private,
	}
	
	protoResp, err := c.client.ReadResource(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	
	return &ReadResourceResponse{
		NewState: req.CurrentState, // For now, just return the current state
		Private:  protoResp.Private,
	}, nil
}

// ImportResource calls the remote ImportResource method
func (c *GRPCClient) ImportResource(ctx context.Context, req *ImportResourceRequest) (*ImportResourceResponse, error) {
	protoReq := &proto.ImportResourceRequest{
		TypeName: req.TypeName,
		Id:       req.Id,
	}
	
	protoResp, err := c.client.ImportResource(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	
	return &ImportResourceResponse{
		ImportedResources: make([]*ImportedResource, len(protoResp.ImportedResources)),
	}, nil
}

// ValidateResourceConfig calls the remote ValidateResourceConfig method
func (c *GRPCClient) ValidateResourceConfig(ctx context.Context, req *ValidateResourceConfigRequest) (*ValidateResourceConfigResponse, error) {
	protoReq := &proto.ValidateResourceConfigRequest{
		TypeName: req.TypeName,
	}
	
	protoResp, err := c.client.ValidateResourceConfig(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	
	return &ValidateResourceConfigResponse{
		Diagnostics: make([]*Diagnostic, len(protoResp.Diagnostics)),
	}, nil
}

// Close closes the gRPC client connection
func (c *GRPCClient) Close() error {
	return c.conn.Close()
}