package secrets

// StaticSecretsManager will manage secrets by internal map.
type StaticSecretsManager struct {

	// secrets contains a key/value list of all managed secrets.
	secrets map[string]string
}

// Obtain will try to get a secrets from internal map by given key.
func (s *StaticSecretsManager) Obtain(key string) (*string, error) {

	keys := generateSecretKeys(key)
	for _, currentKey := range keys {
		if secret, ok := s.secrets[currentKey]; ok {
			return &secret, nil
		}
	}
	return nil, newSecretNotFoundError(key)
}
