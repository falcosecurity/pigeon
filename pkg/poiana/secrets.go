package poiana

import (
	"github.com/1Password/connect-sdk-go/connect"
)

type SecretsProvider interface {
	GetSecret(string) (string, error)
}

type OnePasswordProvider struct {
	connect.Client
	vault string
}

type Credential struct {
	Username string `opfield:"username"`
	Password string `opfield:"password"`
}

func newOnePasswordProvider(vault string) (*OnePasswordProvider, error) {
	client, err := connect.NewClientFromEnvironment()
	if err != nil {
		return nil, err
	}
	return &OnePasswordProvider{
		client,
		vault,
	}, nil
}

func (pp *OnePasswordProvider) GetSecret(key string) (string, error) {
	var cred Credential
	err := pp.Client.LoadStructFromItemByTitle(&cred, key, pp.vault)
	if err != nil {
		return "", nil
	}
	return cred.Password, nil
}
