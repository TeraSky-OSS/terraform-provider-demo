#!/bin/bash

# Define variables
PLUGIN_DIR="$HOME/.terraform.d/plugins/registry.terraform.io/local/carstore/1.0.0/darwin_arm64"
PROVIDER_BINARY="terraform-provider-carstore"
PROJECT_DIR="./"
TEST_DIR="../terraform-module"

# Create the plugin directory
mkdir -p "$PLUGIN_DIR"

# Build the Go provider
cd "$PROJECT_DIR" || {
    echo "Project directory not found: $PROJECT_DIR"
    exit 1
}

go build -o "$PLUGIN_DIR/$PROVIDER_BINARY"

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Provider built successfully at $PLUGIN_DIR/$PROVIDER_BINARY"
else
    echo "Failed to build the provider"
    exit 1
fi

# Clean up Terraform files in the test directory
cd "$TEST_DIR" || {
    echo "Test directory not found: $TEST_DIR"
    exit 1
}

rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup

# Confirm cleanup
echo "Cleaned up Terraform files in $TEST_DIR"
