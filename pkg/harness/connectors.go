package harness

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Base URL for Harness API
const baseURL = "https://app.harness.io/ng/api"

// Default account, org, and project identifiers
// In a real implementation, these would be configurable
// TODO: remove these
const (
	defaultAccountID = "wFHXHD0RRQWoO8tIZT5YVw"
	defaultOrgID     = "default"
	defaultProjectID = "nitisha"
)

// ConnectorClient provides methods to interact with Harness API for connectors
type ConnectorClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewConnectorClient creates a new connector client
func NewConnectorClient() *ConnectorClient {
	return &ConnectorClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// ListConnectorRequest represents the request to list connectors
type ListConnectorRequest struct {
	Categories           []string `json:"categories,omitempty"`
	ConnectorNames       []string `json:"connectorNames,omitempty"`
	ConnectorIdentifiers []string `json:"connectorIdentifiers,omitempty"`
	Types                []string `json:"types,omitempty"`
	FilterType           string   `json:"filterType"`
}

// ListConnectors retrieves connectors from the Harness API
func (c *ConnectorClient) ListConnectors(ctx context.Context, connectorNames []string, connectorIDs []string, types []string) ([]byte, error) {
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
	params.Add("getDefaultFromOtherRepo", "true")
	params.Add("getDistinctFromBranches", "true")
	params.Add("onlyFavorites", "false")
	params.Add("pageIndex", "0")
	params.Add("pageSize", "50")
	params.Add("sortOrders", "orderType=ASC")

	// Create the request body
	reqBody := ListConnectorRequest{
		Categories:           []string{"CLOUD_PROVIDER"},
		ConnectorNames:       connectorNames,
		ConnectorIdentifiers: connectorIDs,
		Types:                types,
		FilterType:           "Connector",
	}

	// Marshal the request body to JSON
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create a new HTTP request
	url := fmt.Sprintf("%s/connectors/listV2?%s", c.baseURL, params.Encode())
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

// ListConnectorsTool creates the list-connector tool
func ListConnectorsTool(client *ConnectorClient) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("list-connector",
		mcp.WithDescription("List available connectors"),
		mcp.WithString("connector_type",
			mcp.Description("Type of connector to filter by")),
		mcp.WithArray("connector_names",
			mcp.Description("List of connector names to filter by"),
			mcp.Items(map[string]interface{}{"type": "string"})),
		mcp.WithArray("connector_ids",
			mcp.Description("List of connector IDs to filter by"),
			mcp.Items(map[string]interface{}{"type": "string"})),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract parameters
		connectorType, err := OptionalParam[string](request, "connector_type")
		if err != nil {
			return nil, err
		}

		connectorNames, err := OptionalStringArrayParam(request, "connector_names")
		if err != nil {
			return nil, err
		}

		connectorIDs, err := OptionalStringArrayParam(request, "connector_ids")
		if err != nil {
			return nil, err
		}

		// Prepare types array if connector_type is specified
		var types []string
		if connectorType != "" {
			types = []string{connectorType}
		}

		// Call Harness API
		jsonResponse, err := client.ListConnectors(ctx, connectorNames, connectorIDs, types)
		if err != nil {
			return nil, err
		}

		// Return the raw JSON response as a text result
		return mcp.NewToolResultText(string(jsonResponse)), nil
	}

	return tool, handler
}
