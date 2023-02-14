package main

import (
	"context"
	"flag"
	"github.com/FedeDP/Pigeon/pkg/pigeon"
	"github.com/google/go-github/v50/github"
	"os"
	"strings"

	"github.com/FedeDP/Pigeon/pkg/config"
	"github.com/sirupsen/logrus"
)

var (
	ghToken  string
	confFile string
	dryRun   bool
	verbose  bool
)

func init() {
	flag.StringVar(&confFile, "conf", "", "path to yaml conf file")
	flag.StringVar(&ghToken, "gh-token", "", "path to github token with admin access on org/repo")
	flag.BoolVar(&dryRun, "dry-run", false, "enable dry run mode")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
}

func initOpts() {
	flag.Parse()

	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Load token file from cli flag or env
	if ghToken == "" {
		ghToken = os.Getenv("GITHUB_TOKEN_FILE")
		if ghToken == "" {
			logrus.Fatal(`Github token must be provided either through "gh-token" flag, or "GITHUB_AUTH_TOKEN" env."`)
		}
	}
	ghTokBytes, err := os.ReadFile(ghToken)
	if err != nil {
		logrus.Fatal(err)
	}
	ghToken = string(ghTokBytes)
	ghToken = strings.Trim(ghToken, "\n")

	if confFile == "" {
		logrus.Fatal(`"conf" flag must be set`)
	}
}

func main() {
	initOpts()

	conf, err := config.FromFile(confFile)
	if err != nil {
		logrus.Fatal(err)
	}

	provider, err := pigeon.NewOnePasswordSecretsProvider()
	if err != nil {
		logrus.Fatal(err)
	}

	ctx := context.Background()
	client := github.NewTokenClient(ctx, ghToken)
	err = conf.Sync(pigeon.NewClientServiceFactory(client), provider, dryRun)
	if err != nil {
		logrus.Fatal(err)
	}
}
