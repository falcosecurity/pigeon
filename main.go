package main

import (
	"context"
	"net/http"
	"os"

	"github.com/FedeDP/GhEnvSet/pkg/poiana"
	"github.com/google/go-github/v49/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
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
	conf, err := poiana.FromData(sampleYAML)
	if err != nil {
		fail(err.Error())
	}

	provider, _ := poiana.NewMockSecretsProvider(map[string]string{
		"TEST_SECRET_KEY2": "ciaociaociao2",
		"TEST_SECRET_KEY":  "ciaociaociao",
	})

	ctx := context.Background()
	client := getClient(ctx)

	err = conf.Loop(client.Actions, client.Actions, provider, client.Actions)
	if err != nil {
		fail(err.Error())
	}
}
