package secrets

import "fmt"

// SecretNotFoundError is returned when a requested secret key does not exist
// in the underlying secrets source.
type SecretNotFoundError struct {
	key string
}

// Error implements the error interface.
func (e *SecretNotFoundError) Error() string {
	return fmt.Sprintf("secret not found: %q", e.key)
}

// Key returns the key that was not found.
func (e *SecretNotFoundError) Key() string {
	return e.key
}

// asSecretNotFoundError returns a SecretNotFoundError for the given key.
func asSecretNotFoundError(key string) error {
	return &SecretNotFoundError{key: key}
}
