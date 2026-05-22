

a CLI to interact with the Honeycomb API.

Primary focus is easy usability for a coding agent or machine use-case.

## Commands

Use the standard Make targets shared across these Go CLI repos:

- Format: `make fmt`
- Lint: `make lint`
- Test: `make test`
- Build: `make build`
- Check: `make check`
- Install: `make install`

`make build` formats first. `make check` runs lint, test, then build.

## Project Structure

- `main.go` — CLI entry point; creates the root urfave/cli v3 app and registers subcommands from `cmd/`
- `cmd/` — CLI command definitions, flag handling, credential resolution, and JSON output helpers
- `api/` — Honeycomb API client, request/response types, and one client method file per resource
- `tests/` — CLI smoke/integration tests that build and execute the binary as a subprocess
- `timefmt/` — shared parsing helpers for human-friendly time ranges and timestamps
- `version.go` — version resolution from linker-injected value or Go build info
- API docs: https://api-docs.honeycomb.io/api

## Code Style

- Wrap errors with `fmt.Errorf("context: %w", err)` for context
- API methods take `context.Context` as first param, return `(*T, error)`
- One file per API resource in `api/` (types + client method together)
- Tests use `_test` package; smoke tests skip via `t.Skip` when env var missing
- Tests invoke the compiled binary as a subprocess (`os/exec`), not the Go API directly; `TestMain` in `cli_test.go` builds the binary once; shared helpers (`runCLI`, `runCLIWithKey`, `parseJSON`) live there too
- One test file per resource (e.g. `boards_cli_test.go`, `columns_cli_test.go`)
- JSON output via `json.NewEncoder(os.Stdout)` with 2-space indent
- Auth via `X-Honeycomb-Team` header; profiles are the preferred local UX (`hccli auth login`, `hccli auth switch`), while `--api-key` and `HONEYCOMB_API_KEY` remain supported for CI/one-off automation
