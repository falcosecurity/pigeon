package pigeon

import (
	"context"

	"github.com/google/go-github/v50/github"
	"github.com/sirupsen/logrus"
)

type clientOrgSecretService struct {
	client *github.Client
	org    string
}

func NewClientOrgSecretService(client *github.Client, org string) ActionsSecretsService {
	return &clientOrgSecretService{client: client, org: org}
}

func (c *clientOrgSecretService) GetPublicKey(ctx context.Context) (*github.PublicKey, *github.Response, error) {
	logrus.Infof("retrieving public key for org '%s'...", c.org)
	return c.client.Actions.GetOrgPublicKey(ctx, c.org)
}

func (c *clientOrgSecretService) ListSecrets(ctx context.Context, opts *github.ListOptions) (*github.Secrets, *github.Response, error) {
	logrus.Infof("listing secrets for org '%s'...", c.org)
	return c.client.Actions.ListOrgSecrets(ctx, c.org, opts)
}

func (c *clientOrgSecretService) DeleteSecret(ctx context.Context, name string) (*github.Response, error) {
	logrus.Infof("deleting secret '%s' for org '%s'...", name, c.org)
	return c.client.Actions.DeleteOrgSecret(ctx, c.org, name)
}

func (c *clientOrgSecretService) CreateOrUpdateSecret(ctx context.Context, eSecret *github.EncryptedSecret) (*github.Response, error) {
	logrus.Infof("creating or updating secret '%s' in org '%s'...", eSecret.Name, c.org)
	return c.client.Actions.CreateOrUpdateOrgSecret(ctx, c.org, eSecret)
}

type clientRepoSecretService struct {
	client *github.Client
	org    string
	repo   string
}

func NewClientRepoSecretService(client *github.Client, org, repo string) ActionsSecretsService {
	return &clientRepoSecretService{client: client, org: org, repo: repo}
}

func (c *clientRepoSecretService) GetPublicKey(ctx context.Context) (*github.PublicKey, *github.Response, error) {
	logrus.Infof("retrieving public key for repo '%s/%s'...", c.org, c.repo)
	return c.client.Actions.GetRepoPublicKey(ctx, c.org, c.repo)
}

func (c *clientRepoSecretService) ListSecrets(ctx context.Context, opts *github.ListOptions) (*github.Secrets, *github.Response, error) {
	logrus.Infof("listing secrets for repo '%s/%s'...", c.org, c.repo)
	return c.client.Actions.ListRepoSecrets(ctx, c.org, c.repo, opts)
}

func (c *clientRepoSecretService) DeleteSecret(ctx context.Context, name string) (*github.Response, error) {
	logrus.Infof("deleting secret '%s' for repo '%s/%s'...", name, c.org, c.repo)
	return c.client.Actions.DeleteRepoSecret(ctx, c.org, c.repo, name)
}

func (c *clientRepoSecretService) CreateOrUpdateSecret(ctx context.Context, eSecret *github.EncryptedSecret) (*github.Response, error) {
	logrus.Infof("creating or updating secret '%s' for repo '%s/%s'...", eSecret.Name, c.org, c.repo)
	return c.client.Actions.CreateOrUpdateRepoSecret(ctx, c.org, c.repo, eSecret)
}
