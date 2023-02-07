package config

import (
	"context"
	"encoding/base64"
	"github.com/FedeDP/GhEnvSet/pkg/poiana"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/jamesruan/sodium"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type GitHubRepoActions struct {
	Variables map[string]string `yaml:"variables"`
	Secrets   []string          `yaml:"secrets"`
}

type GitHubRepo struct {
	Actions GitHubRepoActions `yaml:"actions"`
}

type GitHubOrgActions struct {
	Variables map[string]string `yaml:"variables"`
	Secrets   []string          `yaml:"secrets"`
}

type GitHubOrg struct {
	Repos   map[string]GitHubRepo
	Actions GitHubOrgActions `yaml:"actions"`
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

func (goa *GitHubOrgActions) syncSecrets(ctx context.Context,
	service poiana.ActionsOrgSecretsService,
	provider poiana.SecretsProvider,
	pKey *github.PublicKey,
	orgName string) error {

	// Step 1: load repo secrets
	logrus.Infof("listing secrets for org '%s'...", orgName)
	secs, _, err := service.ListOrgSecrets(ctx, orgName, nil)
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
			logrus.Infof("deleting secret '%s' for org '%s'...", existentSec.Name, orgName)
			_, err = service.DeleteOrgSecret(ctx, orgName, existentSec.Name)
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
		logrus.Infof("adding/updating secret '%s' in org '%s'...", secName, orgName)
		secValue, err := provider.GetSecret(secName)
		if err != nil {
			return err
		}

		secBytes := sodium.Bytes(secValue)
		encSecBytes := secBytes.SealedBox(sodium.BoxPublicKey{Bytes: keyBytes})
		encSecBytesB64 := base64.StdEncoding.EncodeToString(([]byte)(encSecBytes))
		_, err = service.CreateOrUpdateOrgSecret(ctx, orgName, &github.EncryptedSecret{
			Name:           secName,
			KeyID:          pKey.GetKeyID(),
			EncryptedValue: encSecBytesB64,
		})
		if err != nil {
			return err
		}
	}
	logrus.Infof("secrets synced for org %s\n", orgName)
	return nil
}

func (goa *GitHubOrgActions) syncVariables(ctx context.Context,
	service poiana.ActionsOrgVarsService,
	orgName string) error {
	// Step 1: load repo variables
	logrus.Infof("listing variables for org '%s'...", orgName)
	vars, _, err := service.ListOrgVariables(ctx, orgName, nil)
	if err != nil {
		return err
	}

	// Step 2: delete all variables that are no more existing
	for _, existentVar := range vars.Variables {
		_, ok := goa.Variables[existentVar.Name]
		if !ok {
			logrus.Infof("deleting variable '%s' for org '%s'...", existentVar.Name, orgName)
			_, err = service.DeleteOrgVariable(ctx, orgName, existentVar.Name)
			if err != nil {
				return err
			}
		}
	}

	// Step 3: add or update all conf-listed variables
	for newVarName, newVarValue := range goa.Variables {
		logrus.Infof("adding/updating variable '%s' in org '%s'...", newVarName, orgName)
		resp, err := service.UpdateOrgVariable(ctx, orgName, &github.ActionsVariable{
			Name:  newVarName,
			Value: newVarValue,
		})
		if err != nil {
			// Update returns StatusNoContent when successful
			if resp.StatusCode != http.StatusNoContent {
				_, err = service.CreateOrgVariable(ctx, orgName, &github.ActionsVariable{
					Name:  newVarName,
					Value: newVarValue,
				})
			}
			if err != nil {
				return err
			}
		}
	}
	logrus.Infof("variables synced for org %s\n", orgName)
	return nil
}

func (gra *GitHubRepoActions) syncSecrets(ctx context.Context,
	service poiana.ActionsRepoSecretsService,
	provider poiana.SecretsProvider,
	pKey *github.PublicKey,
	orgName, repoName string) error {

	// Step 1: load repo secrets
	logrus.Infof("listing secrets for repo '%s/%s'...", orgName, repoName)
	secs, _, err := service.ListRepoSecrets(ctx, orgName, repoName, nil)
	if err != nil {
		return err
	}

	// Step 2: delete all secrets that are no more existing
	found := false
	for _, existentSec := range secs.Secrets {
		for _, newSec := range gra.Secrets {
			if newSec == existentSec.Name {
				found = true
				break
			}
		}
		if !found {
			logrus.Infof("deleting secret '%s' for repo '%s/%s'...", existentSec.Name, orgName, repoName)
			_, err = service.DeleteRepoSecret(ctx, orgName, repoName, existentSec.Name)
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
	for _, secName := range gra.Secrets {
		logrus.Infof("adding/updating secret '%s' in repo '%s/%s'...", secName, orgName, repoName)
		secValue, err := provider.GetSecret(secName)
		if err != nil {
			return err
		}

		secBytes := sodium.Bytes(secValue)
		encSecBytes := secBytes.SealedBox(sodium.BoxPublicKey{Bytes: keyBytes})
		encSecBytesB64 := base64.StdEncoding.EncodeToString(([]byte)(encSecBytes))
		_, err = service.CreateOrUpdateRepoSecret(ctx, orgName, repoName, &github.EncryptedSecret{
			Name:           secName,
			KeyID:          pKey.GetKeyID(),
			EncryptedValue: encSecBytesB64,
		})
		if err != nil {
			return err
		}
	}
	logrus.Infof("secrets synced for repo %s/%s\n", orgName, repoName)
	return nil
}

func (gra *GitHubRepoActions) syncVariables(ctx context.Context,
	service poiana.ActionsRepoVarsService,
	orgName,
	repoName string) error {
	// Step 1: load repo variables
	logrus.Infof("listing variables for repo '%s/%s'...", orgName, repoName)
	vars, _, err := service.ListRepoVariables(ctx, orgName, repoName, nil)
	if err != nil {
		return err
	}

	// Step 2: delete all variables that are no more existing
	for _, existentVar := range vars.Variables {
		_, ok := gra.Variables[existentVar.Name]
		if !ok {
			logrus.Infof("deleting variable '%s' for repo '%s/%s'...", existentVar.Name, orgName, repoName)
			_, err = service.DeleteRepoVariable(ctx, orgName, repoName, existentVar.Name)
			if err != nil {
				return err
			}
		}
	}

	// Step 3: add or update all conf-listed variables
	for newVarName, newVarValue := range gra.Variables {
		logrus.Infof("adding/updating variable '%s' in repo '%s/%s'...", newVarName, orgName, repoName)
		resp, err := service.UpdateRepoVariable(ctx, orgName, repoName, &github.ActionsVariable{
			Name:  newVarName,
			Value: newVarValue,
		})
		if err != nil {
			// Update returns StatusNoContent when successful
			if resp.StatusCode != http.StatusNoContent {
				_, err = service.CreateRepoVariable(ctx, orgName, repoName, &github.ActionsVariable{
					Name:  newVarName,
					Value: newVarValue,
				})
			}
			if err != nil {
				return err
			}
		}
	}
	logrus.Infof("variables synced for repo %s/%s\n", orgName, repoName)
	return nil
}

func (g *GithubConfig) Loop(
	vService poiana.ActionsVarsService,
	sService poiana.ActionsSecretsService,
	provider poiana.SecretsProvider,
	pKeyProvider poiana.PublicKeyProvider,
	dryRun bool,
) error {
	ctx := context.Background()
	var (
		err  error
		pKey *github.PublicKey
	)

	logrus.Debugln("Starting the loop")
	for orgName, org := range g.Orgs {
		logrus.Infof("retrieving public key for org '%s'...", orgName)
		if dryRun {
			logrus.Infoln("Would have synced org secrets")
		} else {
			pKey, _, err = pKeyProvider.GetOrgPublicKey(ctx, orgName)
			if err == nil {
				err = org.Actions.syncSecrets(ctx, sService, provider, pKey, orgName)
			}
			if err != nil {
				return err
			}
		}

		if dryRun {
			logrus.Infoln("Would have synced org variables")
		} else {
			err = org.Actions.syncVariables(ctx, vService, orgName)
			if err != nil {
				return err
			}
		}

		// todo: also remove all secrets and vars for all repos not present
		// in the YAML config
		for repoName, repo := range org.Repos {
			// fetch encryption key
			logrus.Infof("retrieving public key for repo '%s/%s'...", orgName, repoName)
			if dryRun {
				logrus.Infoln("Would have synced repo secrets")
			} else {
				pKey, _, err = pKeyProvider.GetRepoPublicKey(ctx, orgName, repoName)
				if err == nil {
					err = repo.Actions.syncSecrets(ctx, sService, provider, pKey, orgName, repoName)
				}
				if err != nil {
					return err
				}
			}

			if dryRun {
				logrus.Infoln("Would have synced repo variables")
			} else {
				err = repo.Actions.syncVariables(ctx, vService, orgName, repoName)
				if err != nil {
					return err
				}
			}
		}
	}
	logrus.Debugln("Ended loop")
	return nil
}
