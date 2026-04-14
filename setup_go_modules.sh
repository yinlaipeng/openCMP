#!/bin/bash

# Script to properly initialize Go modules after adding LDAP dependency
# This script should be run in an environment with Go installed

cd /Users/aurora/Desktop/xtwork/git/openCMP/backend

echo "Initializing Go modules after LDAP dependency addition..."
echo "This requires Go to be installed in your environment."

# Initialize or update go.sum with the new dependencies
if command -v go &> /dev/null; then
    echo "Go is available, running go mod tidy..."
    go mod tidy

    echo "Verifying the module..."
    go mod verify

    echo "Dependencies have been updated. You can now run the server with:"
    echo "go run cmd/server/main.go"
else
    echo "Go is not installed in this environment."
    echo "Please run the following commands in a Go-enabled environment:"
    echo "cd /Users/aurora/Desktop/xtwork/git/openCMP/backend"
    echo "go mod tidy"
    echo "go run cmd/server/main.go"
fi