package pigeon

import (
	"context"
	"net/http"

	"github.com/google/go-github/v50/github"
	"github.com/sirupsen/logrus"
)

type clientOrgVariableService struct {
	client *github.Client
	org    string
}

func NewClientOrgVariableService(client *github.Client, org string) ActionsVarsService {
	return &clientOrgVariableService{client: client, org: org}
}

func (c *clientOrgVariableService) ListVariables(ctx context.Context, opts *github.ListOptions) (*github.ActionsVariables, *github.Response, error) {
	logrus.Infof("listing variables for org '%s'...", c.org)
	return c.client.Actions.ListOrgVariables(ctx, c.org, opts)
}

func (c *clientOrgVariableService) DeleteVariable(ctx context.Context, name string) (*github.Response, error) {
	logrus.Infof("deleting variable '%s' for org '%s'...", name, c.org)
	return c.client.Actions.DeleteOrgVariable(ctx, c.org, name)
}

func (c *clientOrgVariableService) CreateOrUpdateVariable(ctx context.Context, variable *github.ActionsVariable) (*github.Response, error) {
	logrus.Infof("creating or updating variable '%s' in org '%s'...", variable.Name, c.org)
	resp, err := c.client.Actions.UpdateOrgVariable(ctx, c.org, variable)
	if err != nil {
		// Update returns StatusNoContent when successful
		if resp.StatusCode != http.StatusNoContent {
			return c.client.Actions.CreateOrgVariable(ctx, c.org, variable)
		}
	}
	return resp, err
}

type clientRepoVariableService struct {
	client *github.Client
	org    string
	repo   string
}

func NewClientRepoVariableService(client *github.Client, org, repo string) ActionsVarsService {
	return &clientRepoVariableService{client: client, org: org, repo: repo}
}

func (c *clientRepoVariableService) ListVariables(ctx context.Context, opts *github.ListOptions) (*github.ActionsVariables, *github.Response, error) {
	logrus.Infof("listing variables for repo '%s/%s'...", c.org, c.repo)
	return c.client.Actions.ListRepoVariables(ctx, c.org, c.repo, opts)
}

func (c *clientRepoVariableService) DeleteVariable(ctx context.Context, name string) (*github.Response, error) {
	logrus.Infof("deleting variable '%s' for repo '%s/%s'...", name, c.org, c.repo)
	return c.client.Actions.DeleteRepoVariable(ctx, c.org, c.repo, name)
}

func (c *clientRepoVariableService) CreateOrUpdateVariable(ctx context.Context, variable *github.ActionsVariable) (*github.Response, error) {
	logrus.Infof("creating or updating variable '%s' for repo '%s/%s'...", variable.Name, c.org, c.repo)
	resp, err := c.client.Actions.UpdateRepoVariable(ctx, c.org, c.repo, variable)
	if err != nil {
		// Update returns StatusNoContent when successful
		if resp.StatusCode != http.StatusNoContent {
			return c.client.Actions.CreateRepoVariable(ctx, c.org, c.repo, variable)
		}
	}
	return resp, err
}
