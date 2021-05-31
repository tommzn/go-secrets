package secrets

import "os"

// EnvironmentSecretsManager will read secrets from environment variables.
type EnvironmentSecretsManager struct {
}

// Obtain will try to read a secrets from environment variable defined by passed key.
func (s *EnvironmentSecretsManager) Obtain(key string) (*string, error) {

	keys := generateSecretKeys(key)
	for _, currentKey := range keys {
		if secret, ok := os.LookupEnv(currentKey); ok {
			return &secret, nil
		}
	}
	return nil, newSecretNotFoundError(key)
}
