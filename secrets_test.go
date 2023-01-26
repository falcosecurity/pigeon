package main

const sampleYAML = `
orgs:
  FedeDP:
    repos:
      GhEnvSet:
        actions:
          variables:
            SOME_VARIABLE2: "ciao"
          secrets:
            - TEST_SECRET_KEY
`

// import (
// 	"context"
// 	"encoding/base64"
// 	"github.com/FedeDP/GhEnvSet/pkg/poiana"
// 	"github.com/google/go-github/v49/github"
// 	"github.com/stretchr/testify/assert"
// 	"testing"
// )

// type MockSecretsService struct {
// 	secrets map[string]*github.EncryptedSecret
// }

// func (m MockSecretsService) ListRepoSecrets(ctx context.Context, owner, repo string) (*github.Secrets, error) {
// 	secs := make([]*github.Secret, 0)
// 	for key, _ := range m.secrets {
// 		secs = append(secs, &github.Secret{
// 			Name: key,
// 		})
// 	}

// 	return &github.Secrets{
// 		TotalCount: len(m.secrets),
// 		Secrets:    secs,
// 	}, nil
// }

// func (m MockSecretsService) DeleteRepoSecret(ctx context.Context, owner, repo, name string) error {
// 	delete(m.secrets, name)
// 	return nil
// }

// func (m MockSecretsService) CreateOrUpdateRepoSecret(ctx context.Context, owner, repo string, eSecret *github.EncryptedSecret) error {
// 	m.secrets[eSecret.Name] = eSecret
// 	return nil
// }

// func newMockSecretsService() poiana.ActionsSecretsService {
// 	mServ := &MockSecretsService{secrets: make(map[string]*github.EncryptedSecret, 0)}
// 	_ = mServ.CreateOrUpdateRepoSecret(context.Background(), "", "", &github.EncryptedSecret{
// 		Name:  "secret0",
// 		KeyID: "testing",
// 	})
// 	return mServ
// }

// func TestSyncServices(t *testing.T) {
// 	ctx := context.Background()
// 	secrets := []string{
// 		"secret1", "secret2",
// 	}

// 	mockServ := newMockSecretsService()
// 	provider, err := poiana.NewMockSecretsProvider(map[string]string{"secret1": "value1", "secret2": "value2"})
// 	assert.NoError(t, err)

// 	keyID := "testing"
// 	key := base64.StdEncoding.EncodeToString([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")) // 32B key
// 	pKey := github.PublicKey{
// 		KeyID: &keyID,
// 		Key:   &key,
// 	}
// 	err = poiana.syncSecrets(ctx, mockServ, provider, &pKey, "", "", secrets)
// 	assert.NoError(t, err)

// 	secs, err := mockServ.ListRepoSecrets(ctx, "", "")
// 	assert.NoError(t, err)

// 	for _, sec := range secs.Secrets {
// 		assert.Contains(t, secrets, sec.Name)
// 	}
// }
