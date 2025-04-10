# harness-mcp MCP server

MCP server for Harness

## Components

### Tools

The server implements one tool:
- list-connector: List available connectors
  - Takes "connector_type", "connector_names", and "connector_ids" as optional arguments
  - Returns a list of connectors in JSON format

## To do

We need to figure out a good way to handle auth, it's just hardcoded for this POC.

## Quickstart

### Install

#### Virtual Environment Setup

1. Create a virtual environment (already done):
```bash
python -m venv venv
```

2. Activate the virtual environment:
```bash
# On macOS/Linux
source venv/bin/activate

# On Windows
venv\Scripts\activate
```

3. Install dependencies:
```bash
pip install -r requirements.txt
```

4. Configure your API key by updating the `.env` file with your Harness API key.

#### Claude Desktop

On MacOS: `~/Library/Application\ Support/Claude/claude_desktop_config.json`
On Windows: `%APPDATA%/Claude/claude_desktop_config.json`

<details>
  <summary>Development/Unpublished Servers Configuration</summary>
  ```
  "mcpServers": {
    "harness-mcp": {
      "command": "uv",
      "args": [
        "--directory",
        "<path-to-harness-mcp>",
        "run",
        "harness-mcp"
      ]
    }
  }
  ```
</details>

<details>
  <summary>Published Servers Configuration</summary>
  ```
  "mcpServers": {
    "harness-mcp": {
      "command": "uvx",
      "args": [
        "harness-mcp"
      ]
    }
  }
  ```
</details>

## Development

### Building and Publishing

To prepare the package for distribution:

1. Sync dependencies and update lockfile:
```bash
uv sync
```

2. Build package distributions:
```bash
uv build
```

This will create source and wheel distributions in the `dist/` directory.

3. Publish to PyPI:
```bash
uv publish
```

Note: You'll need to set PyPI credentials via environment variables or command flags:
- Token: `--token` or `UV_PUBLISH_TOKEN`
- Or username/password: `--username`/`UV_PUBLISH_USERNAME` and `--password`/`UV_PUBLISH_PASSWORD`

### Debugging

Since MCP servers run over stdio, debugging can be challenging. For the best debugging
experience, we strongly recommend using the [MCP Inspector](https://github.com/modelcontextprotocol/inspector).


You can launch the MCP Inspector via [`npm`](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm) with this command:

```bash
npx @modelcontextprotocol/inspector uv --directory <path-to-harness-mcp> run harness-mcp
```


Upon launching, the Inspector will display a URL that you can access in your browser to begin debugging.