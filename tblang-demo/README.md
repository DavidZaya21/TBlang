# TBLang Demo - Infrastructure as Code

This directory demonstrates how to use TBLang to manage AWS infrastructure.

## What is TBLang?

TBLang is a domain-specific language for Infrastructure as Code that provides:
- Simple, readable syntax
- Plugin-based architecture for multiple cloud providers
- State management and dependency resolution
- Terraform-like workflow (plan → apply → destroy)

## Installation

TBLang is installed system-wide at `/usr/local/bin/tblang` with plugins in `/usr/local/lib/tblang/plugins/`.

## Usage

### 1. Check Installation
```bash
tblang version
tblang plugins list
```

### 2. Plan Infrastructure
```bash
tblang plan infrastructure.tbl
```

### 3. Apply Changes
```bash
tblang apply infrastructure.tbl
```

### 4. View Current State
```bash
tblang show
```

### 5. Destroy Infrastructure
```bash
tblang destroy infrastructure.tbl
```

## TBLang Syntax

### Cloud Provider Configuration
```tblang
cloud_vendor "aws" {
    region = "us-east-1"
}
```

### Variable Declarations
```tblang
declare vpc_config = {
    cidr_block: "10.0.0.0/16"
    enable_dns_hostnames: true
    tags: {
        Name: "my-vpc"
        Environment: "production"
    }
}
```

### Resource Creation
```tblang
declare my_vpc = vpc("vpc-name", vpc_config);
```

## Current Infrastructure

This demo creates a complete AWS network infrastructure:

### VPC
- **Name**: `production-vpc`
- **CIDR**: `10.0.0.0/16`
- **DNS**: Hostnames and support enabled
- **Region**: `eu-west-1`

### Subnets
- **Public Subnet**: `production-public-subnet-1a`
  - CIDR: `10.0.1.0/24`
  - AZ: `eu-west-1a`
  - Public IP mapping: Enabled
  
- **Private Subnet**: `production-private-subnet-1b`
  - CIDR: `10.0.2.0/24`
  - AZ: `eu-west-1b`
  - Public IP mapping: Disabled
  
- **Database Subnet**: `production-db-subnet-1c`
  - CIDR: `10.0.3.0/24`
  - AZ: `eu-west-1c`
  - Public IP mapping: Disabled

### Tags
All resources are tagged with:
- `Environment`: production
- `Project`: tblang-demo
- `ManagedBy`: TBLang
- `Name`: Resource-specific name
- `Type`: Resource type (public/private/database)

## State Management

TBLang maintains state in `.tblang/tblang.tfstate` to track:
- Resource names and types
- AWS resource IDs
- Resource attributes and status

## Cleanup

To destroy the infrastructure:
```bash
tblang destroy infrastructure.tbl
```