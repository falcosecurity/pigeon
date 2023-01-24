package poiana

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
func NewOnePasswordSecretsProvider() (*onePasswordSecretsProvider, error) {
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
