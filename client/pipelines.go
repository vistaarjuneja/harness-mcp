package client

import (
	"context"
	"fmt"

	"github.com/harness/harness-mcp/client/dto"
)

const (
	pipelinePath          = "pipeline/api/pipelines/%s"
	pipelineListPath      = "pipeline/api/pipelines/list"
	pipelineExecutionPath = "pipeline/api/pipelines/execution/url"
)

type PipelineService struct {
	client *Client
}

func (p *PipelineService) Get(ctx context.Context, scope dto.Scope, pipelineID string) (*dto.Entity[dto.PipelineData], error) {
	path := fmt.Sprintf(pipelinePath, pipelineID)

	// Prepare query parameters
	params := make(map[string]string)
	addScope(scope, params)

	// Initialize the response object
	response := &dto.Entity[dto.PipelineData]{}

	// Make the GET request
	err := p.client.Get(ctx, path, params, map[string]string{}, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *PipelineService) List(ctx context.Context, scope dto.Scope, opts *dto.PipelineListOptions) (*dto.ListOutput[dto.PipelineListItem], error) {
	// Prepare query parameters
	params := make(map[string]string)
	addScope(scope, params)

	// Set default pagination and add pagination parameters if opts is provided
	if opts != nil {
		setDefaultPagination(&opts.PaginationOptions)
		params["page"] = fmt.Sprintf("%d", opts.Page)
		params["size"] = fmt.Sprintf("%d", opts.Size)

		// Add optional parameters if provided
		if opts.SearchTerm != "" {
			params["searchTerm"] = opts.SearchTerm
		}
	}

	// Create request body - this is required
	requestBody := map[string]string{
		"filterType": "PipelineSetup",
	}

	// Initialize the response object
	response := &dto.ListOutput[dto.PipelineListItem]{}

	// Make the POST request
	err := p.client.Post(ctx, pipelineListPath, params, requestBody, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *PipelineService) ListExecutions(ctx context.Context, scope dto.Scope, opts *dto.PipelineExecutionOptions) (*dto.ListOutput[dto.PipelineExecution], error) {
	setDefaultPagination(&opts.PaginationOptions)
	return nil, nil
}

func (p *PipelineService) FetchExecutionURL(ctx context.Context, scope dto.Scope, pipelineID, planExecutionID string) (string, error) {
	path := pipelineExecutionPath

	// Prepare query parameters
	params := make(map[string]string)
	addScope(scope, params)
	params["pipelineIdentifier"] = pipelineID
	params["planExecutionId"] = planExecutionID

	// Initialize the response object
	urlResponse := &dto.Entity[string]{}

	// Make the GET request
	err := p.client.Get(ctx, path, params, nil, urlResponse)
	if err != nil {
		return "", fmt.Errorf("failed to fetch execution URL: %w", err)
	}

	return urlResponse.Data, nil
}
