package poiana

import (
	"context"
	"github.com/google/go-github/v50/github"
)

type ActionsRepoVarsService interface {
	ListRepoVariables(ctx context.Context, owner, repo string, opts *github.ListOptions) (*github.ActionsVariables, *github.Response, error)
	DeleteRepoVariable(ctx context.Context, owner, repo, name string) (*github.Response, error)
	CreateRepoVariable(ctx context.Context, owner, repo string, variable *github.ActionsVariable) (*github.Response, error)
	UpdateRepoVariable(ctx context.Context, owner, repo string, variable *github.ActionsVariable) (*github.Response, error)
}

type ActionsOrgVarsService interface {
	ListOrgVariables(ctx context.Context, owner string, opts *github.ListOptions) (*github.ActionsVariables, *github.Response, error)
	DeleteOrgVariable(ctx context.Context, owner, name string) (*github.Response, error)
	CreateOrgVariable(ctx context.Context, owner string, variable *github.ActionsVariable) (*github.Response, error)
	UpdateOrgVariable(ctx context.Context, owner string, variable *github.ActionsVariable) (*github.Response, error)
}

type ActionsVarsService interface {
	ActionsRepoVarsService
	ActionsOrgVarsService
}
