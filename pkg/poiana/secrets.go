package poiana

import (
	"context"
	"github.com/google/go-github/v49/github"
)

// SecretsProvider retrieves secrets with a given key
type SecretsProvider interface {
	// GetSecret returns a secret with the given key.
	// Returns a non-nil error in case of failure
	GetSecret(string) (string, error)
}

type ActionsSecretsService interface {
	ListRepoSecrets(ctx context.Context, owner, repo string) (*github.Secrets, error)
	DeleteRepoSecret(ctx context.Context, owner, repo, name string) error
	CreateOrUpdateRepoSecret(ctx context.Context, owner, repo string, eSecret *github.EncryptedSecret) error
}

func (s *actionsService) ListRepoSecrets(ctx context.Context, owner, repo string) (*github.Secrets, error) {
	secrets, _, err := s.ActionsService.ListRepoSecrets(ctx, owner, repo, nil)
	return secrets, err
}

func (s *actionsService) DeleteRepoSecret(ctx context.Context, owner, repo, name string) error {
	_, err := s.ActionsService.DeleteRepoSecret(ctx, owner, repo, name)
	return err
}

func (s *actionsService) CreateOrUpdateRepoSecret(ctx context.Context, owner, repo string, eSecret *github.EncryptedSecret) error {
	_, err := s.ActionsService.CreateOrUpdateRepoSecret(ctx, owner, repo, eSecret)
	return err
}
