package poiana

import (
	"context"
	"encoding/base64"
	"io"
	"os"
	"strings"

	"github.com/google/go-github/v49/github"
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

type GitHubOrg struct {
	Repos map[string]GitHubRepo
}

type GithubConfig struct {
	Orgs map[string]GitHubOrg `yaml:"orgs"`
}

func (e *GithubConfig) Decode(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(e)
}

func (e *GithubConfig) Encode(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(e)
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
	err := conf.Decode(strings.NewReader(yamlData))
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func syncSecrets(ctx context.Context,
	service ActionsSecretsService,
	provider SecretsProvider,
	pKey *github.PublicKey,
	orgName, repoName string,
	secrets []string) error {

	// Step 1: load repo secrets
	logrus.Infof("listing secrets for repo '%s/%s'...", orgName, repoName)
	secs, err := service.ListRepoSecrets(ctx, orgName, repoName)
	if err != nil {
		return err
	}

	// Step 2: delete all secrets that are no more existing
	found := false
	for _, existentSec := range secs.Secrets {
		for _, newSec := range secrets {
			if newSec == existentSec.Name {
				found = true
				break
			}
		}
		if !found {
			logrus.Infof("deleting secret '%s' for repo '%s/%s'...", existentSec.Name, orgName, repoName)
			err = service.DeleteRepoSecret(ctx, orgName, repoName, existentSec.Name)
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
	for _, secName := range secrets {
		logrus.Infof("adding/updating secret '%s' in repo '%s/%s'...", secName, orgName, repoName)
		secValue, err := provider.GetSecret(secName)
		if err != nil {
			return err
		}

		secBytes := sodium.Bytes(secValue)
		encSecBytes := secBytes.SealedBox(sodium.BoxPublicKey{Bytes: keyBytes})
		encSecBytesB64 := base64.StdEncoding.EncodeToString(([]byte)(encSecBytes))
		if err != nil {
			return err
		}
		err = service.CreateOrUpdateRepoSecret(ctx, orgName, repoName, &github.EncryptedSecret{
			Name:           secName,
			KeyID:          pKey.GetKeyID(),
			EncryptedValue: encSecBytesB64,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func syncVariables(ctx context.Context,
	service ActionsVarsService,
	orgName,
	repoName string,
	variables map[string]string) error {
	// Step 1: load repo variables
	logrus.Infof("listing variables for repo '%s/%s'...", orgName, repoName)
	vars, err := service.ListRepoVariables(ctx, orgName, repoName)
	if err != nil {
		return err
	}

	// Step 2: delete all variables that are no more existing
	for _, existentVar := range vars.Variables {
		_, ok := variables[existentVar.Name]
		if !ok {
			logrus.Infof("deleting variable '%s' for repo '%s/%s'...", existentVar.Name, orgName, repoName)
			err = service.DeleteRepoVariable(ctx, orgName, repoName, existentVar.Name)
			if err != nil {
				return err
			}
		}
	}

	// Step 3: add or update all conf-listed variables
	for newVarName, newVarValue := range variables {
		logrus.Infof("adding/updating variable '%s' in repo '%s/%s'...", newVarName, orgName, repoName)
		err = service.CreateOrUpdateRepoVariable(ctx, orgName, repoName, &Variable{
			Name:  newVarName,
			Value: newVarValue,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GithubConfig) Loop(vService ActionsVarsService, sService ActionsSecretsService, provider SecretsProvider, pKeyProvider PublicKeyProvider) error {
	ctx := context.Background()

	for orgName, org := range g.Orgs {
		// todo: also remove all secrets and vars for all repos not present
		// in the YAML config
		for repoName, repo := range org.Repos {
			// fetch encryption key
			logrus.Infof("retrieving public key for repo '%s/%s'...", orgName, repoName)
			pKey, err := pKeyProvider.GetPublicKey(ctx, orgName, repoName)
			if err == nil {
				err = syncSecrets(ctx, sService, provider, pKey, orgName, repoName, repo.Actions.Secrets)
			}
			if err != nil {
				return err
			}
			logrus.Infof("secrets synced for %s/%s\n", orgName, repoName)

			err = syncVariables(ctx, vService, orgName, repoName, repo.Actions.Variables)
			if err != nil {
				return err
			}
			logrus.Infof("variables synced for %s/%s\n", orgName, repoName)
		}
	}

	return nil
}
