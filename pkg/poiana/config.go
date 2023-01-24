package poiana

import (
	"io"

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
