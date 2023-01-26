package main

import (
	"context"
	"encoding/base64"
	"github.com/FedeDP/GhEnvSet/pkg/config"
	"github.com/FedeDP/GhEnvSet/pkg/poiana"
	"github.com/google/go-github/v49/github"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testYAML = `
orgs:
  FedeDP:
    repos:
      GhEnvSet:
        actions:
          variables:
            var1: "value1"
            var2: "value2"
          secrets:
            - secret0
            - secret1
            - secret2
`

type MockVariableService struct {
	variables map[string]string
}

func (m MockVariableService) ListRepoVariables(ctx context.Context, owner, repo string) (*poiana.Variables, error) {
	vars := make([]*poiana.Variable, 0)
	for key, val := range m.variables {
		vars = append(vars, &poiana.Variable{
			Name:  key,
			Value: val,
		})
	}

	return &poiana.Variables{
		TotalCount: len(m.variables),
		Variables:  vars,
	}, nil
}

func (m MockVariableService) DeleteRepoVariable(ctx context.Context, owner, repo, name string) error {
	delete(m.variables, name)
	return nil
}

func (m MockVariableService) CreateOrUpdateRepoVariable(ctx context.Context, owner, repo string, variable *poiana.Variable) error {
	m.variables[variable.Name] = variable.Value
	return nil
}

func newMockVariableService() poiana.ActionsVarsService {
	mServ := &MockVariableService{variables: make(map[string]string, 0)}
	// Initial variable set
	_ = mServ.CreateOrUpdateRepoVariable(context.Background(), "", "", &poiana.Variable{
		Name:      "var0",
		Value:     "value0",
		CreatedAt: github.Timestamp{},
		UpdatedAt: github.Timestamp{},
	})
	return mServ
}

type MockSecretsService struct {
	secrets map[string]*github.EncryptedSecret
}

func (m MockSecretsService) ListRepoSecrets(ctx context.Context, owner, repo string) (*github.Secrets, error) {
	secs := make([]*github.Secret, 0)
	for key, _ := range m.secrets {
		secs = append(secs, &github.Secret{
			Name: key,
		})
	}

	return &github.Secrets{
		TotalCount: len(m.secrets),
		Secrets:    secs,
	}, nil
}

func (m MockSecretsService) DeleteRepoSecret(ctx context.Context, owner, repo, name string) error {
	delete(m.secrets, name)
	return nil
}

func (m MockSecretsService) CreateOrUpdateRepoSecret(ctx context.Context, owner, repo string, eSecret *github.EncryptedSecret) error {
	m.secrets[eSecret.Name] = eSecret
	return nil
}

func newMockSecretsService() poiana.ActionsSecretsService {
	mServ := &MockSecretsService{secrets: make(map[string]*github.EncryptedSecret, 0)}
	_ = mServ.CreateOrUpdateRepoSecret(context.Background(), "", "", &github.EncryptedSecret{
		Name:  "secret0",
		KeyID: "testing",
	})
	return mServ
}

type MockPublicKeyProvider struct{}

func (pk *MockPublicKeyProvider) GetPublicKey(ctx context.Context, orgName string, repoName string) (*github.PublicKey, error) {
	keyID := "testing"
	key := base64.StdEncoding.EncodeToString([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")) // 32B key
	pKey := github.PublicKey{
		KeyID: &keyID,
		Key:   &key,
	}
	return &pKey, nil
}

func newMockPublicKeyProvider() poiana.PublicKeyProvider {
	return &MockPublicKeyProvider{}
}

func TestMainLoop(t *testing.T) {
	ctx := context.Background()

	conf, err := config.FromData(testYAML)
	assert.NoError(t, err)

	mockVarServ := newMockVariableService()
	mockSecServ := newMockSecretsService()
	mockPKeyProv := newMockPublicKeyProvider()
	provider, err := poiana.NewMockSecretsProvider(map[string]string{"secret0": "value0", "secret1": "value1", "secret2": "value2"})
	assert.NoError(t, err)

	err = conf.Loop(mockVarServ, mockSecServ, provider, mockPKeyProv, false)
	assert.NoError(t, err)

	// Check repo variables
	vars, err := mockVarServ.ListRepoVariables(ctx, "", "")
	assert.NoError(t, err)
	assert.Equal(t, vars.TotalCount, len(conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Variables))
	for _, v := range vars.Variables {
		assert.Equal(t, conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Variables[v.Name], v.Value)
	}

	// Check repo secrets
	secs, err := mockSecServ.ListRepoSecrets(ctx, "", "")
	assert.NoError(t, err)
	assert.Equal(t, secs.TotalCount, len(conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Secrets))
	for _, sec := range secs.Secrets {
		assert.Contains(t, conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Secrets, sec.Name)
	}
}
