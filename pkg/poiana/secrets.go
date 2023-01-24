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

type SecretsProvider interface {
	GetSecret(string) (string, error)
}

type onePasswordProvider struct {
	c connect.Client
}

func NewOnePasswordProvider() (*onePasswordProvider, error) {
	client, err := connect.NewClientFromEnvironment()
	if err != nil {
		return nil, err
	}
	return &onePasswordProvider{client}, nil
}

func (pp *onePasswordProvider) GetSecret(key string) (string, error) {
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
