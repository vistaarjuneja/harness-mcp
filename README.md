# Harness MCP Server

The Harness MCP Server is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) server that provides seamless integration with Harness APIs, enabling advanced automation and interaction capabilities for developers and tools.

## Components

### Tools

The server implements one tool:
- list-connector: List available connectors
  - Takes "connector_type", "connector_names", and "connector_ids" as optional arguments
  - Returns a list of connectors in JSON format

## Prerequisites

1. You will need to have Go 1.23 or later installed on your system.
2. A Harness API key for authentication.

## Quickstart

### Build from source

1. Clone the repository:
```bash
git clone https://github.com/yourusername/harness-mcp.git
cd harness-mcp
```

2. Build the binary:
```bash
go build -o harness-mcp-server ./cmd/harness-mcp-server
```

3. Run the server:
```bash
HARNESS_API_KEY=your_api_key ./harness-mcp-server stdio
```

### Claude Desktop Configuration

On MacOS: `~/Library/Application\ Support/Claude/claude_desktop_config.json`  
On Windows: `%APPDATA%/Claude/claude_desktop_config.json`

<details>
  <summary>Server Configuration</summary>
  
  ```json
  {
    "mcpServers": {
      "harness": {
        "command": "/path/to/harness-mcp-server",
        "args": ["stdio"],
        "env": {
          "HARNESS_API_KEY": "<YOUR_API_KEY>"
        }
      }
    }
  }
  ```
</details>

## Usage with Claude Code

```bash
HARNESS_API_KEY=your_api_key /path/to/harness-mcp-server stdio
```

## Development

### Command Line Arguments

The Harness MCP Server supports the following command line arguments:

- `--toolsets`: Comma-separated list of tool groups to enable (default: "all")
- `--read-only`: Run the server in read-only mode
- `--log-file`: Path to log file for debugging
- `--enable-command-logging`: Enable logging of all commands
- `--version`: Show version information
- `--help`: Show help message

### Environment Variables

Environment variables are prefixed with `HARNESS_`:

- `HARNESS_API_KEY`: Harness API key
- `HARNESS_TOOLSETS`: Comma-separated list of toolsets to enable
- `HARNESS_READ_ONLY`: Set to "true" to run in read-only mode
- `HARNESS_LOG_FILE`: Path to log file
- `HARNESS_ENABLE_COMMAND_LOGGING`: Set to "true" to enable command logging

### Authentication

The server uses a Harness API key for authentication. This can be set via the `HARNESS_API_KEY` environment variable.

## Debugging

Since MCP servers run over stdio, debugging can be challenging. For the best debugging experience, we strongly recommend using the [MCP Inspector](https://github.com/modelcontextprotocol/inspector).

You can launch the MCP Inspector with this command:

```bash
npx @modelcontextprotocol/inspector /path/to/harness-mcp-server stdio
```

Upon launching, the Inspector will display a URL that you can access in your browser to begin debugging.

## To do

We need to figure out a good way to handle auth configuration, it's just hardcoded for this POC.