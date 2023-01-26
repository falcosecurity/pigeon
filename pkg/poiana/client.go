package poiana

import (
	"context"

	"github.com/google/go-github/v49/github"
)

// Variable represents a repository action variable.
type Variable struct {
	Name      string           `json:"name"`
	Value     string           `json:"value"`
	CreatedAt github.Timestamp `json:"created_at"`
	UpdatedAt github.Timestamp `json:"updated_at"`
}

type Variables struct {
	TotalCount int         `json:"total_count"`
	Variables  []*Variable `json:"variables"`
}

type actionsService struct {
	*github.ActionsService
	client *github.Client
}

type Client struct {
	*github.Client
	Actions *actionsService
}

func NewClient(c *github.Client) *Client {
	return &Client{
		Client: c,
		Actions: &actionsService{
			ActionsService: c.Actions,
			client:         c,
		},
	}
}

func (a *actionsService) GetPublicKey(ctx context.Context, orgName string, repoName string) (*github.PublicKey, error) {
	pKey, _, err := a.GetRepoPublicKey(ctx, orgName, repoName)
	return pKey, err
}
