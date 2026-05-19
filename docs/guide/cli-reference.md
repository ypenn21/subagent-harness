# CLI Reference

## Usage

```
appteam [options]
```

Run without arguments to start the interactive wizard.

## Flags

| Flag | Description |
|------|-------------|
| `-h`, `--help` | Show help message and exit |
| `-v`, `--version` | Show version number and exit |
| `-r`, `--regenerate` | Regenerate all files from `.appteam/settings.json` without running the wizard |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `NO_COLOR` | Set to any value to disable ANSI color output (per [no-color.org](https://no-color.org/)) |

## Color Output

The wizard uses ANSI escape codes for colored output when running in a terminal. Colors are automatically disabled when:

- stdout is not a TTY (e.g., output is piped or redirected)
- The `NO_COLOR` environment variable is set

```bash
# Disable colors explicitly
NO_COLOR=1 appteam

# Pipe-safe (colors auto-disabled)
appteam | tee output.txt
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (unknown flag, missing settings.json for `-r`, generation failure) |

## Examples

```bash
# Interactive wizard
appteam

# Show version
appteam -v

# Regenerate from saved settings
appteam -r

# Non-interactive in CI/scripting
NO_COLOR=1 appteam -r
```
