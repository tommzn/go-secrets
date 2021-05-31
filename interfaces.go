package secrets

// SecretsManager is a generic interface to read secrets from different sources.
type SecretsManager interface {

	// Obtain will try to read a secret for given key.
	Obtain(key string) (*string, error)
}
