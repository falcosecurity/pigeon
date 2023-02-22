package config

import (
	"context"
	"encoding/base64"
	"os"
	"strings"

	"github.com/falcosecurity/pigeon/pkg/pigeon"

	"github.com/google/go-github/v50/github"
	"github.com/jamesruan/sodium"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type GitHubActionsConfig struct {
	Variables map[string]string `yaml:"variables"`
	Secrets   []string          `yaml:"secrets"`
}

type GitHubRepo struct {
	Actions GitHubActionsConfig `yaml:"actions"`
}

type GitHubOrg struct {
	Repos   map[string]GitHubRepo
	Actions GitHubActionsConfig `yaml:"actions"`
}

type GithubConfig struct {
	Orgs map[string]GitHubOrg `yaml:"orgs"`
}

func FromFile(fileName string) (*GithubConfig, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return FromData(string(b))
}

func FromData(yamlData string) (*GithubConfig, error) {
	var conf GithubConfig
	err := yaml.NewDecoder(strings.NewReader(yamlData)).Decode(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func (goa *GitHubActionsConfig) syncSecrets(
	ctx context.Context,
	service pigeon.ActionsSecretsService,
	provider pigeon.SecretsProvider,
	pKey *github.PublicKey) error {

	// Step 1: load repo secrets
	secs, _, err := service.ListSecrets(ctx, nil)
	if err != nil {
		return err
	}

	// Step 2: delete all secrets that are no more existing
	found := false
	for _, existentSec := range secs.Secrets {
		for _, newSec := range goa.Secrets {
			if newSec == existentSec.Name {
				found = true
				break
			}
		}
		if !found {
			_, err = service.DeleteSecret(ctx, existentSec.Name)
			if err != nil {
				return err
			}
		}
	}

	keyBytes, err := base64.StdEncoding.DecodeString(pKey.GetKey())
	if err != nil {
		return err
	}

	// Step 3: add or update all conf-listed secrets
	for _, secName := range goa.Secrets {
		secValue, err := provider.GetSecret(secName)
		if err != nil {
			return err
		}

		secBytes := sodium.Bytes(secValue)
		encSecBytes := secBytes.SealedBox(sodium.BoxPublicKey{Bytes: keyBytes})
		encSecBytesB64 := base64.StdEncoding.EncodeToString(([]byte)(encSecBytes))
		_, err = service.CreateOrUpdateSecret(ctx, &github.EncryptedSecret{
			Name:           secName,
			KeyID:          pKey.GetKeyID(),
			EncryptedValue: encSecBytesB64,
		})
		if err != nil {
			return err
		}
	}
	// logrus.Infof("secrets synced for org %s\n", orgName)
	return nil
}

func (goa *GitHubActionsConfig) syncVariables(
	ctx context.Context,
	service pigeon.ActionsVarsService) error {
	// Step 1: load repo variables

	vars, _, err := service.ListVariables(ctx, nil)
	if err != nil {
		return err
	}

	// Step 2: delete all variables that are no more existing
	for _, existentVar := range vars.Variables {
		_, ok := goa.Variables[existentVar.Name]
		if !ok {
			_, err = service.DeleteVariable(ctx, existentVar.Name)
			if err != nil {
				return err
			}
		}
	}

	// Step 3: add or update all conf-listed variables
	for newVarName, newVarValue := range goa.Variables {
		_, err := service.CreateOrUpdateVariable(ctx, &github.ActionsVariable{
			Name:  newVarName,
			Value: newVarValue,
		})
		if err != nil {
			return err
		}
	}
	// logrus.Infof("variables synced for org %s\n", orgName)
	return nil
}

func (g *GitHubActionsConfig) Sync(
	ctx context.Context,
	provider pigeon.SecretsProvider,
	vars pigeon.ActionsVarsService,
	secrets pigeon.ActionsSecretsService,
	dryRun bool) error {
	if dryRun {
		logrus.Infoln("skipping secrets sync due to dry run")
		logrus.Infoln("skipping variables sync due to dry run")
		return nil
	}

	pKey, _, err := secrets.GetPublicKey(ctx)
	if err != nil {
		return err
	}
	err = g.syncSecrets(ctx, secrets, provider, pKey)
	if err != nil {
		return err
	}
	return g.syncVariables(ctx, vars)
}

func (g *GithubConfig) Sync(f pigeon.ServiceFactory, p pigeon.SecretsProvider, dryRun bool) error {
	logrus.Debugf("starting the synching loop")
	ctx := context.Background()
	for orgName, org := range g.Orgs {
		logrus.Debugf("synching org %s", orgName)
		err := org.Actions.Sync(
			ctx, p,
			f.NewOrgVariableService(orgName),
			f.NewOrgSecretService(orgName),
			dryRun)
		if err != nil {
			return err
		}

		// note: we want to ignore the repos that are not in the config
		for repoName, repo := range org.Repos {
			logrus.Debugf("synching repo %s/%s", orgName, repoName)
			err := repo.Actions.Sync(
				ctx, p,
				f.NewRepoVariableService(orgName, repoName),
				f.NewRepoSecretService(orgName, repoName),
				dryRun)
			if err != nil {
				return err
			}
		}
	}
	logrus.Debugf("ending the synching loop")
	return nil
}
