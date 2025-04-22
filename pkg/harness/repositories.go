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

// CreateRepositoryRequest represents the request to create a repository
type CreateRepositoryRequest struct {
	DefaultBranch string `json:"default_branch,omitempty"`
	Description   string `json:"description,omitempty"`
	ForkID        int    `json:"fork_id,omitempty"`
	GitIgnore     string `json:"git_ignore,omitempty"`
	Identifier    string `json:"identifier"`
	IsPublic      bool   `json:"is_public"`
	License       string `json:"license,omitempty"`
	ParentRef     string `json:"parent_ref,omitempty"`
	Readme        bool   `json:"readme"`
	UID           string `json:"uid,omitempty"`
}

// CreateRepository creates a new repository in Harness
func (c *ConnectorClient) CreateRepository(ctx context.Context, request CreateRepositoryRequest) ([]byte, error) {
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
	url := fmt.Sprintf("%s/v1/repos?%s", c.config.CodeBaseURL, params.Encode())
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

// CreateRepositoryTool creates the create-repository tool
func CreateRepositoryTool(client *ConnectorClient) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("create-repository",
		mcp.WithDescription("Create a new repository in Harness"),
		mcp.WithString("identifier",
			mcp.Description("Unique identifier for the repository"),
			mcp.Required()),
		mcp.WithString("default_branch",
			mcp.Description("Default branch name")),
		mcp.WithString("description",
			mcp.Description("Repository description")),
		mcp.WithNumber("fork_id",
			mcp.Description("ID of the repository to fork from")),
		mcp.WithString("git_ignore",
			mcp.Description("Git ignore template")),
		mcp.WithBoolean("is_public",
			mcp.Description("Whether the repository is public"),
			mcp.Required()),
		mcp.WithString("license",
			mcp.Description("License template")),
		mcp.WithString("parent_ref",
			mcp.Description("Parent reference for the repository")),
		mcp.WithBoolean("readme",
			mcp.Description("Whether to initialize with a README"),
			mcp.Required()),
		mcp.WithString("uid",
			mcp.Description("Unique ID for the repository")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract required parameters
		identifier, err := requiredParam[string](request, "identifier")
		if err != nil {
			return nil, err
		}

		isPublic, err := requiredParam[bool](request, "is_public")
		if err != nil {
			return nil, err
		}

		readme, err := requiredParam[bool](request, "readme")
		if err != nil {
			return nil, err
		}

		// Extract optional parameters
		defaultBranch, err := OptionalParam[string](request, "default_branch")
		if err != nil {
			return nil, err
		}

		description, err := OptionalParam[string](request, "description")
		if err != nil {
			return nil, err
		}

		forkID, err := OptionalParam[int](request, "fork_id")
		if err != nil {
			return nil, err
		}

		gitIgnore, err := OptionalParam[string](request, "git_ignore")
		if err != nil {
			return nil, err
		}

		license, err := OptionalParam[string](request, "license")
		if err != nil {
			return nil, err
		}

		parentRef, err := OptionalParam[string](request, "parent_ref")
		if err != nil {
			return nil, err
		}

		uid, err := OptionalParam[string](request, "uid")
		if err != nil {
			return nil, err
		}

		// Create the request
		createRepoRequest := CreateRepositoryRequest{
			Identifier:    identifier,
			IsPublic:      isPublic,
			Readme:        readme,
			DefaultBranch: defaultBranch,
			Description:   description,
			ForkID:        forkID,
			GitIgnore:     gitIgnore,
			License:       license,
			ParentRef:     parentRef,
			UID:           uid,
		}

		// Call Harness API
		jsonResponse, err := client.CreateRepository(ctx, createRepoRequest)
		if err != nil {
			return nil, err
		}

		// Return the raw JSON response as a text result
		return mcp.NewToolResultText(string(jsonResponse)), nil
	}

	return tool, handler
}

// GetRepository retrieves a repository by its identifier
func (c *ConnectorClient) GetRepository(ctx context.Context, repoIdentifier string) ([]byte, error) {
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
	url := fmt.Sprintf("%s/v1/repos/%s?%s", c.config.CodeBaseURL, repoIdentifier, params.Encode())
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

// GetRepositoryTool creates the get-repository tool
func GetRepositoryTool(client *ConnectorClient) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("get-repository",
		mcp.WithDescription("Get repository information by identifier"),
		mcp.WithString("repo_identifier",
			mcp.Description("The identifier of the repository to retrieve"),
			mcp.Required()),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract required parameters
		repoIdentifier, err := requiredParam[string](request, "repo_identifier")
		if err != nil {
			return nil, err
		}

		// Call Harness API
		jsonResponse, err := client.GetRepository(ctx, repoIdentifier)
		if err != nil {
			return nil, err
		}

		// Return the raw JSON response as a text result
		return mcp.NewToolResultText(string(jsonResponse)), nil
	}

	return tool, handler
}

// ListCommits lists commits for a repository
func (c *ConnectorClient) ListCommits(ctx context.Context, repoIdentifier string, params url.Values) ([]byte, error) {
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
	url := fmt.Sprintf("%s/v1/repos/%s/commits?%s", c.config.CodeBaseURL, repoIdentifier, params.Encode())
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

// GetCommit retrieves a specific commit
func (c *ConnectorClient) GetCommit(ctx context.Context, repoIdentifier, commitSHA string) ([]byte, error) {
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
	url := fmt.Sprintf("%s/v1/repos/%s/commits/%s?%s", c.config.CodeBaseURL, repoIdentifier, commitSHA, params.Encode())
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

// ListCommitsTool creates the list-commits tool
func ListCommitsTool(client *ConnectorClient) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("list-commits",
		mcp.WithDescription("List commits for a repository"),
		mcp.WithString("repo_identifier",
			mcp.Description("The identifier of the repository"),
			mcp.Required()),
		mcp.WithString("git_ref",
			mcp.Description("Git reference (branch, tag, etc.)")),
		mcp.WithString("after",
			mcp.Description("Commit SHA to start listing from")),
		mcp.WithString("path",
			mcp.Description("Filter commits by file path")),
		mcp.WithNumber("since",
			mcp.Description("Timestamp to filter commits created after")),
		mcp.WithNumber("until",
			mcp.Description("Timestamp to filter commits created before")),
		mcp.WithString("committer",
			mcp.Description("Filter by committer")),
		mcp.WithNumber("page",
			mcp.Description("Page number")),
		mcp.WithNumber("limit",
			mcp.Description("Number of items per page")),
		mcp.WithBoolean("include_stats",
			mcp.Description("Whether to include commit stats")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract required parameter
		repoIdentifier, err := requiredParam[string](request, "repo_identifier")
		if err != nil {
			return nil, err
		}

		// Build query parameters from optional parameters
		params := url.Values{}
		
		// Add optional string parameters if provided
		optionalStringParams := []string{
			"git_ref", "after", "path", "committer",
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
			"since", "until", "page", "limit",
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
		
		// Add optional boolean parameters
		includeStats, err := OptionalParam[bool](request, "include_stats")
		if err != nil {
			return nil, err
		}
		if includeStats {
			params.Add("include_stats", "true")
		}

		// Call Harness API
		jsonResponse, err := client.ListCommits(ctx, repoIdentifier, params)
		if err != nil {
			return nil, err
		}

		// Return the raw JSON response as a text result
		return mcp.NewToolResultText(string(jsonResponse)), nil
	}

	return tool, handler
}

// GetCommitTool creates the get-commit tool
func GetCommitTool(client *ConnectorClient) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("get-commit",
		mcp.WithDescription("Get a specific commit"),
		mcp.WithString("repo_identifier",
			mcp.Description("The identifier of the repository"),
			mcp.Required()),
		mcp.WithString("commit_sha",
			mcp.Description("The commit SHA"),
			mcp.Required()),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract required parameters
		repoIdentifier, err := requiredParam[string](request, "repo_identifier")
		if err != nil {
			return nil, err
		}

		commitSHA, err := requiredParam[string](request, "commit_sha")
		if err != nil {
			return nil, err
		}

		// Call Harness API
		jsonResponse, err := client.GetCommit(ctx, repoIdentifier, commitSHA)
		if err != nil {
			return nil, err
		}

		// Return the raw JSON response as a text result
		return mcp.NewToolResultText(string(jsonResponse)), nil
	}

	return tool, handler
}
