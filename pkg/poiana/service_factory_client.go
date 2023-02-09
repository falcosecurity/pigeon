package poiana

import "github.com/google/go-github/v50/github"

type clientServiceFactory struct {
	client *github.Client
}

func NewClientServiceFactory(client *github.Client) ServiceFactory {
	return &clientServiceFactory{client: client}
}

func (c *clientServiceFactory) NewOrgVariableService(org string) ActionsVarsService {
	return NewClientOrgVariableService(c.client, org)
}

func (c *clientServiceFactory) NewOrgSecretService(org string) ActionsSecretsService {
	return NewClientOrgSecretService(c.client, org)
}

func (c *clientServiceFactory) NewRepoVariableService(org, repo string) ActionsVarsService {
	return NewClientRepoVariableService(c.client, org, repo)
}

func (c *clientServiceFactory) NewRepoSecretService(org, repo string) ActionsSecretsService {
	return NewClientRepoSecretService(c.client, org, repo)
}
