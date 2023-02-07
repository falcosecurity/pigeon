package poiana

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/v50/github"
)

type mockSecretsProvider struct {
	m map[string]string
}

// NewMockSecretsProvider creates a secret provider that retrieves
// secrets from a key-value map.
func NewMockSecretsProvider(values map[string]string) (SecretsProvider, error) {
	return &mockSecretsProvider{m: values}, nil
}

func (m *mockSecretsProvider) GetSecret(key string) (string, error) {
	v, ok := m.m[key]
	if !ok {
		return "", fmt.Errorf("secret not found: %s", key)
	}
	return v, nil
}

type MockPublicKeyProvider struct{}

func (pk *MockPublicKeyProvider) GetRepoPublicKey(ctx context.Context, orgName string, repoName string) (*github.PublicKey, *github.Response, error) {
	keyID := "testing"
	key := base64.StdEncoding.EncodeToString([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")) // 32B key
	pKey := github.PublicKey{
		KeyID: &keyID,
		Key:   &key,
	}
	return &pKey, nil, nil
}

func NewMockPublicKeyProvider() PublicKeyProvider {
	return &MockPublicKeyProvider{}
}

type MockSecretsService struct {
	secrets map[string]*github.EncryptedSecret
}

func (m MockSecretsService) ListRepoSecrets(_ context.Context, _, _ string, _ *github.ListOptions) (*github.Secrets, *github.Response, error) {
	secs := make([]*github.Secret, 0)
	for key, _ := range m.secrets {
		secs = append(secs, &github.Secret{
			Name: key,
		})
	}

	return &github.Secrets{
		TotalCount: len(m.secrets),
		Secrets:    secs,
	}, nil, nil
}

func (m MockSecretsService) DeleteRepoSecret(_ context.Context, _, _, name string) (*github.Response, error) {
	delete(m.secrets, name)
	return nil, nil
}

func (m MockSecretsService) CreateOrUpdateRepoSecret(_ context.Context, _, _ string, eSecret *github.EncryptedSecret) (*github.Response, error) {
	m.secrets[eSecret.Name] = eSecret
	return nil, nil
}

func NewMockSecretsService() ActionsSecretsService {
	mServ := &MockSecretsService{secrets: make(map[string]*github.EncryptedSecret, 0)}
	_, _ = mServ.CreateOrUpdateRepoSecret(context.Background(), "", "", &github.EncryptedSecret{
		Name:  "secret0",
		KeyID: "testing",
	})
	return mServ
}
