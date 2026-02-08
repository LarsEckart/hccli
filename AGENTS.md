

a CLI to interact with the Honeycomb API.

Primary focus is easy usability for a coding agent or machine use-case.

## Commands

- Install: `make install` (installs hccli binary to $GOPATH/bin)
- Build: `make build`
- Test: `make test`
- Format: `make fmt` (uses goimports)
- Lint: `make lint` (uses golangci-lint v2)
- All checks: `make check` (format check + lint + test)

## Project Structure

- `main.go` — CLI entry point, subcommand wiring (urfave/cli v3)
- `api/` — Honeycomb API client and types
- API docs: https://api-docs.honeycomb.io/api

## Code Style

- Wrap errors with `fmt.Errorf("context: %w", err)` for context
- API methods take `context.Context` as first param, return `(*T, error)`
- One file per API resource in `api/` (types + client method together)
- Tests use `_test` package; smoke tests skip via `t.Skip` when env var missing
- Tests invoke the compiled binary as a subprocess (`os/exec`), not the Go API directly; `TestMain` in `cli_test.go` builds the binary once; shared helpers (`runCLI`, `runCLIWithKey`, `parseJSON`) live there too
- One test file per resource (e.g. `boards_cli_test.go`, `columns_cli_test.go`)
- JSON output via `json.NewEncoder(os.Stdout)` with 2-space indent
- Auth via `X-Honeycomb-Team` header; key from `--api-key` flag or `HONEYCOMB_API_KEY` env
