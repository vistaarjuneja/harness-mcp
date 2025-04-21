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
	params.Add("accountIdentifier", defaultAccountID)
	params.Add("orgIdentifier", defaultOrgID)
	params.Add("projectIdentifier", defaultProjectID)

	// Marshal the request body to JSON
	reqBodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create a new HTTP request
	url := fmt.Sprintf("https://app.harness.io/code/api/v1/repos?%s", params.Encode())
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
	params.Add("accountIdentifier", defaultAccountID)
	params.Add("orgIdentifier", defaultOrgID)
	params.Add("projectIdentifier", defaultProjectID)

	// Create a new HTTP request
	url := fmt.Sprintf("https://app.harness.io/code/api/v1/repos/%s?%s", repoIdentifier, params.Encode())
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
