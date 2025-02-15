# Car Store Terraform Provider

This project implements a custom Terraform provider for managing a car store system. It consists of multiple components working together to provide infrastructure as code capabilities for car store management.

## Project Structure
```
.
├── api/     # Python API server
├── terraform-provider    # Terraform Provider implementation
└── terraform-module/    # Terraform Module for testing
```

## Components

### 1. API Server (Python)
Located in the `api/` directory, this is a Python-based REST API server that handles the actual car store operations. It serves as the backend that the Terraform provider communicates with.

### 2. Terraform Provider (Go)
Located in the `terraform-provider/` directory, this is the core provider implementation written in Go. It:
- Implements the necessary provider interfaces using the Terraform Plugin Framework
- Handles CRUD operations for car store resources
- Communicates with the Python API server to execute operations
- Manages state and schema validation

### 3. Terraform Test Module
Located in the `terraform-module/` directory, contains example Terraform configurations and acceptance tests to verify the provider's functionality.

## Getting Started

### Prerequisites
- Go 1.19 or later
- Python 3.8 or later
- Terraform 1.0 or later

### Building the Provider
```bash
cd terraform-provider
./build-and-cleanup.sh
```

### Running the API Server
```bash
cd api
# First time setup:
./setup.sh    # Creates virtual environment, installs dependencies, and starts the API server

# For subsequent runs:
./run.sh # Enter the virtual environment and start the API server
```

### Running Terraform test module
```bash
cd terraform-module
# First time setup:
./setup_terraform_plugin_cache.sh # Creates the terraform plugin cache


## Provider Configuration
```hcl
terraform {
  required_providers {
    carstore = {
      source = "registry.terraform.io/local/carstore"
    }
  }
}

provider "carstore" {
  # provider configuration options
}
```

## Resources
The provider currently supports the following resources:
- `carstore_car` - Manages car entries in the store

Add other resources as they are implemented...

