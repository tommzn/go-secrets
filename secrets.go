// Package secrets provides a generic interface to obtain secrets from different sources.
package secrets

import (
	"os/user"
	"path/filepath"
	"strings"

	config "github.com/tommzn/go-config"
)

// NewSecretsManager returns a new default secrets manager, which will read secrets from environment variables.
func NewSecretsManager() SecretsManager {
	return &EnvironmentSecretsManager{}
}

// NewStaticSecretsManager returns a secrets manager which contains passed secrets. Useful e.g. for testing.
func NewStaticSecretsManager(secrets map[string]string) SecretsManager {
	return &StaticSecretsManager{secrets: secrets}
}

// NewFileSecretsManager returns a new secrets manager for the given credentials file.
// A leading ~/ in the path is expanded to the current user's home directory.
func NewFileSecretsManager(fileName string) SecretsManager {
	return &FileSecretsManager{
		secretsFile: expandHome(fileName),
	}
}

// expandHome replaces a leading ~ with the current user's home directory.
func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		if usr, err := user.Current(); err == nil {
			return filepath.Join(usr.HomeDir, path[2:])
		}
	}
	return path
}

// NewDockerSecretsManager returns a new secrets manager for Docker or K8s.
func NewDockerSecretsManager(secretsPath string) SecretsManager {
	return &DockerSecretsManager{secretsPath: secretsPath}
}

// NewSecretsManagerByConfig will create a new secrets manager based on the
// given config. If there are no config values for secrets, the default
// (environment-based) secrets manager is returned.
func NewSecretsManagerByConfig(conf config.Config) SecretsManager {

	managerType := conf.Get("secrets.source", nil)
	if managerType == nil || *managerType != "docker" {
		return NewSecretsManager()
	}
	// Get always returns the non-nil default here, so no nil-check is needed.
	secretsPath := conf.Get("secrets.path", config.AsStringPtr(DOCKER_SECRETS_PATH))
	return NewDockerSecretsManager(*secretsPath)
}
