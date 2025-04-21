package harness

import (
	"context"

	"github.com/harness/harness-mcp/pkg/toolsets"
)

// Default tools to enable
var DefaultTools = []string{"all"}

// InitToolsets initializes and returns the toolset groups
func InitToolsets(passedToolsets []string, readOnly bool) (*toolsets.ToolsetGroup, error) {
	// Create connector client
	connectorClient := NewConnectorClient()

	// Create a toolset group
	tsg := toolsets.NewToolsetGroup(readOnly)

	// Create the connectors toolset
	connectors := toolsets.NewToolset("connectors", "Harness Connector related tools").
		AddReadTools(
			toolsets.NewServerTool(ListConnectorsTool(connectorClient)),
		)

	// Add toolsets to the group
	tsg.AddToolset(connectors)

	// Enable requested toolsets
	if err := tsg.EnableToolsets(passedToolsets); err != nil {
		return nil, err
	}

	return tsg, nil
}

// SetupContextWithApiKey sets up the context with the API key
func SetupContextWithApiKey(ctx context.Context, apiKey string) context.Context {
	return context.WithValue(ctx, "apiKey", apiKey)
}
