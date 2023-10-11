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

package pigeon

import "github.com/google/go-github/v50/github"

type clientServiceFactory struct {
	client *github.Client
}

func NewClientServiceFactory(client *github.Client) ServiceFactory {
	return &clientServiceFactory{client: client}
}

func (c *clientServiceFactory) NewOrgVariableService(org string) ActionsVarsService {
	return NewClientOrgVariableService(c.client, org)
}

func (c *clientServiceFactory) NewOrgSecretService(org string) ActionsSecretsService {
	return NewClientOrgSecretService(c.client, org)
}

func (c *clientServiceFactory) NewRepoVariableService(org, repo string) ActionsVarsService {
	return NewClientRepoVariableService(c.client, org, repo)
}

func (c *clientServiceFactory) NewRepoSecretService(org, repo string) ActionsSecretsService {
	return NewClientRepoSecretService(c.client, org, repo)
}
