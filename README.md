[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/go-config.svg)](https://pkg.go.dev/github.com/tommzn/go-secrets)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/go-secrets)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tommzn/go-secrets)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/go-secrets)](https://goreportcard.com/report/github.com/tommzn/go-secrets)
[![Actions Status](https://github.com/tommzn/go-secrets/actions/workflows/go.pkg.auto-ci.yml/badge.svg)](https://github.com/tommzn/go-secrets/actions)

# go-secrets

A Go library that provides a unified interface for reading secrets from multiple backends. Switch between environment variables, Docker/Kubernetes mounted secrets, and credentials files without changing application code.

## Installation

```bash
go get github.com/tommzn/go-secrets
```

## Backends

| Backend | Constructor | Use case |
|---------|-------------|----------|
| Environment variables | `NewSecretsManager()` | Default; twelve-factor apps |
| Docker / Kubernetes | `NewDockerSecretsManager(path)` | Mounted secret files in containers |
| Credentials file | `NewFileSecretsManager(file)` | Local development, CI runners |
| Static (in-memory) | `NewStaticSecretsManager(map)` | Unit tests |

All backends implement the same interface:

```go
type SecretsManager interface {
    Obtain(key string) (*string, error)
}
```

`Obtain` returns a pointer to the secret value, or a `*SecretNotFoundError` when the key does not exist.

## Usage

### Environment variables

```go
manager := secrets.NewSecretsManager()

value, err := manager.Obtain("DATABASE_URL")
if err != nil {
    // handle *secrets.SecretNotFoundError
}
fmt.Println(*value)
```

Key lookup is case-insensitive: `Obtain("db_password")` also matches `DB_PASSWORD` and `Db_Password` (the environment manager tries the original key, lowercase, and uppercase variants).

### Docker / Kubernetes mounted secrets

```go
// Use the default mount path /run/secrets
manager := secrets.NewDockerSecretsManager(secrets.DOCKER_SECRETS_PATH)

// Or a custom path
manager = secrets.NewDockerSecretsManager("/var/run/secrets/myapp")

value, err := manager.Obtain("db-password")
```

Each secret must be a separate file inside the secrets directory. The filename is the key. Key lookup is case-insensitive (the original key, lowercase, and uppercase variants are tried).

### Credentials file

```go
// ~ is expanded to the current user's home directory
manager := secrets.NewFileSecretsManager("~/.credentials")

value, err := manager.Obtain("db_password")
```

**Credentials file format** — one entry per line, key and base64-encoded value separated by `:`:

```
db_password:cGFzc3dvcmQxMjM=
api_key:c2VjcmV0a2V5eHl6
```

Encode a value:

```bash
echo -n "mysecretvalue" | base64
```

**Security requirements for the credentials file:**
- Permissions must be `0600` or stricter (no group or world read/write/execute; owner must have read access)
- The path must not be a symlink
- Key lookup in credentials files is **case-sensitive** — the key in the file must match exactly.

### Static (for tests)

```go
manager := secrets.NewStaticSecretsManager(map[string]string{
    "API_KEY": "test-key-123",
    "DB_URL":  "postgres://localhost/testdb",
})

value, err := manager.Obtain("API_KEY")

// Remove all secrets from memory when done
manager.(*secrets.StaticSecretsManager).Clear()
```

### Export secrets to environment variables

`ExportToEnvironment` loads a list of secrets by key and sets each one as an environment variable. It returns the keys that could not be looked up or set, so callers can decide whether a missing secret is fatal.

```go
manager := secrets.NewSecretsManager()
failed := secrets.ExportToEnvironment([]string{"API_KEY", "DB_PASSWORD"}, manager)
if len(failed) > 0 {
    // handle missing secrets, e.g. log.Fatal or fall back to defaults
}
```

### Select backend from configuration

The library integrates with [go-config](https://github.com/tommzn/go-config). Set `secrets.source` in your config file to select a backend:

```yaml
secrets:
  source: docker
  path: /run/secrets
```

```go
conf, _ := config.NewFileConfigSource(&configFile).Load()
manager := secrets.NewSecretsManagerByConfig(conf)
```

Supported values for `secrets.source`:

| Value | Backend |
|-------|---------|
| `docker` | `DockerSecretsManager` using path from `secrets.path` (default: `/run/secrets`) |
| _(absent)_ | `EnvironmentSecretsManager` |

## Error handling

```go
value, err := manager.Obtain("MISSING_KEY")
if err != nil {
    var notFound *secrets.SecretNotFoundError
    if errors.As(err, &notFound) {
        // key does not exist in this backend
    }
}
```

## Security notes

- **Credentials file**: rejected if it has group or world permissions, or if the path resolves to a symlink.
- **Docker backend**: key names containing path separators (`/`, `\`) or traversal sequences (`.`, `..`) are rejected to prevent directory traversal attacks.
- **Error messages**: secret key names are never included in error messages to avoid leaking them into logs.
- **Base64 in credentials files** is encoding, not encryption. Protect the file with filesystem permissions and avoid committing it to version control.

## License

MIT
