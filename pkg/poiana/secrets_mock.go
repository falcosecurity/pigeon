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

func (pk *MockPublicKeyProvider) GetOrgPublicKey(ctx context.Context, org string) (*github.PublicKey, *github.Response, error) {
	return pk.GetRepoPublicKey(ctx, org, "")
}

func NewMockPublicKeyProvider() PublicKeyProvider {
	return &MockPublicKeyProvider{}
}

type MockSecretsService struct {
	repoSecrets map[string]*github.EncryptedSecret
	orgSecrets  map[string]*github.EncryptedSecret
}

func (m MockSecretsService) ListRepoSecrets(_ context.Context, _, _ string, _ *github.ListOptions) (*github.Secrets, *github.Response, error) {
	secs := make([]*github.Secret, 0)
	for key, _ := range m.repoSecrets {
		secs = append(secs, &github.Secret{
			Name: key,
		})
	}

	return &github.Secrets{
		TotalCount: len(m.repoSecrets),
		Secrets:    secs,
	}, nil, nil
}

func (m MockSecretsService) DeleteRepoSecret(_ context.Context, _, _, name string) (*github.Response, error) {
	delete(m.repoSecrets, name)
	return nil, nil
}

func (m MockSecretsService) CreateOrUpdateRepoSecret(_ context.Context, _, _ string, eSecret *github.EncryptedSecret) (*github.Response, error) {
	m.repoSecrets[eSecret.Name] = eSecret
	return nil, nil
}

func (m MockSecretsService) ListOrgSecrets(ctx context.Context, owner string, opts *github.ListOptions) (*github.Secrets, *github.Response, error) {
	secs := make([]*github.Secret, 0)
	for key, _ := range m.orgSecrets {
		secs = append(secs, &github.Secret{
			Name: key,
		})
	}

	return &github.Secrets{
		TotalCount: len(m.orgSecrets),
		Secrets:    secs,
	}, nil, nil
}

func (m MockSecretsService) DeleteOrgSecret(ctx context.Context, owner, name string) (*github.Response, error) {
	delete(m.orgSecrets, name)
	return nil, nil
}

func (m MockSecretsService) CreateOrUpdateOrgSecret(ctx context.Context, owner string, eSecret *github.EncryptedSecret) (*github.Response, error) {
	m.orgSecrets[eSecret.Name] = eSecret
	return nil, nil
}

func NewMockSecretsService() ActionsSecretsService {
	mServ := &MockSecretsService{
		repoSecrets: make(map[string]*github.EncryptedSecret, 0),
		orgSecrets:  make(map[string]*github.EncryptedSecret, 0),
	}
	_, _ = mServ.CreateOrUpdateRepoSecret(context.Background(), "", "", &github.EncryptedSecret{
		Name:  "repoSecret0",
		KeyID: "repoTesting",
	})
	_, _ = mServ.CreateOrUpdateOrgSecret(context.Background(), "", &github.EncryptedSecret{
		Name:  "orgSecret0",
		KeyID: "orgTesting",
	})
	return mServ
}
