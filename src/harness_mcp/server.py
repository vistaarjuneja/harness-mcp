import asyncio

from mcp.server.models import InitializationOptions
import mcp.types as types
from mcp.server import NotificationOptions, Server
from pydantic import AnyUrl
import json
from .connectors import list_connectors
import mcp.server.stdio

server = Server("harness-mcp")

@server.list_tools()
async def handle_list_tools() -> list[types.Tool]:
    """
    List available tools.
    Each tool specifies its arguments using JSON Schema validation.
    """
    return [
        types.Tool(
            name="list-connector",
            description="List available connectors",
            inputSchema={
                "type": "object",
                "properties": {
                    "connector_type": {"type": "string"},
                    "connector_names": {"type": "array", "items": {"type": "string"}},
                    "connector_ids": {"type": "array", "items": {"type": "string"}},
                },
                "required": [],
            },
        )
    ]

@server.call_tool()
async def handle_call_tool(
    name: str, arguments: dict | None
) -> list[types.TextContent | types.ImageContent | types.EmbeddedResource]:
    """
    Handle tool execution requests.
    Tools can modify server state and notify clients of changes.
    """
    if name == "list-connector":
        connector_type = arguments.get("connector_type") if arguments else None
        connector_names = arguments.get("connector_names") if arguments else None
        connector_ids = arguments.get("connector_ids") if arguments else None
        
        response = list_connectors(
            connector_names=connector_names,
            connector_identifiers=connector_ids,
            types=[connector_type] if connector_type else None
        )
        
        return [
            types.TextContent(
                type="text",
                text=json.dumps(response, indent=2),
            )
        ]
    else:
        raise ValueError(f"Unknown tool: {name}")

async def main():
    # Run the server using stdin/stdout streams
    async with mcp.server.stdio.stdio_server() as (read_stream, write_stream):
        await server.run(
            read_stream,
            write_stream,
            InitializationOptions(
                server_name="harness-mcp",
                server_version="0.1.0",
                capabilities=server.get_capabilities(
                    notification_options=NotificationOptions(),
                    experimental_capabilities={},
                ),
            ),
        )