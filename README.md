[![Certified Shovelware](https://justin.searls.co/img/shovelware.svg)](https://justin.searls.co/shovelware/)

# hccli

A machine-friendly CLI for [Honeycomb](https://www.honeycomb.io/) observability.

## Installation

### From source

```bash
make install
```

### Download a release binary

Pre-built binaries for Linux, macOS, and Windows (amd64/arm64) are available on the [Releases](https://github.com/LarsEckart/hccli/releases) page.

## Authentication

`hccli` supports named Honeycomb profiles, similar to `gh auth switch`. This is useful when you use separate Honeycomb accounts for work and personal projects.

```bash
printf '%s\n' "$WORK_HONEYCOMB_API_KEY" | hccli auth login --profile work --api-key-stdin
printf '%s\n' "$PERSONAL_HONEYCOMB_API_KEY" | hccli auth login --profile personal --api-key-stdin

hccli auth list
hccli auth switch work
hccli auth status
hccli --profile personal boards
```

Profiles are stored in `~/.config/hccli/config.json` on Linux/macOS (`os.UserConfigDir` on each platform) with `0600` permissions. The file contains API keys, so do not commit or share it.

Credential precedence is:

1. `--api-key`
2. `HONEYCOMB_API_KEY`
3. `--profile`
4. `HCCLI_PROFILE`
5. current project's local profile from `hccli auth switch <profile> --local`
6. global active profile from `hccli auth switch <profile>`

For CI and one-off use, you can still provide your [Honeycomb API key](https://docs.honeycomb.io/get-started/configure/environments/manage-api-keys/) directly:

```bash
export HONEYCOMB_API_KEY=your-key-here
hccli auth whoami
```

### Auth commands

```bash
hccli auth login --profile work --api-key-stdin   # store or update a profile
hccli auth list                                   # show profiles without revealing keys
hccli auth switch work                            # set global active profile
hccli auth switch work --local                    # set active profile for this project
hccli auth status                                 # show and validate active credentials
hccli auth whoami                                 # raw /1/auth response
hccli auth whoami-v2                              # raw /2/auth response for management keys
hccli auth logout work                            # remove a profile
```

Use `HCCLI_CONFIG_DIR` to override the profile storage directory, primarily for tests and sandboxed agent runs.

## Commands

Run `hccli --help` for full command reference.

## Large Output

When JSON output exceeds 30KB, hccli writes the full output to a temp file and prints a warning to stderr:

```
⚠️  Output is large (47.3KB). Full output written to: /tmp/hccli-abc123.json

💡 To reduce output size:
  • Use fewer --breakdown flags
  • Use a shorter --time-range
  • Add filters to narrow results
```

The full JSON is still written to stdout, so piping to `jq` works normally.
