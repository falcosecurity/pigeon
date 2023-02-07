package poiana

import (
	"context"

	"github.com/google/go-github/v50/github"
)

// SecretsProvider retrieves secrets with a given key
type SecretsProvider interface {
	// GetSecret returns a secret with the given key.
	// Returns a non-nil error in case of failure
	GetSecret(string) (string, error)
}

type PublicKeyProvider interface {
	GetRepoPublicKey(ctx context.Context, orgName string, repoName string) (*github.PublicKey, *github.Response, error)
}

type ActionsSecretsService interface {
	ListRepoSecrets(ctx context.Context, owner, repo string, opts *github.ListOptions) (*github.Secrets, *github.Response, error)
	DeleteRepoSecret(ctx context.Context, owner, repo, name string) (*github.Response, error)
	CreateOrUpdateRepoSecret(ctx context.Context, owner, repo string, eSecret *github.EncryptedSecret) (*github.Response, error)
}
