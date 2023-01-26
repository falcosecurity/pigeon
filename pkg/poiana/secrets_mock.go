package poiana

import "fmt"

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
