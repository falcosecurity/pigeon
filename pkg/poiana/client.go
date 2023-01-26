package poiana

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"os"

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

func (a *actionsService) GetPublicKey(ctx context.Context, orgName string, repoName string) (*github.PublicKey, error) {
	pKey, _, err := a.GetRepoPublicKey(ctx, orgName, repoName)
	return pKey, err
}

func NewClient(ctx context.Context, tokenFile string) (*Client, error) {
	ghTokBytes, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}
	ghTok := string(ghTokBytes)
	var tc *http.Client
	if ghTok != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: ghTok},
		)
		tc = oauth2.NewClient(ctx, ts)
	}
	ghCl := github.NewClient(tc)
	return &Client{
		Client: ghCl,
		Actions: &actionsService{
			ActionsService: ghCl.Actions,
			client:         ghCl,
		},
	}, nil
}
