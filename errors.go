package secrets

// SecretNotFoundError is returned when a requested secret key does not exist
// in the underlying secrets source.
type SecretNotFoundError struct {
	key string
}

// Error implements the error interface.
func (e *SecretNotFoundError) Error() string {
	return "secret not found"
}

// asSecretNotFoundError returns a SecretNotFoundError for the given key.
func asSecretNotFoundError(key string) error {
	return &SecretNotFoundError{key: key}
}
