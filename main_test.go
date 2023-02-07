package main

import (
	"context"
	"testing"

	"github.com/FedeDP/GhEnvSet/pkg/config"
	"github.com/FedeDP/GhEnvSet/pkg/poiana"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	provider, err := poiana.NewMockSecretsProvider(map[string]string{
		"orgSecret0":  "orgValue0",
		"repoSecret0": "repoValue0",
		"repoSecret1": "repoValue1",
		"repoSecret2": "repoValue2",
	})
	require.NoError(t, err)

	// sync the org and the repos with the config
	require.Len(t, conf.Orgs, 1)
	for orgName, org := range conf.Orgs {
		mockVarServ := poiana.NewMockVariableService()
		mockSecServ := poiana.NewMockSecretsService()
		require.Equal(t, orgName, "FedeDP")
		require.NoError(t, org.Actions.Sync(ctx, provider, mockVarServ, mockSecServ, false))

		// Check variables
		vars, _, err := mockVarServ.ListVariables(ctx, nil)
		assert.NoError(t, err)
		assert.Equal(t, vars.TotalCount, len(conf.Orgs["FedeDP"].Actions.Variables))
		for _, v := range vars.Variables {
			assert.Equal(t, conf.Orgs["FedeDP"].Actions.Variables[v.Name], v.Value)
		}

		// Check secrets
		secs, _, err := mockSecServ.ListSecrets(ctx, nil)
		assert.NoError(t, err)
		assert.Equal(t, secs.TotalCount, len(conf.Orgs["FedeDP"].Actions.Secrets))
		for _, sec := range secs.Secrets {
			assert.Contains(t, conf.Orgs["FedeDP"].Actions.Secrets, sec.Name)
		}

		require.Len(t, org.Repos, 1)
		for repoName, repo := range org.Repos {
			mockVarServ := poiana.NewMockVariableService()
			mockSecServ := poiana.NewMockSecretsService()
			require.Equal(t, repoName, "GhEnvSet")
			require.NoError(t, repo.Actions.Sync(ctx, provider, mockVarServ, mockSecServ, false))

			// Check variables
			vars, _, err = mockVarServ.ListVariables(ctx, nil)
			assert.NoError(t, err)
			assert.Equal(t, vars.TotalCount, len(conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Variables))
			for _, v := range vars.Variables {
				assert.Equal(t, conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Variables[v.Name], v.Value)
			}

			// Check secrets
			secs, _, err = mockSecServ.ListSecrets(ctx, nil)
			assert.NoError(t, err)
			assert.Equal(t, secs.TotalCount, len(conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Secrets))
			for _, sec := range secs.Secrets {
				assert.Contains(t, conf.Orgs["FedeDP"].Repos["GhEnvSet"].Actions.Secrets, sec.Name)
			}
		}
	}
}
