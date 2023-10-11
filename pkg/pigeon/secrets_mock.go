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
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/v50/github"
)

type mockSecretsProvider struct {
	m map[string]string
}

// NewMockSecretsProvider creates a secret provider that retrieves
// secrets from a key-value map.
func NewMockSecretsProvider(values map[string]string) (SecretsProvider, error) {
	return &mockSecretsProvider{m: values}, nil
}

func (m *mockSecretsProvider) GetSecret(key string) (string, error) {
	v, ok := m.m[key]
	if !ok {
		return "", fmt.Errorf("secret not found: %s", key)
	}
	return v, nil
}

type mockSecretsService struct {
	secrets map[string]*github.EncryptedSecret
}

func (m *mockSecretsService) GetPublicKey(ctx context.Context) (*github.PublicKey, *github.Response, error) {
	keyID := "testing"
	key := base64.StdEncoding.EncodeToString([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")) // 32B key
	pKey := github.PublicKey{
		KeyID: &keyID,
		Key:   &key,
	}
	return &pKey, nil, nil
}

func (m *mockSecretsService) ListSecrets(_ context.Context, _ *github.ListOptions) (*github.Secrets, *github.Response, error) {
	secs := make([]*github.Secret, 0)
	for key := range m.secrets {
		secs = append(secs, &github.Secret{
			Name: key,
		})
	}

	return &github.Secrets{
		TotalCount: len(m.secrets),
		Secrets:    secs,
	}, nil, nil
}

func (m *mockSecretsService) DeleteSecret(_ context.Context, name string) (*github.Response, error) {
	delete(m.secrets, name)
	return nil, nil
}

func (m *mockSecretsService) CreateOrUpdateSecret(_ context.Context, eSecret *github.EncryptedSecret) (*github.Response, error) {
	m.secrets[eSecret.Name] = eSecret
	return nil, nil
}

func NewMockSecretsService() ActionsSecretsService {
	return &mockSecretsService{
		secrets: make(map[string]*github.EncryptedSecret, 0),
	}
}
