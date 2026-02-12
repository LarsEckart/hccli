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

Provide your [Honeycomb API key](https://docs.honeycomb.io/get-started/configure/environments/manage-api-keys/) via the `--api-key` flag or the `HONEYCOMB_API_KEY` environment variable.

```bash
export HONEYCOMB_API_KEY=your-key-here
hccli auth
```

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
