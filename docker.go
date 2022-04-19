package secrets

import (
	"fmt"
	"io/ioutil"
)

// DOCKER_SECRETS_PATH defined the default path to look for mounted secrets in Docker or K8s.
const DOCKER_SECRETS_PATH = "/run/secrets"

// DockerSecretsManager will read secrets from files mounted by Docker or K8s.
type DockerSecretsManager struct {
	secretsPath string
}

// Obtain will try to read secrets from mountes secrets files.
func (s *DockerSecretsManager) Obtain(key string) (*string, error) {

	keys := generateSecretKeys(key)
	for _, currentKey := range keys {
		fullPath := generateSecretFilePath(s.secretsPath, currentKey)
		if secret, err := ioutil.ReadFile(fullPath); err == nil {
			secretStr := string(secret)
			return &secretStr, nil
		}
	}
	return nil, asSecretNotFoundError(key)
}

// generateSecretFilePath creates the path to a mounted secrets file.
func generateSecretFilePath(path, filename string) string {
	return fmt.Sprintf("%s/%s", path, filename)
}
