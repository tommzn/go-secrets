package secrets

type SecretNotFoundError struct {
	key string
}

func (e *SecretNotFoundError) Error() string {
	return "Secret not found: " + e.key
}

type Base64DecodeError struct {
	msg string
}

func (e *Base64DecodeError) Error() string {
	return e.msg
}

// asSecretNotFoundError is a helper to provide same error for not existing secrets across all secret sources.
func asSecretNotFoundError(key string) error {
	return &SecretNotFoundError{key: key}
}

// asBase64DecodeError creates an error for base64 decoding failures.
func asBase64DecodeError(err error) error {
	if err == nil {
		return nil
	}
	return &Base64DecodeError{msg: err.Error()}
}
