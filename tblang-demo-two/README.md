# TBLang Demo Two

Infrastructure configuration with:
- **1 VPC**: 10.1.0.0/16
- **2 Subnets**: 
  - Public subnet (10.1.1.0/24) in us-east-1a
  - Private subnet (10.1.2.0/24) in us-east-1b
- **1 Security Group**: Web server security group with HTTP, HTTPS, and SSH access

## AWS Configuration
- **Profile**: kyaw-zin
- **Region**: us-east-1

## Commands

### Plan
```bash
cd tblang-demo-two
tblang plan infrastructure.tbl
```

### Apply
```bash
tblang apply infrastructure.tbl
# You will be prompted: "Do you want to perform these actions? (yes/no):"
# Type: yes
```

### Show State
```bash
tblang show
```

### Destroy
```bash
tblang destroy infrastructure.tbl
# You will be prompted: "Do you really want to destroy all resources? (yes/no):"
# Type: yes
```

## Verify in AWS

### Check VPC
```bash
aws ec2 describe-vpcs --profile kyaw-zin --filters "Name=tag:Project,Values=tblang-demo-two" --query 'Vpcs[*].[VpcId,CidrBlock,Tags[?Key==`Name`].Value|[0]]' --output table
```

### Check Subnets
```bash
aws ec2 describe-subnets --profile kyaw-zin --filters "Name=tag:Project,Values=tblang-demo-two" --query 'Subnets[*].[SubnetId,CidrBlock,AvailabilityZone,Tags[?Key==`Name`].Value|[0]]' --output table
```

### Check Security Group
```bash
aws ec2 describe-security-groups --profile kyaw-zin --filters "Name=tag:Project,Values=tblang-demo-two" --query 'SecurityGroups[*].[GroupId,GroupName,Description]' --output table
```
