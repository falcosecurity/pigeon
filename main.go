package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/FedeDP/GhEnvSet/pkg/poiana"
	"github.com/google/go-github/v49/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/jamesruan/sodium"
)

func fail(err string) {
	logrus.Fatal(err)
	os.Exit(1)
}

const sampleYAML = `
orgs:
  FedeDP:
    repos:
      GhEnvSet:
        actions:
          variables:
            SOME_VARIABLE2: "ciao"
          secrets:
            - TEST_SECRET_KEY
`

func syncSecrets(ctx context.Context, client *poiana.Client, provider poiana.SecretsProvider, orgName, repoName string, secrets []string) error {
	// Step 1: load repo secrets
	logrus.Infof("listing secrets for repo '%s/%s'...", orgName, repoName)
	secs, _, err := client.Actions.ListRepoSecrets(ctx, orgName, repoName, nil)
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
			_, err = client.Actions.DeleteRepoSecret(ctx, orgName, repoName, existentSec.Name)
			if err != nil {
				return err
			}
		}
	}

	// Step 3: fetch encryption key
	logrus.Infof("retrieving public key for repo '%s/%s'...", orgName, repoName)
	pKey, _, err := client.Actions.GetRepoPublicKey(ctx, orgName, repoName)
	if err != nil {
		return err
	}

	keyBytes, err := base64.StdEncoding.DecodeString(pKey.GetKey())
	if err != nil {
		return err
	}

	// Step 4: add or update all conf-listed secrets
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
		_, err = client.Actions.CreateOrUpdateRepoSecret(ctx, orgName, repoName, &github.EncryptedSecret{
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

func syncVariables(ctx context.Context, client *poiana.Client, orgName, repoName string, variables map[string]string) error {
	// Step 1: load repo variables
	logrus.Infof("listing variables for repo '%s/%s'...", orgName, repoName)
	vars, _, err := client.Actions.ListRepoVariables(ctx, orgName, repoName, nil)
	if err != nil {
		return err
	}

	// Step 2: delete all variables that are no more existing
	for _, existentVar := range vars.Variables {
		_, ok := variables[existentVar.Name]
		if !ok {
			logrus.Infof("deleting variable '%s' for repo '%s/%s'...", existentVar.Name, orgName, repoName)
			_, err = client.Actions.DeleteRepoVariable(ctx, orgName, repoName, existentVar.Name)
			if err != nil {
				return err
			}
		}
	}

	// Step 3: add or update all conf-listed variables
	for newVarName, newVarValue := range variables {
		logrus.Infof("adding/updating variable '%s' in repo '%s/%s'...", newVarName, orgName, repoName)
		_, err = client.Actions.CreateOrUpdateRepoVariable(ctx, orgName, repoName, &poiana.Variable{
			Name:  newVarName,
			Value: newVarValue,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func getClient(ctx context.Context) *poiana.Client {
	ghTok := os.Getenv("GH_TOKEN")
	var tc *http.Client
	if ghTok != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: ghTok},
		)
		tc = oauth2.NewClient(ctx, ts)
	}
	return poiana.NewClient(github.NewClient(tc))
}

func main() {
	conf := poiana.GithubConfig{}
	err := conf.Decode(strings.NewReader(sampleYAML))
	if err != nil {
		fail(err.Error())
	}

	provider, _ := poiana.NewMockSecretsProvider(map[string]string{
		"TEST_SECRET_KEY2": "ciaociaociao2",
		"TEST_SECRET_KEY":  "ciaociaociao",
	})

	ctx := context.Background()
	client := getClient(ctx)
	for orgName, org := range conf.Orgs {
		// todo: also remove all secrets and vars for all repos not present
		// in the YAML config
		for repoName, repo := range org.Repos {
			err = syncSecrets(ctx, client, provider, orgName, repoName, repo.Actions.Secrets)
			if err != nil {
				fail(err.Error())
			}
			logrus.Infof("secrets synced for %s/%s\n", orgName, repoName)

			err = syncVariables(ctx, client, orgName, repoName, repo.Actions.Variables)
			if err != nil {
				fail(err.Error())
			}
			logrus.Infof("variables synced for %s/%s\n", orgName, repoName)
		}
	}
}
