package pigeon

import "fmt"

type mockServiceFactory struct {
	secs map[string]ActionsSecretsService
	vars map[string]ActionsVarsService
}

func NewMockServiceFactory() ServiceFactory {
	return mockServiceFactory{
		secs: make(map[string]ActionsSecretsService),
		vars: make(map[string]ActionsVarsService),
	}
}

func (c mockServiceFactory) NewOrgVariableService(org string) ActionsVarsService {
	if _, ok := c.vars[org]; !ok {
		c.vars[org] = NewMockVariableService()
	}
	return c.vars[org]
}

func (c mockServiceFactory) NewOrgSecretService(org string) ActionsSecretsService {
	if _, ok := c.secs[org]; !ok {
		c.secs[org] = NewMockSecretsService()
	}
	return c.secs[org]
}

func (c mockServiceFactory) NewRepoVariableService(org, repo string) ActionsVarsService {
	key := fmt.Sprintf("%s-%s", org, repo)
	if _, ok := c.vars[key]; !ok {
		c.vars[key] = NewMockVariableService()
	}
	return c.vars[key]
}

func (c mockServiceFactory) NewRepoSecretService(org, repo string) ActionsSecretsService {
	key := fmt.Sprintf("%s-%s", org, repo)
	if _, ok := c.secs[key]; !ok {
		c.secs[key] = NewMockSecretsService()
	}
	return c.secs[key]
}
