package harness

import (
	"context"
	"fmt"

	"github.com/harness/harness-mcp/client"
	"github.com/harness/harness-mcp/cmd/harness-mcp-server/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// FetchLogDownloadURLTool creates a tool for fetching log download URLs for pipeline executions
func FetchLogDownloadURLTool(config *config.Config, client *client.Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("fetch_log_download_url",
			mcp.WithDescription("Fetch the download URL for pipeline execution logs in Harness."),
			mcp.WithString("plan_execution_id",
				mcp.Required(),
				mcp.Description("The ID of the plan execution"),
			),
			WithScope(config, true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			planExecutionID, err := requiredParam[string](request, "plan_execution_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			scope, err := fetchScope(config, request, true)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			logDownloadURL, err := client.Logs.DownloadLogs(ctx, scope, planExecutionID)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch log download URL: %w", err)
			}

			instruction := fmt.Sprintf("You can use the command \"curl -o logs.zip %s\" to download the zip for the logs.", logDownloadURL)

			return mcp.NewToolResultText(instruction), nil
		}
}
