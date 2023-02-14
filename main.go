package main

import (
	"context"
	"flag"

	"github.com/FedeDP/Pigeon/pkg/config"
	"github.com/FedeDP/Pigeon/pkg/poiana"
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
	flag.StringVar(&ghToken, "gh-token", "", "path to poiana github token with admin access on org/repo")
	flag.BoolVar(&dryRun, "dry-run", false, "enable dry run mode")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
}

func initOpts() {
	flag.Parse()

	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if ghToken == "" {
		logrus.Fatal(`"gh-token" flag must be set`)
	}
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

	provider, err := poiana.NewOnePasswordSecretsProvider()
	if err != nil {
		logrus.Fatal(err)
	}

	ctx := context.Background()
	client, err := poiana.NewClient(ctx, ghToken)
	if err != nil {
		logrus.Fatal(err)
	}

	err = conf.Sync(poiana.NewClientServiceFactory(client), provider, dryRun)
	if err != nil {
		logrus.Fatal(err)
	}
}
