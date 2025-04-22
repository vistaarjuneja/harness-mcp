package harness

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// PullRequestLabel represents a label in a pull request
type PullRequestLabel struct {
	LabelID int    `json:"label_id,omitempty"`
	Value   string `json:"value,omitempty"`
	ValueID int    `json:"value_id,omitempty"`
}

// CreatePullRequestRequest represents the request to create a pull request
type CreatePullRequestRequest struct {
	BypassRules   bool              `json:"bypass_rules,omitempty"`
	Description   string            `json:"description,omitempty"`
	IsDraft       bool              `json:"is_draft,omitempty"`
	Labels        []PullRequestLabel `json:"labels,omitempty"`
	ReviewerIDs   []int             `json:"reviewer_ids,omitempty"`
	SourceBranch  string            `json:"source_branch"`
	SourceRepoRef string            `json:"source_repo_ref,omitempty"`
	TargetBranch  string            `json:"target_branch"`
	Title         string            `json:"title"`
}

// ListPullRequests lists all pull requests for a repository
func (c *ConnectorClient) ListPullRequests(ctx context.Context, repoIdentifier string, params url.Values) ([]byte, error) {
	// Get API key from context
	apiKey, err := GetApiKeyFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Add required identifiers if not present
	if params == nil {
		params = url.Values{}
	}
	if params.Get("accountIdentifier") == "" {
		params.Add("accountIdentifier", c.config.AccountID)
	}
	if params.Get("orgIdentifier") == "" {
		params.Add("orgIdentifier", c.config.OrgID)
	}
	if params.Get("projectIdentifier") == "" {
		params.Add("projectIdentifier", c.config.ProjectID)
	}

	// Create a new HTTP request
	url := fmt.Sprintf("%s/v1/repos/%s/pullreq?%s", c.config.CodeBaseURL, repoIdentifier, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Make the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if the response status code is not 2xx
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// Return the raw JSON response
	return body, nil
}

// GetPullRequest retrieves a specific pull request
func (c *ConnectorClient) GetPullRequest(ctx context.Context, repoIdentifier string, pullreqNumber int) ([]byte, error) {
	// Get API key from context
	apiKey, err := GetApiKeyFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Construct the query parameters
	params := url.Values{}
	params.Add("accountIdentifier", c.config.AccountID)
	params.Add("orgIdentifier", c.config.OrgID)
	params.Add("projectIdentifier", c.config.ProjectID)

	// Create a new HTTP request
	url := fmt.Sprintf("%s/v1/repos/%s/pullreq/%d?%s", c.config.CodeBaseURL, repoIdentifier, pullreqNumber, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Make the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if the response status code is not 2xx
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// Return the raw JSON response
	return body, nil
}

// CreatePullRequest creates a new pull request
func (c *ConnectorClient) CreatePullRequest(ctx context.Context, repoIdentifier string, request CreatePullRequestRequest) ([]byte, error) {
	// Get API key from context
	apiKey, err := GetApiKeyFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Construct the query parameters
	params := url.Values{}
	params.Add("accountIdentifier", c.config.AccountID)
	params.Add("orgIdentifier", c.config.OrgID)
	params.Add("projectIdentifier", c.config.ProjectID)

	// Marshal the request body to JSON
	reqBodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create a new HTTP request
	url := fmt.Sprintf("%s/v1/repos/%s/pullreq?%s", c.config.CodeBaseURL, repoIdentifier, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Make the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if the response status code is not 2xx
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// Return the raw JSON response
	return body, nil
}

// ListPullRequestsTool creates the list-pullreq tool
func ListPullRequestsTool(client *ConnectorClient) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("list-pullreq",
		mcp.WithDescription("List pull requests for a repository"),
		mcp.WithString("repo_identifier",
			mcp.Description("The identifier of the repository"),
			mcp.Required()),
		mcp.WithString("state",
			mcp.Description("Filter by state (open, closed, all)")),
		mcp.WithString("source_repo_ref",
			mcp.Description("Filter by source repository reference")),
		mcp.WithString("source_branch",
			mcp.Description("Filter by source branch")),
		mcp.WithString("target_branch",
			mcp.Description("Filter by target branch")),
		mcp.WithString("query",
			mcp.Description("Search query")),
		mcp.WithNumber("created_by",
			mcp.Description("Filter by creator user ID")),
		mcp.WithString("order",
			mcp.Description("Order of results (asc, desc)")),
		mcp.WithString("sort",
			mcp.Description("Sort field (created, updated)")),
		mcp.WithNumber("page",
			mcp.Description("Page number")),
		mcp.WithNumber("limit",
			mcp.Description("Number of items per page")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract required parameter
		repoIdentifier, err := requiredParam[string](request, "repo_identifier")
		if err != nil {
			return nil, err
		}

		// Build query parameters from optional parameters
		params := url.Values{}
		
		// Add optional parameters if provided
		optionalStringParams := []string{
			"state", "source_repo_ref", "source_branch", "target_branch", 
			"query", "order", "sort",
		}
		
		for _, param := range optionalStringParams {
			value, err := OptionalParam[string](request, param)
			if err != nil {
				return nil, err
			}
			if value != "" {
				params.Add(param, value)
			}
		}
		
		// Add optional number parameters
		optionalNumberParams := []string{
			"created_by", "page", "limit",
		}
		
		for _, param := range optionalNumberParams {
			value, err := OptionalParam[float64](request, param)
			if err != nil {
				return nil, err
			}
			if value != 0 {
				params.Add(param, fmt.Sprintf("%v", value))
			}
		}

		// Call Harness API
		jsonResponse, err := client.ListPullRequests(ctx, repoIdentifier, params)
		if err != nil {
			return nil, err
		}

		// Return the raw JSON response as a text result
		return mcp.NewToolResultText(string(jsonResponse)), nil
	}

	return tool, handler
}

// GetPullRequestTool creates the get-pullreq tool
func GetPullRequestTool(client *ConnectorClient) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("get-pullreq",
		mcp.WithDescription("Get a specific pull request"),
		mcp.WithString("repo_identifier",
			mcp.Description("The identifier of the repository"),
			mcp.Required()),
		mcp.WithNumber("pullreq_number",
			mcp.Description("The pull request number"),
			mcp.Required()),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract required parameters
		repoIdentifier, err := requiredParam[string](request, "repo_identifier")
		if err != nil {
			return nil, err
		}

		pullreqNumberFloat, err := requiredParam[float64](request, "pullreq_number")
		if err != nil {
			return nil, err
		}
		pullreqNumber := int(pullreqNumberFloat)

		// Call Harness API
		jsonResponse, err := client.GetPullRequest(ctx, repoIdentifier, pullreqNumber)
		if err != nil {
			return nil, err
		}

		// Return the raw JSON response as a text result
		return mcp.NewToolResultText(string(jsonResponse)), nil
	}

	return tool, handler
}

// CreatePullRequestTool creates the create-pullreq tool
func CreatePullRequestTool(client *ConnectorClient) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("create-pullreq",
		mcp.WithDescription("Create a new pull request"),
		mcp.WithString("repo_identifier",
			mcp.Description("The identifier of the repository"),
			mcp.Required()),
		mcp.WithString("title",
			mcp.Description("Title of the pull request"),
			mcp.Required()),
		mcp.WithString("source_branch",
			mcp.Description("Source branch name"),
			mcp.Required()),
		mcp.WithString("target_branch",
			mcp.Description("Target branch name"),
			mcp.Required()),
		mcp.WithString("description",
			mcp.Description("Description of the pull request")),
		mcp.WithBoolean("is_draft",
			mcp.Description("Whether the pull request is a draft")),
		mcp.WithBoolean("bypass_rules",
			mcp.Description("Whether to bypass rules")),
		mcp.WithString("source_repo_ref",
			mcp.Description("Source repository reference for cross-repo PRs")),
		mcp.WithArray("reviewer_ids",
			mcp.Description("IDs of users to be added as reviewers"),
			mcp.Items(map[string]interface{}{"type": "number"})),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract required parameters
		repoIdentifier, err := requiredParam[string](request, "repo_identifier")
		if err != nil {
			return nil, err
		}

		title, err := requiredParam[string](request, "title")
		if err != nil {
			return nil, err
		}

		sourceBranch, err := requiredParam[string](request, "source_branch")
		if err != nil {
			return nil, err
		}

		targetBranch, err := requiredParam[string](request, "target_branch")
		if err != nil {
			return nil, err
		}

		// Extract optional parameters
		description, err := OptionalParam[string](request, "description")
		if err != nil {
			return nil, err
		}

		isDraft, err := OptionalParam[bool](request, "is_draft")
		if err != nil {
			return nil, err
		}

		bypassRules, err := OptionalParam[bool](request, "bypass_rules")
		if err != nil {
			return nil, err
		}

		sourceRepoRef, err := OptionalParam[string](request, "source_repo_ref")
		if err != nil {
			return nil, err
		}

		// Get reviewer IDs
		reviewerIDsFloat, err := OptionalParam[[]interface{}](request, "reviewer_ids")
		if err != nil {
			return nil, err
		}
		
		// Convert reviewer IDs from float64 to int
		reviewerIDs := []int{}
		if len(reviewerIDsFloat) > 0 {
			for _, id := range reviewerIDsFloat {
				if floatID, ok := id.(float64); ok {
					reviewerIDs = append(reviewerIDs, int(floatID))
				} else {
					return nil, fmt.Errorf("reviewer_ids must be numbers")
				}
			}
		}

		// Create the request
		createPRRequest := CreatePullRequestRequest{
			Title:         title,
			SourceBranch:  sourceBranch,
			TargetBranch:  targetBranch,
			Description:   description,
			IsDraft:       isDraft,
			BypassRules:   bypassRules,
			SourceRepoRef: sourceRepoRef,
			ReviewerIDs:   reviewerIDs,
		}

		// Call Harness API
		jsonResponse, err := client.CreatePullRequest(ctx, repoIdentifier, createPRRequest)
		if err != nil {
			return nil, err
		}

		// Return the raw JSON response as a text result
		return mcp.NewToolResultText(string(jsonResponse)), nil
	}

	return tool, handler
}
