package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/harness/harness-mcp/client/dto"
)

const (
	pullRequestBasePath = "code/api/v1/repos"
	pullRequestGetPath  = pullRequestBasePath + "/%s/pullreq/%d"
	pullRequestListPath = pullRequestBasePath + "/%s/pullreq"
)

type PullRequestService struct {
	client *Client
}

func (p *PullRequestService) Get(ctx context.Context, scope dto.Scope, repoID string, prNumber int) (*dto.PullRequest, error) {
	path := fmt.Sprintf(pullRequestGetPath, repoID, prNumber)
	params := make(map[string]string)
	addScope(scope, params)

	pr := new(dto.PullRequest)
	err := p.client.Get(ctx, path, params, nil, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	return pr, nil
}

func (p *PullRequestService) List(ctx context.Context, scope dto.Scope, repoID string, opts *dto.PullRequestOptions) ([]*dto.PullRequest, error) {
	path := fmt.Sprintf(pullRequestListPath, repoID)
	params := make(map[string]string)
	addScope(scope, params)

	// Add query parameters from options
	if opts != nil {
		if len(opts.State) > 0 {
			params["state"] = strings.Join(opts.State, ",")
		}
		if opts.SourceRepoRef != "" {
			params["source_repo_ref"] = opts.SourceRepoRef
		}
		if opts.SourceBranch != "" {
			params["source_branch"] = opts.SourceBranch
		}
		if opts.TargetBranch != "" {
			params["target_branch"] = opts.TargetBranch
		}
		if opts.Query != "" {
			params["query"] = opts.Query
		}
		if len(opts.CreatedBy) > 0 {
			createdByStrings := make([]string, len(opts.CreatedBy))
			for i, id := range opts.CreatedBy {
				createdByStrings[i] = fmt.Sprintf("%d", id)
			}
			params["created_by"] = strings.Join(createdByStrings, ",")
		}
		if opts.Order != "" {
			params["order"] = opts.Order
		}
		if opts.Sort != "" {
			params["sort"] = opts.Sort
		}
		if opts.CreatedLt > 0 {
			params["created_lt"] = fmt.Sprintf("%d", opts.CreatedLt)
		}
		if opts.CreatedGt > 0 {
			params["created_gt"] = fmt.Sprintf("%d", opts.CreatedGt)
		}
		if opts.UpdatedLt > 0 {
			params["updated_lt"] = fmt.Sprintf("%d", opts.UpdatedLt)
		}
		if opts.UpdatedGt > 0 {
			params["updated_gt"] = fmt.Sprintf("%d", opts.UpdatedGt)
		}
		if opts.Page > 0 {
			params["page"] = fmt.Sprintf("%d", opts.Page)
		}
		if opts.Limit > 0 {
			params["limit"] = fmt.Sprintf("%d", opts.Limit)
		}
		if opts.AuthorID > 0 {
			params["author_id"] = fmt.Sprintf("%d", opts.AuthorID)
		}
		if opts.IncludeChecks {
			params["include_checks"] = "true"
		}
	}

	var prs []*dto.PullRequest
	err := p.client.Get(ctx, path, params, nil, &prs)
	if err != nil {
		return nil, fmt.Errorf("failed to list pull requests: %w", err)
	}

	return prs, nil
}
