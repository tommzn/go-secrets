[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/go-secrets.svg)](https://pkg.go.dev/github.com/tommzn/go-secrets)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/go-secrets)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tommzn/go-secrets)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/go-secrets)](https://goreportcard.com/report/github.com/tommzn/go-secrets)
[![Actions Status](https://github.com/tommzn/go-secrets/actions/workflows/go.pkg.auto-ci.yml/badge.svg)](https://github.com/tommzn/go-secrets/actions)

# go-secrets

`go-secrets` provides a unified Go interface to read secrets from multiple backends. All backends implement the same `SecretsManager` interface, so the calling code does not need to know where a secret comes from.

## Sources

| Backend | Constructor | Description |
|---|---|---|
| Environment variables | `NewSecretsManager()` | Reads secrets from OS environment variables (default). |
| Static / in-memory | `NewStaticSecretsManager(map)` | Holds a fixed map of secrets. Useful for tests. |
| Credentials file | `NewFileSecretsManager(path)` | Reads `key:base64value` pairs from a credentials file. |
| Docker / K8s mounts | `NewDockerSecretsManager(path)` | Reads secrets from files mounted by Docker or Kubernetes. |

All backends perform **case-insensitive key lookup**: if a secret is not found under the exact key, the lower-case and upper-case variants are tried automatically.

## Installation

```bash
go get github.com/tommzn/go-secrets
```

## Usage

### Environment variables (default)

```go
manager := secrets.NewSecretsManager()
value, err := manager.Obtain("DB_PASSWORD")
if err != nil {
    log.Fatal(err)
}
fmt.Println(*value)
```

### Static secrets (e.g. for testing)

```go
manager := secrets.NewStaticSecretsManager(map[string]string{
    "API_KEY": "supersecret",
})
value, _ := manager.Obtain("api_key") // case-insensitive
```

### Credentials file

The file must contain one `key:base64encodedvalue` entry per line:

```
DB_PASSWORD:c3VwZXJzZWNyZXQ=
API_KEY:bXlhcGlrZXk=
```

```go
manager := secrets.NewFileSecretsManager("~/.credentials")
value, err := manager.Obtain("DB_PASSWORD")
```

### Docker / Kubernetes mounted secrets

```go
manager := secrets.NewDockerSecretsManager("/run/secrets")
value, err := manager.Obtain("db-password")
```

The default path `/run/secrets` is also available as the constant `secrets.DOCKER_SECRETS_PATH`.

### Config-driven creation

```go
// config value secrets.source = "docker" creates a DockerSecretsManager,
// everything else (or no value) falls back to the environment manager.
manager := secrets.NewSecretsManagerByConfig(conf)
```

### Exporting secrets to environment variables

```go
secrets.ExportToEnvironment([]string{"DB_PASSWORD", "API_KEY"}, manager)
```

## Interface

```go
type SecretsManager interface {
    Obtain(key string) (*string, error)
}
```
