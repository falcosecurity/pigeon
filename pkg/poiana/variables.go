package poiana

import (
	"context"
	"github.com/google/go-github/v50/github"
)

type ActionsVarsService interface {
	ListRepoVariables(ctx context.Context, owner, repo string, opts *github.ListOptions) (*github.ActionsVariables, *github.Response, error)
	DeleteRepoVariable(ctx context.Context, owner, repo, name string) (*github.Response, error)
	CreateRepoVariable(ctx context.Context, owner, repo string, variable *github.ActionsVariable) (*github.Response, error)
	UpdateRepoVariable(ctx context.Context, owner, repo string, variable *github.ActionsVariable) (*github.Response, error)
}
