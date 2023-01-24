package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/FedeDP/GhEnvSet/pkg/poiana"
	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"
)

func fail(err string) {
	println(err)
	os.Exit(1)
}

const sampleYAML = `
orgs:
  falcosecurity:
    repos:
      .github:
        actions:
          variables:
            GH_TOKEN: "ciao"
          secrets:
            - S3_ACCESSEKY
`

func syncSecrets(ctx context.Context, client *poiana.Client, orgName, repoName string, secrets []string) error {
	// Step 1: load repo secrets
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
			_, err = client.Actions.DeleteRepoSecret(ctx, orgName, repoName, existentSec.Name)
			if err != nil {
				return err
			}
		}
	}

	// Step 3: fetch encryption key
	pKey, _, err := client.Actions.GetRepoPublicKey(ctx, orgName, repoName)
	if err != nil {
		return err
	}

	// Step 4: add or update all conf-listed secrets
	for _, sec := range secrets {
		// TODO
		// encrypted_value string
		//	Value for your secret, encrypted with LibSodium using the public key retrieved from the Get a repository public key endpoint.
		//
		//	key_id string
		//	ID of the key you used to encrypt the secret.
		_, err = client.Actions.CreateOrUpdateRepoSecret(ctx, orgName, repoName, &github.EncryptedSecret{
			Name:           sec,
			KeyID:          pKey.GetKeyID(),
			EncryptedValue: "", // TODO: fetch from 1password and encrypt using pkey
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func syncVariables(ctx context.Context, client *poiana.Client, orgName, repoName string, variables map[string]string) error {
	// Step 1: load repo variables
	vars, _, err := client.Actions.ListRepoVariables(ctx, orgName, repoName, nil)
	if err != nil {
		return err
	}

	// Step 2: delete all variables that are no more existing
	found := false
	for _, existentVar := range vars.Variables {
		for newVarName, _ := range variables {
			if newVarName == existentVar.Name {
				found = true
				break
			}
		}
		if !found {
			_, err = client.Actions.DeleteRepoVariable(ctx, orgName, repoName, existentVar.Name)
			if err != nil {
				return err
			}
		}
	}

	// Step 3: add or update all conf-listed variables
	for newVarName, newVarValue := range variables {
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

func main() {
	conf := poiana.GithubConfig{}
	err := conf.Decode(strings.NewReader(sampleYAML))
	if err != nil {
		fail(err.Error())
	}

	ctx := context.Background()
	fmt.Println("Start")

	ghTok := os.Getenv("GH_TOKEN")
	var tc *http.Client
	if ghTok != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: ghTok},
		)
		tc = oauth2.NewClient(ctx, ts)
	}
	client := poiana.NewClient(github.NewClient(tc))

	fmt.Println("Looping orgs")
	for orgName, org := range conf.Orgs {
		for repoName, repo := range org.Repos {
			err = syncSecrets(ctx, client, orgName, repoName, repo.Actions.Secrets)
			if err != nil {
				fail(err.Error())
			}
			fmt.Printf("Secrets synced for %s/%s\n", orgName, repoName)

			err = syncVariables(ctx, client, orgName, repoName, repo.Actions.Variables)
			if err != nil {
				fail(err.Error())
			}
			fmt.Printf("Variables synced for %s/%s\n", orgName, repoName)
		}
	}
	fmt.Println("End")
}
