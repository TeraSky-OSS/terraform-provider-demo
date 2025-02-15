#!/bin/bash

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo "Error: Python 3 is not installed or not in PATH"
    exit 1
fi

# Create virtual environment
python3 -m venv venv

# Verify venv was created
if [ ! -f "venv/bin/activate" ]; then
    echo "Error: Failed to create virtual environment"
    exit 1
fi

# Activate virtual environment
source venv/bin/activate

# Install requirements
if [ -f "requirements.txt" ]; then
    pip install -r requirements.txt
else
    echo "Creating requirements.txt"
    echo "flask==3.0.2" > requirements.txt
    pip install -r requirements.txt
fi

echo "Setup complete! Virtual environment is activated."
echo "To start the server, run: python api/server.py" 