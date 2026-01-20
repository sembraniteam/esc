# Repository Guidelines

## Project Structure & Module Organization
- `cmd/esc/main.go` is the CLI entry point.
- `internal/cli/` holds command handlers wired into the CLI.
- `internal/config/` contains configuration loading, schema, validation, and writing.
- `internal/sshx/` includes SSH client/auth helpers.
- `internal/filex/` provides filesystem utilities.
- Dependencies and Go version live in `go.mod`/`go.sum` (Go 1.25.5 in `go.mod`).

## Build, Test, and Development Commands
- `go build ./cmd/esc` builds the `esc` CLI binary.
- `go run ./cmd/esc` runs the CLI directly from source.
- `go test ./...` runs all Go tests (none present yet).
- `go vet ./...` runs static checks for common issues.

## Coding Style & Naming Conventions
- Follow standard Go formatting: tabs for indentation and `gofmt` output.
- Package names are short, lowercase, and descriptive (e.g., `cli`, `config`).
- File names are lowercase with underscores when needed (e.g., `load.go`).
- Prefer explicit, error-returning functions over panics in library code.
- Use `cobra`-style command naming if adding CLI commands (see `internal/cli`).

## Testing Guidelines
- No `_test.go` files are present; add tests alongside packages when introducing logic.
- Use Go’s `testing` package; name files `*_test.go` and test funcs `TestXxx`.
- Run `go test ./...` before opening a PR.

## Commit & Pull Request Guidelines
- Git history is not available in this workspace; no formal commit style is enforced.
- Suggested commit style: short, imperative subject (e.g., `Add config validation`).
- PRs should include: purpose, approach, and any CLI output or examples of use.
- Link related issues if applicable; call out any user-visible behavior changes.

## Configuration & Security Notes
- Configuration logic lives under `internal/config/`; keep validation centralized there.
- Avoid storing secrets in config defaults; read secrets from environment when needed.
