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

type ActionsSecretsService interface {
	GetPublicKey(ctx context.Context) (*github.PublicKey, *github.Response, error)
	ListSecrets(ctx context.Context, opts *github.ListOptions) (*github.Secrets, *github.Response, error)
	DeleteSecret(ctx context.Context, name string) (*github.Response, error)
	CreateOrUpdateSecret(ctx context.Context, eSecret *github.EncryptedSecret) (*github.Response, error)
}
