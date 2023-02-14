package config

import (
	"context"
	"github.com/FedeDP/Pigeon/pkg/pigeon"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testYAML = `
orgs:
  myOrg:
    actions:
      variables:
        orgVar1: "orgValue1"
      secrets:
        - orgSecret0
    repos:
      myRepo:
        actions:
          variables:
            repoVar1: "repoValue1"
            repoVar2: "repoValue2"
          secrets:
            - repoSecret0
            - repoSecret1
            - repoSecret2
`

func TestConfigSync(t *testing.T) {
	org := "myOrg"
	repo := "myRepo"
	ctx := context.Background()

	conf, err := FromData(testYAML)
	assert.NoError(t, err)

	// Check correct yaml parsing
	assert.Contains(t, conf.Orgs, org)
	assert.Contains(t, conf.Orgs[org].Actions.Secrets, "orgSecret0")
	assert.Contains(t, conf.Orgs[org].Actions.Variables, "orgVar1")
	assert.Contains(t, conf.Orgs[org].Repos, repo)

	factory := pigeon.NewMockServiceFactory()
	provider, err := pigeon.NewMockSecretsProvider(map[string]string{
		"orgSecret0":  "orgValue0",
		"repoSecret0": "repoValue0",
		"repoSecret1": "repoValue1",
		"repoSecret2": "repoValue2",
	})
	assert.NoError(t, err)

	err = conf.Sync(factory, provider, false)
	assert.NoError(t, err)

	// Org
	// Check variables
	orgVars := factory.NewOrgVariableService(org)
	vars, _, err := orgVars.ListVariables(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, vars.TotalCount, len(conf.Orgs[org].Actions.Variables))
	for _, v := range vars.Variables {
		assert.Equal(t, conf.Orgs[org].Actions.Variables[v.Name], v.Value)
	}

	// Check secrets
	orgSecs := factory.NewOrgSecretService(org)
	secs, _, err := orgSecs.ListSecrets(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, secs.TotalCount, len(conf.Orgs[org].Actions.Secrets))
	for _, sec := range secs.Secrets {
		assert.Contains(t, conf.Orgs[org].Actions.Secrets, sec.Name)
	}

	// Repo
	// Check variables
	repoVars := factory.NewRepoVariableService(org, repo)
	vars, _, err = repoVars.ListVariables(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, vars.TotalCount, len(conf.Orgs[org].Repos[repo].Actions.Variables))
	for _, v := range vars.Variables {
		assert.Equal(t, conf.Orgs[org].Repos[repo].Actions.Variables[v.Name], v.Value)
	}

	// Check secrets
	repoSecs := factory.NewRepoSecretService(org, repo)
	secs, _, err = repoSecs.ListSecrets(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, secs.TotalCount, len(conf.Orgs[org].Repos[repo].Actions.Secrets))
	for _, sec := range secs.Secrets {
		assert.Contains(t, conf.Orgs[org].Repos[repo].Actions.Secrets, sec.Name)
	}
}
