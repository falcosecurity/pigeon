package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/FedeDP/GhEnvSet/pkg/poiana"
	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"
)

func fail(err string) {
	println(err)
	os.Exit(1)
}

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
	client := poiana.NewClient(github.NewClient(tc))

	// fmt.Println("Have client")
	// env, _, err := client.Repositories.CreateUpdateEnvironment(ctx, "FedeDP", "test-infra", "stocazzo2",
	// 	&github.CreateUpdateEnvironment{
	// 		WaitTimer:              nil,
	// 		Reviewers:              nil,
	// 		DeploymentBranchPolicy: nil,
	// 	})
	// if err == nil {
	// 	fmt.Println(env.CreatedAt)
	// }

	vars, _, err := client.Actions.ListRepoVariables(ctx, "FedeDP", "test-infra", nil)
	if err != nil {
		fail(err.Error())
	}
	for _, s := range vars.Variables {
		fmt.Printf("%s: %s\n", s.Name, s.Value)
	}

	println("creating...")
	_, err = client.Actions.CreateOrUpdateRepoVariable(ctx, "FedeDP", "test-infra", &poiana.Variable{
		Name:  "topkek2",
		Value: "stocazzo2",
	})
	if err != nil {
		fail(err.Error())
	}

	vars, _, err = client.Actions.ListRepoVariables(ctx, "FedeDP", "test-infra", nil)
	if err != nil {
		fail(err.Error())
	}
	for _, s := range vars.Variables {
		fmt.Printf("  %s: %s\n", s.Name, s.Value)
	}

	println("deleting...")
	_, err = client.Actions.DeleteRepoVariable(ctx, "FedeDP", "test-infra", "topkek2")
	if err != nil {
		fail(err.Error())
	}

	vars, _, err = client.Actions.ListRepoVariables(ctx, "FedeDP", "test-infra", nil)
	if err != nil {
		fail(err.Error())
	}
	for _, s := range vars.Variables {
		fmt.Printf("%s: %s\n", s.Name, s.Value)
	}

	fmt.Println("End")
}
