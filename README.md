[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/go-config.svg)](https://pkg.go.dev/github.com/tommzn/go-secrets)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/go-secrets)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tommzn/go-secrets)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/go-secrets)](https://goreportcard.com/report/github.com/tommzn/go-secrets)
[![Actions Status](https://github.com/tommzn/go-secrets/actions/workflows/go.pkg.auto-ci.yml/badge.svg)](https://github.com/tommzn/go-secrets/actions)

# Secrets Manager
Provides a generic interface to obtain secrets from different sources.

## Sources
Following sources are available:
- Static managed secrets, e.g. for testing
- Secrets read from environment variables
- Secrets read from mounted secret files in Docker or K8s 


