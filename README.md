# Harness MCP Server

The Harness MCP Server is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) server that provides seamless integration with Harness APIs, enabling advanced automation and interaction capabilities for developers and tools.

## Components

### Tools

The server implements several toolsets:

#### Pipelines Toolset
- `get_pipeline`: Get details of a specific pipeline
- `list_pipelines`: List pipelines in a repository
- `get_execution`: Get details of a specific pipeline execution
- `list_executions`: List pipeline executions
- `fetch_execution_url`: Fetch the execution URL for a pipeline execution

#### Pull Requests Toolset
- `get_pull_request`: Get details of a specific pull request
- `list_pull_requests`: List pull requests in a repository
- `get_pull_request_checks`: Get status checks for a specific pull request
- `create_pull_request`: Create a new pull request

#### Repositories Toolset
- `get_repository`: Get details of a specific repository
- `list_repositories`: List repositories

#### Logs Toolset
- `download_execution_logs`: Download logs for a pipeline execution

## Prerequisites

1. You will need to have Go 1.22 or later installed on your system.
2. A Harness API key for authentication.

## Quickstart

### Build from source

1. Clone the repository:
```bash
git clone https://github.com/vistaarjuneja/harness-mcp.git
cd harness-mcp
```

2. Build the binary:
```bash
go build -o cmd/harness-mcp-server/harness-mcp-server ./cmd/harness-mcp-server
```

3. Run the server:
```bash
HARNESS_API_KEY=your_api_key HARNESS_ACCOUNT_ID=your_account_id HARNESS_ORG_ID=your_org_id HARNESS_PROJECT_ID=your_project_id ./cmd/harness-mcp-server/harness-mcp-server stdio
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
          "HARNESS_API_KEY": "<YOUR_API_KEY>",
          "HARNESS_ACCOUNT_ID": "<YOUR_ACCOUNT_ID>",
          "HARNESS_ORG_ID": "<YOUR_ORG_ID>",
          "HARNESS_PROJECT_ID": "<YOUR_PROJECT_ID>"
        }
      }
    }
  }
  ```
</details>

## Usage with Windsurf

To use the Harness MCP Server with Windsurf:

1. Add the server configuration to your Windsurf config file:

```json
{
  "mcpServers": {
    "harness": {
      "command": "/path/to/harness-mcp-server",
      "args": ["stdio"],
      "env": {
        "HARNESS_API_KEY": "<YOUR_API_KEY>",
        "HARNESS_ACCOUNT_ID": "<YOUR_ACCOUNT_ID>",
        "HARNESS_ORG_ID": "<YOUR_ORG_ID>",
        "HARNESS_PROJECT_ID": "<YOUR_PROJECT_ID>",
        "HARNESS_BASE_URL": "<YOUR_BASE_URL>",
      }
    }
  }
}
```

## Development

### Command Line Arguments

The Harness MCP Server supports the following command line arguments:

- `--toolsets`: Comma-separated list of tool groups to enable (default: "all")
- `--read-only`: Run the server in read-only mode
- `--log-file`: Path to log file for debugging
- `--log-level`: Set the logging level (debug, info, warn, error)
- `--version`: Show version information
- `--help`: Show help message
- `--base-url`: Base URL for Harness (default: "https://app.harness.io")

### Environment Variables

Environment variables are prefixed with `HARNESS_`:

- `HARNESS_API_KEY`: Harness API key (required)
- `HARNESS_ACCOUNT_ID`: Harness account ID (required)
- `HARNESS_ORG_ID`: Harness organization ID (optional, but required for some operations)
- `HARNESS_PROJECT_ID`: Harness project ID (optional, but required for some operations)
- `HARNESS_TOOLSETS`: Comma-separated list of toolsets to enable (default: "all")
- `HARNESS_READ_ONLY`: Set to "true" to run in read-only mode
- `HARNESS_LOG_FILE`: Path to log file
- `HARNESS_LOG_LEVEL`: Set the logging level (debug, info, warn, error)
- `HARNESS_BASE_URL`: Base URL for Harness (default: "https://app.harness.io")

### Authentication

The server uses a Harness API key for authentication. This can be set via the `HARNESS_API_KEY` environment variable.

## Debugging

Since MCP servers run over stdio, debugging can be challenging. For the best debugging experience, we strongly recommend using the [MCP Inspector](https://github.com/modelcontextprotocol/inspector).

You can launch the MCP Inspector with this command:

```bash
npx @modelcontextprotocol/inspector /path/to/harness-mcp-server stdio
```

Upon launching, the Inspector will display a URL that you can access in your browser to begin debugging.