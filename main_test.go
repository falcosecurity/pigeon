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

func TestMainLoop(t *testing.T) {
	ctx := context.Background()

	conf, err := config.FromData(testYAML)
	assert.NoError(t, err)

	mockVarServ := poiana.NewMockVariableService()
	mockSecServ := poiana.NewMockSecretsService()
	mockPKeyProv := poiana.NewMockPublicKeyProvider()
	provider, err := poiana.NewMockSecretsProvider(map[string]string{"secret0": "value0", "secret1": "value1", "secret2": "value2"})
	assert.NoError(t, err)

	err = conf.Loop(mockVarServ, mockSecServ, provider, mockPKeyProv, false)
	assert.NoError(t, err)

	// Check repo variables
	vars, _, err := mockVarServ.ListRepoVariables(ctx, "", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, vars.TotalCount, len(conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Variables))
	for _, v := range vars.Variables {
		assert.Equal(t, conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Variables[v.Name], v.Value)
	}

	// Check repo secrets
	secs, _, err := mockSecServ.ListRepoSecrets(ctx, "", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, secs.TotalCount, len(conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Secrets))
	for _, sec := range secs.Secrets {
		assert.Contains(t, conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Secrets, sec.Name)
	}
}
