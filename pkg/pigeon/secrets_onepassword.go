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

import (
	"fmt"
	"os"

	"github.com/1Password/connect-sdk-go/connect"
)

var (
	// todo: switch to "credential", "notes", or do a cascade attempt
	onePasswordItemFieldName = "password"
)

type onePasswordSecretsProvider struct {
	c connect.Client
}

// NewOnePasswordSecretsProvider creates a secret provider that retrieves
// secrets from 1Password. It requires three environment variables to be set:
//   - `OP_CONNECT_TOKEN`: the API token to be used to authenticate the client
//     to your 1Password Connect instance.
//   - `OP_CONNECT_HOST`: the hostname of your 1Password Connect instance.
//   - `OP_VAULT`: a vault UUID from which retrieving the secrets.
func NewOnePasswordSecretsProvider() (SecretsProvider, error) {
	client, err := connect.NewClientFromEnvironment()
	if err != nil {
		return nil, err
	}
	return &onePasswordSecretsProvider{client}, nil
}

func (pp *onePasswordSecretsProvider) GetSecret(key string) (string, error) {
	vault, vaultFount := os.LookupEnv("OP_VAULT")
	if !vaultFount {
		return "", fmt.Errorf("the OP_VAULT env variable is not set")
	}
	item, err := pp.c.GetItemByTitle(key, vault)
	if err != nil {
		return "", nil
	}
	return item.GetValue(onePasswordItemFieldName), nil
}
