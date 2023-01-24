package poiana

// SecretsProvider retrieves secrets with a given key
type SecretsProvider interface {
	// GetSecret returns a secret with the given key.
	// Returns a non-nil error in case of failure
	GetSecret(string) (string, error)
}
