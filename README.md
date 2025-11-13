# TBLang - Plugin-based Infrastructure as Code

TBLang is a domain-specific language for Infrastructure as Code, designed with a plugin architecture similar to Terraform. The core handles language parsing and orchestration, while provider plugins implement cloud-specific logic.

## Architecture

### Core Project (`core/`)
- **Language & Grammar**: TBLang syntax and ANTLR grammar
- **Compiler**: Parses `.tbl` files and builds dependency graphs
- **Engine**: Orchestrates plan/apply/destroy operations
- **State Management**: Tracks infrastructure state
- **Plugin System**: Discovers and communicates with provider plugins

### Plugin Projects (`plugin/*/`)
- **AWS Plugin** (`plugin/aws/`): Amazon Web Services provider
- **Azure Plugin** (coming soon): Microsoft Azure provider
- **GCP Plugin** (coming soon): Google Cloud Platform provider

## Quick Start

### 1. Build Everything
```bash
# Build core and all plugins
make build

# Or build individually
make build-core
make build-aws-plugin
```

### 2. Install
```bash
# Install core and plugins
make install
```

### 3. Create Infrastructure
Create a `main.tbl` file:
```tblang
cloud_vendor "aws" {
    region = "us-east-1"
    account_id = "123456789012"
}

declare vpc_config = {
    cidr_block: "10.0.0.0/16"
    enable_dns_hostnames: true
    enable_dns_support: true
    tags: {
        Environment: "production"
        Project: "tblang-demo"
    }
}

declare main_vpc = vpc("production-vpc", vpc_config);
```

### 4. Run TBLang Commands
```bash
# Plan infrastructure
tblang plan main.tbl

# Apply changes
tblang apply main.tbl

# Show current state
tblang show

# Destroy infrastructure
tblang destroy main.tbl
```

## Development

### Prerequisites
- Go 1.21+
- ANTLR4 (for parser generation)
- AWS CLI configured (for AWS plugin)

### Development Setup
```bash
# Set up development environment
make dev-setup

# Generate parser from grammar (if grammar changes)
cd core && make parser
```

### Building
```bash
# Build everything
make build

# Build only core
make build-core

# Build only AWS plugin
make build-aws-plugin
```

### Testing
```bash
# Run all tests
make test

# Test individual components
cd core && make test
cd plugin/aws && make test
```

## Project Structure

```
tblang/
├── core/                    # TBLang Core Project
│   ├── cmd/tblang/         # Main CLI
│   ├── internal/
│   │   ├── compiler/       # Language compiler
│   │   ├── engine/         # Core engine
│   │   ├── state/          # State management
│   │   ├── ast/            # AST definitions
│   │   └── graph/          # Dependency graph
│   ├── pkg/plugin/         # Plugin protocol
│   ├── grammar/            # ANTLR grammar
│   └── go.mod
├── plugin/                 # Provider Plugins
│   └── aws/                # AWS Provider Plugin
│       ├── internal/provider/
│       ├── main.go
│       └── go.mod
├── main.tbl                # Example TBLang file
└── Makefile                # Root build system
```

## Plugin Development

To create a new provider plugin:

1. **Create Plugin Project**:
```bash
mkdir plugin/myprovider
cd plugin/myprovider
go mod init github.com/tblang/provider-myprovider
```

2. **Implement Provider Interface**:
```go
import "github.com/tblang/core/pkg/plugin"

type MyProvider struct{}

func (p *MyProvider) GetSchema(ctx context.Context, req *plugin.GetSchemaRequest) (*plugin.GetSchemaResponse, error) {
    // Return provider and resource schemas
}

func (p *MyProvider) ApplyResourceChange(ctx context.Context, req *plugin.ApplyResourceChangeRequest) (*plugin.ApplyResourceChangeResponse, error) {
    // Create/update resources
}

// ... implement other methods
```

3. **Create Plugin Main**:
```go
func main() {
    provider := &MyProvider{}
    server := plugin.NewServer(provider)
    server.Serve()
}
```

4. **Build and Install**:
```bash
go build -o tblang-provider-myprovider
cp tblang-provider-myprovider ../core/.tblang/plugins/
```

## Language Syntax

### Cloud Provider Configuration
```tblang
cloud_vendor "aws" {
    region = "us-east-1"
    account_id = "123456789012"
}
```

### Variable Declarations
```tblang
declare vpc_config = {
    cidr_block: "10.0.0.0/16"
    enable_dns_hostnames: true
}
```

### Resource Creation
```tblang
declare my_vpc = vpc("vpc-name", vpc_config);
declare my_subnet = subnet("subnet-name", {
    vpc_id: my_vpc
    cidr_block: "10.0.1.0/24"
});
```

## Commands

- `tblang plan <file.tbl>` - Show planned changes
- `tblang apply <file.tbl>` - Apply infrastructure changes
- `tblang destroy <file.tbl>` - Destroy infrastructure
- `tblang show` - Show current state
- `tblang graph <file.tbl>` - Show dependency graph
- `tblang plugins list` - List available plugins
- `tblang version` - Show version information

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details.