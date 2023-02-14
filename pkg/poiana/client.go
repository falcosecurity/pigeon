package poiana

import (
	"context"
	"os"
	"strings"

	"github.com/google/go-github/v50/github"
)

func NewClient(ctx context.Context, tokenFile string) (*github.Client, error) {
	ghTokBytes, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}
	ghTok := string(ghTokBytes)
	ghTok = strings.Trim(ghTok, "\n")
	return github.NewTokenClient(ctx, ghTok), nil
}
