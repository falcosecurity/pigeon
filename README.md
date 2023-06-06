# Pigeon 

[![Falco Infra Repository](https://github.com/falcosecurity/evolution/blob/main/repos/badges/falco-infra-blue.svg)](https://github.com/falcosecurity/evolution/blob/main/REPOSITORIES.md#infra-scope) [![Incubating](https://img.shields.io/badge/status-incubating-orange?style=for-the-badge)](https://github.com/falcosecurity/evolution/blob/main/REPOSITORIES.md#incubating)

[![Latest](https://img.shields.io/github/v/release/falcosecurity/pigeon)](https://github.com/falcosecurity/pigeon/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/falcosecurity/pigeon)](https://goreportcard.com/report/github.com/falcosecurity/pigeon)
[![CI build](https://github.com/falcosecurity/pigeon/actions/workflows/ci.yaml/badge.svg)](https://github.com/falcosecurity/pigeon/actions/workflows/ci.yaml)

Script to set github org/repo actions variables and secrets from a config file.  

## Cli options

* `--conf` -> **MANDATORY**: yaml config for the run
* `--gh-token` -> path to github token file
* `--dry-run` -> don't actually do any change
* `--verbose` -> enable verbose logging

Github token file can also be passed from `GITHUB_TOKEN_FILE` env variable.    
> **NOTE:** github token file **MUST** be set from either the cli flag or env variable!

## Example Config

```yaml
orgs:
  foo:
    actions:
      variables:
        orgVar1: "orgValue1"
      secrets:
        - orgSecret0
    repos:
      bar:
        actions:
          variables:
            repoVar1: "repoValue1"
            repoVar2: "repoValue2"
          secrets:
            - repoSecret0
```
