// Package secrets provides a generic interface to obtain secrets from different sources.
package secrets

import (
	config "github.com/tommzn/go-config"
)

// NewSecretsManager returns a new default secrets mananger, which will read secrets from environment variables.
func NewSecretsManager() SecretsManager {
	return &EnvironmentSecretsManager{}
}

// NewStaticSecretsManager returns a secrets manager which contains passed secrets. Useful e.g. for testing.
func NewStaticSecretsManager(secrets map[string]string) SecretsManager {
	return &StaticSecretsManager{secrets: secrets}
}

// NewFileSecretsManager returns a new secretsmanager for given file.
func NewFileSecretsManager(fileName string) SecretsManager {
	return &FileSecretsManager{
		secretsFile: fileName,
	}
}

// NewDockerecretsManager returns a new secrets manager for Docker or K8s.
func NewDockerecretsManager(secretsPath string) SecretsManager {
	return &DockerSecretsManager{secretsPath: secretsPath}
}

// NewSecretsManagerByConfig will create a new secrets manager by given config.
// If there's no config values for secrets, a default secrets manager will be returned.
func NewSecretsManagerByConfig(conf config.Config) SecretsManager {

	if managerType := conf.Get("secrets.source", nil); managerType != nil {
		if *managerType == "docker" {
			secretsPath := conf.Get("secrets.path", config.AsStringPtr(DOCKER_SECRETS_PATH))
			return NewDockerecretsManager(*secretsPath)
		}
	}
	return NewSecretsManager()
}
