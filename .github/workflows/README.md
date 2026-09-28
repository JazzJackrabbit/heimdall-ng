# GitHub Actions

CI runs on push and pull requests against `main`. No network or API keys are
involved.

## Workflows

- `golangci-lint.yml` runs `golangci-lint`.
- `build.yml` runs `go build` and `go vet`.

## Provider tests

The `providers/*_test.go` tests make live API calls and skip when their
provider's API key environment variable is unset. Run them locally with the
keys set, for example:

```
just test-anthropic
```
