package harness

import (
	"fmt"

	"github.com/harness/harness-mcp/client/dto"
	"github.com/harness/harness-mcp/cmd/harness-mcp-server/config"
	"github.com/mark3labs/mcp-go/mcp"
)

// WithScope adds org_id and project_id as optional parameters if they are not already defined in the
// config.
func WithScope(config *config.Config, required bool) mcp.ToolOption {
	var opt mcp.PropertyOption
	if required {
		opt = mcp.Required()
	}
	return func(tool *mcp.Tool) {
		if config.OrgID == "" {
			mcp.WithString("org_id",
				mcp.Description("The ID of the organization."),
				opt,
			)
		}
		if config.ProjectID == "" {
			mcp.WithString("project_id",
				mcp.Description("The ID of the project."),
				opt,
			)
		}
	}
}

// fetchScope fetches the scope from the config and MCP request.
// It looks in the config first and then in the request (if it was defined).
// If orgID and projectID are required fields, it will return an error if they are not present.
func fetchScope(config *config.Config, request mcp.CallToolRequest, required bool) (dto.Scope, error) {
	scope := dto.Scope{}
	// account ID is always required
	if config.AccountID == "" {
		return scope, fmt.Errorf("account ID is required")
	}
	scope.AccountID = config.AccountID
	// org ID and project ID may or may not be required for APIs. If they are required, we return an error
	// if not present.
	scope.OrgID = config.OrgID
	scope.ProjectID = config.ProjectID

	if scope.OrgID == "" {
		// try to fetch it from the MCP request
		org, _ := OptionalParam[string](request, "org_id")
		scope.OrgID = org
	}
	if scope.ProjectID == "" {
		// try to fetch it from the MCP request
		project, _ := OptionalParam[string](request, "project_id")
		scope.ProjectID = project
	}

	if required {
		if scope.OrgID == "" || scope.ProjectID == "" {
			return scope, fmt.Errorf("org ID and project ID are required")
		}
	}

	return scope, nil
}
