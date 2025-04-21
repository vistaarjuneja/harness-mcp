#!/bin/bash

set -e

# Build the harness-mcp-server binary
echo "Building harness-mcp-server..."
go build -o harness-mcp-server ./cmd/harness-mcp-server

echo "Build complete. You can run the server with:"
echo "HARNESS_API_KEY=your_api_key ./harness-mcp-server stdio"