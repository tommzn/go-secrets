package secrets

import (
	"os"
	"path/filepath"
	"strings"
)

// DOCKER_SECRETS_PATH defines the default path to look for mounted secrets in Docker or K8s.
const DOCKER_SECRETS_PATH = "/run/secrets"

// DockerSecretsManager will read secrets from files mounted by Docker or K8s.
type DockerSecretsManager struct {
	secretsPath string
}

// Obtain will try to read secrets from mounted secrets files.
func (s *DockerSecretsManager) Obtain(key string) (*string, error) {

	keys := generateSecretKeys(key)
	for _, currentKey := range keys {
		if !isValidSecretFileName(currentKey) {
			continue
		}
		fullPath := generateSecretFilePath(s.secretsPath, currentKey)
		if secret, err := os.ReadFile(fullPath); err == nil {
			secretStr := strings.TrimRight(string(secret), "\r\n")
			return &secretStr, nil
		}
	}
	return nil, asSecretNotFoundError(key)
}

// generateSecretFilePath creates the path to a mounted secrets file.
func generateSecretFilePath(path, filename string) string {
	return filepath.Join(path, filename)
}

// isValidSecretFileName rejects filenames containing path separators or traversal sequences.
func isValidSecretFileName(filename string) bool {
	return !strings.Contains(filename, "/") &&
		!strings.Contains(filename, "\\") &&
		filename != "." &&
		filename != ".."
}
