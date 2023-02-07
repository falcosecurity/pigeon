package poiana

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"os"

	"github.com/google/go-github/v50/github"
)

func NewClient(ctx context.Context, tokenFile string) (*github.Client, error) {
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
	return github.NewClient(tc), nil
}
