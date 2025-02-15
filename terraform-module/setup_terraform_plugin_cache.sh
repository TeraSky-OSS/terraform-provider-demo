#!/bin/bash

# Define variables
TERRAFORM_CONFIG="$HOME/.terraformrc"
PLUGIN_DIR="$HOME/.terraform.d/plugins"
BIN_DIR="$HOME/.terraform.d/bin"

# Create the bin directory
mkdir -p "$BIN_DIR"

# Create the plugins directory
mkdir -p "$PLUGIN_DIR"

# Write the configuration to the .terraformrc file
cat <<EOF > "$TERRAFORM_CONFIG"
disable_checkpoint = true
plugin_cache_dir   = "$BIN_DIR"

provider_installation {
  filesystem_mirror {
    path    = "$PLUGIN_DIR"
    include = ["registry.terraform.io/*/*", "terraform.local/*/*"]
  }
  direct {
    exclude = []
  }
}
EOF

# Display a success message
echo "Terraform plugin cache configuration created at $TERRAFORM_CONFIG"
echo "Directories created: $BIN_DIR and $PLUGIN_DIR"
