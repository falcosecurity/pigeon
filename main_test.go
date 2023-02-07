package main

import (
	"context"
	"github.com/FedeDP/GhEnvSet/pkg/config"
	"github.com/FedeDP/GhEnvSet/pkg/poiana"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testYAML = `
orgs:
  FedeDP:
    actions:
      variables:
        orgVar1: "orgValue1"
      secrets:
        - orgSecret0
    repos:
      GhEnvSet:
        actions:
          variables:
            repoVar1: "repoValue1"
            repoVar2: "repoValue2"
          secrets:
            - repoSecret0
            - repoSecret1
            - repoSecret2
`

func TestMainLoop(t *testing.T) {
	ctx := context.Background()

	conf, err := config.FromData(testYAML)
	assert.NoError(t, err)

	mockVarServ := poiana.NewMockVariableService()
	mockSecServ := poiana.NewMockSecretsService()
	mockPKeyProv := poiana.NewMockPublicKeyProvider()
	provider, err := poiana.NewMockSecretsProvider(map[string]string{"orgSecret0": "orgValue0", "repoSecret0": "repoValue0", "repoSecret1": "repoValue1", "repoSecret2": "repoValue2"})
	assert.NoError(t, err)

	err = conf.Loop(mockVarServ, mockSecServ, provider, mockPKeyProv, false)
	assert.NoError(t, err)

	// Org
	// Check variables
	vars, _, err := mockVarServ.ListOrgVariables(ctx, "", nil)
	assert.NoError(t, err)
	assert.Equal(t, vars.TotalCount, len(conf.Orgs["FedeDP"].Actions.Variables))
	for _, v := range vars.Variables {
		assert.Equal(t, conf.Orgs["FedeDP"].Actions.Variables[v.Name], v.Value)
	}

	// Check secrets
	secs, _, err := mockSecServ.ListOrgSecrets(ctx, "", nil)
	assert.NoError(t, err)
	assert.Equal(t, secs.TotalCount, len(conf.Orgs["FedeDP"].Actions.Secrets))
	for _, sec := range secs.Secrets {
		assert.Contains(t, conf.Orgs["FedeDP"].Actions.Secrets, sec.Name)
	}

	// Repo
	// Check variables
	vars, _, err = mockVarServ.ListRepoVariables(ctx, "", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, vars.TotalCount, len(conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Variables))
	for _, v := range vars.Variables {
		assert.Equal(t, conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Variables[v.Name], v.Value)
	}

	// Check secrets
	secs, _, err = mockSecServ.ListRepoSecrets(ctx, "", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, secs.TotalCount, len(conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Secrets))
	for _, sec := range secs.Secrets {
		assert.Contains(t, conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Secrets, sec.Name)
	}
}
