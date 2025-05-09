package harness

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/harness/harness-mcp/client"
	"github.com/harness/harness-mcp/client/dto"
	"github.com/harness/harness-mcp/cmd/harness-mcp-server/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetPipelineTool(config *config.Config, client *client.Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_pipeline",
			mcp.WithDescription("Get details of a specific pipeline in a Harness repository."),
			mcp.WithString("pipeline_id",
				mcp.Required(),
				mcp.Description("The ID of the pipeline"),
			),
			WithScope(config, true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			pipelineID, err := requiredParam[string](request, "pipeline_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			scope, err := fetchScope(config, request, true)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			data, err := client.Pipelines.Get(ctx, scope, pipelineID)
			if err != nil {
				return nil, fmt.Errorf("failed to get pipeline: %w", err)
			}

			r, err := json.Marshal(data.Data)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal pipeline: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

func ListPipelinesTool(config *config.Config, client *client.Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_pipelines",
			mcp.WithDescription("List pipelines in a Harness repository."),
			mcp.WithString("search_term",
				mcp.Description("Optional search term to filter pipelines"),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination"),
			),
			mcp.WithNumber("size",
				mcp.Description("Number of items per page"),
			),
			WithScope(config, true),
			WithPagination(),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			scope, err := fetchScope(config, request, true)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			page, size, err := fetchPagination(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			searchTerm, err := OptionalParam[string](request, "search_term")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			opts := &dto.PipelineListOptions{
				SearchTerm: searchTerm,
				PaginationOptions: dto.PaginationOptions{
					Page: page,
					Size: size,
				},
			}

			data, err := client.Pipelines.List(ctx, scope, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list pipelines: %w", err)
			}

			r, err := json.Marshal(data)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal pipeline list: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
