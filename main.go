package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

func main() {
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
	client := github.NewClient(tc)

	fmt.Println("Have client")
	env, _, err := client.Repositories.CreateUpdateEnvironment(ctx, "FedeDP", "test-infra", "stocazzo2",
		&github.CreateUpdateEnvironment{
			WaitTimer:              nil,
			Reviewers:              nil,
			DeploymentBranchPolicy: nil,
		})
	if err == nil {
		fmt.Println(env.CreatedAt)
	}

	client.Actions.CreateOrUpdateEnvSecret(ctx, "FedeDP", "stocazzo2")
	repos, _, err := client.Repositories.List(ctx, "falcosecurity", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Have repos")
	for i, repo := range repos {
		fmt.Printf("%d) %s\n", i, *repo.Name)
	}
	fmt.Println("End")
}
