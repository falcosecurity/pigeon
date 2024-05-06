// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"github.com/pkg/errors"
	"os"
	"strings"
	"unicode"

	"github.com/falcosecurity/pigeon/pkg/pigeon"
	"github.com/google/go-github/v50/github"

	"github.com/falcosecurity/pigeon/pkg/config"
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

func getTokenFromFile(path string) (string, error) {
	token, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "error reading token file")
	}

	return removeNonPrintableChars(string(token)), nil
}

func removeNonPrintableChars(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case unicode.IsPrint(r):
			return r
		default:
			return -1
		}
	}, s)
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
	var err error
	ghToken, err = getTokenFromFile(ghToken)
	if err != nil {
		logrus.Fatal(err)
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
