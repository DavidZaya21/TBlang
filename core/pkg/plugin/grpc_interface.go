package plugin

import (
	"context"
	"encoding/json"

	"github.com/tblang/core/pkg/plugin/proto"
)

// GRPCProviderPlugin implements the ProviderPlugin interface using gRPC
type GRPCProviderPlugin interface {
	// GetSchema returns the provider and resource schemas
	GetSchema(ctx context.Context, req *GetSchemaRequest) (*GetSchemaResponse, error)
	
	// Configure configures the provider with the given configuration
	Configure(ctx context.Context, req *ConfigureRequest) (*ConfigureResponse, error)
	
	// PlanResourceChange plans changes for a resource
	PlanResourceChange(ctx context.Context, req *PlanResourceChangeRequest) (*PlanResourceChangeResponse, error)
	
	// ApplyResourceChange applies changes to a resource
	ApplyResourceChange(ctx context.Context, req *ApplyResourceChangeRequest) (*ApplyResourceChangeResponse, error)
	
	// ReadResource reads the current state of a resource
	ReadResource(ctx context.Context, req *ReadResourceRequest) (*ReadResourceResponse, error)
	
	// ImportResource imports an existing resource
	ImportResource(ctx context.Context, req *ImportResourceRequest) (*ImportResourceResponse, error)
	
	// ValidateResourceConfig validates a resource configuration
	ValidateResourceConfig(ctx context.Context, req *ValidateResourceConfigRequest) (*ValidateResourceConfigResponse, error)
}

// Helper functions to convert between proto and interface types

func ProtoToGetSchemaResponse(p *proto.GetSchemaResponse) *GetSchemaResponse {
	resp := &GetSchemaResponse{
		ResourceSchemas:   make(map[string]*Schema),
		DataSourceSchemas: make(map[string]*Schema),
		Diagnostics:       make([]*Diagnostic, len(p.Diagnostics)),
	}
	
	if p.Provider != nil {
		resp.Provider = ProtoToSchema(p.Provider)
	}
	
	for name, schema := range p.ResourceSchemas {
		resp.ResourceSchemas[name] = ProtoToSchema(schema)
	}
	
	for name, schema := range p.DataSourceSchemas {
		resp.DataSourceSchemas[name] = ProtoToSchema(schema)
	}
	
	for i, diag := range p.Diagnostics {
		resp.Diagnostics[i] = ProtoToDiagnostic(diag)
	}
	
	return resp
}

func ProtoToSchema(p *proto.Schema) *Schema {
	schema := &Schema{
		Version: p.Version,
	}
	
	if p.Block != nil {
		schema.Block = ProtoToSchemaBlock(p.Block)
	}
	
	return schema
}

func ProtoToSchemaBlock(p *proto.SchemaBlock) *SchemaBlock {
	block := &SchemaBlock{
		Attributes: make(map[string]*Attribute),
		BlockTypes: make(map[string]*BlockType),
	}
	
	for name, attr := range p.Attributes {
		block.Attributes[name] = ProtoToAttribute(attr)
	}
	
	for name, blockType := range p.BlockTypes {
		block.BlockTypes[name] = ProtoToBlockType(blockType)
	}
	
	return block
}

func ProtoToAttribute(p *proto.Attribute) *Attribute {
	return &Attribute{
		Type:        p.Type,
		Description: p.Description,
		Required:    p.Required,
		Optional:    p.Optional,
		Computed:    p.Computed,
		Sensitive:   p.Sensitive,
	}
}

func ProtoToBlockType(p *proto.BlockType) *BlockType {
	blockType := &BlockType{
		NestingMode: p.NestingMode,
		MinItems:    p.MinItems,
		MaxItems:    p.MaxItems,
	}
	
	if p.Block != nil {
		blockType.Block = ProtoToSchemaBlock(p.Block)
	}
	
	return blockType
}

func ProtoToDiagnostic(p *proto.Diagnostic) *Diagnostic {
	return &Diagnostic{
		Severity: p.Severity,
		Summary:  p.Summary,
		Detail:   p.Detail,
	}
}

func ConfigureRequestToProto(req *ConfigureRequest) *proto.ConfigureRequest {
	protoReq := &proto.ConfigureRequest{
		TerraformVersion: req.TerraformVersion,
	}
	
	if req.Config != nil {
		if jsonData, err := json.Marshal(req.Config); err == nil {
			protoReq.Config = &proto.DynamicValue{Json: jsonData}
		}
	}
	
	return protoReq
}

func ProtoToConfigureResponse(p *proto.ConfigureResponse) *ConfigureResponse {
	resp := &ConfigureResponse{
		Diagnostics: make([]*Diagnostic, len(p.Diagnostics)),
	}
	
	for i, diag := range p.Diagnostics {
		resp.Diagnostics[i] = ProtoToDiagnostic(diag)
	}
	
	return resp
}

func ApplyResourceChangeRequestToProto(req *ApplyResourceChangeRequest) *proto.ApplyResourceChangeRequest {
	protoReq := &proto.ApplyResourceChangeRequest{
		TypeName:       req.TypeName,
		PlannedPrivate: req.PlannedPrivate,
	}
	
	if req.PriorState != nil {
		if jsonData, err := json.Marshal(req.PriorState); err == nil {
			protoReq.PriorState = &proto.DynamicValue{Json: jsonData}
		}
	}
	
	if req.PlannedState != nil {
		if jsonData, err := json.Marshal(req.PlannedState); err == nil {
			protoReq.PlannedState = &proto.DynamicValue{Json: jsonData}
		}
	}
	
	if req.Config != nil {
		if jsonData, err := json.Marshal(req.Config); err == nil {
			protoReq.Config = &proto.DynamicValue{Json: jsonData}
		}
	}
	
	return protoReq
}

func ProtoToApplyResourceChangeResponse(p *proto.ApplyResourceChangeResponse) *ApplyResourceChangeResponse {
	resp := &ApplyResourceChangeResponse{
		Private:     p.Private,
		Diagnostics: make([]*Diagnostic, len(p.Diagnostics)),
	}
	
	if p.NewState != nil && len(p.NewState.Json) > 0 {
		var state interface{}
		if err := json.Unmarshal(p.NewState.Json, &state); err == nil {
			resp.NewState = state
		}
	}
	
	for i, diag := range p.Diagnostics {
		resp.Diagnostics[i] = ProtoToDiagnostic(diag)
	}
	
	return resp
}