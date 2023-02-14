package pigeon

import (
	"context"

	"github.com/google/go-github/v50/github"
)

type ActionsVarsService interface {
	ListVariables(ctx context.Context, opts *github.ListOptions) (*github.ActionsVariables, *github.Response, error)
	DeleteVariable(ctx context.Context, name string) (*github.Response, error)
	CreateOrUpdateVariable(ctx context.Context, variable *github.ActionsVariable) (*github.Response, error)
}
