# logslice

Fast structured log query tool that parses and filters JSON logs from files or stdin with a minimal query DSL.

---

## Installation

```bash
go install github.com/yourname/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logslice.git && cd logslice && go build -o logslice .
```

---

## Usage

```bash
# Filter logs from a file
logslice -f app.log 'level == "error"'

# Pipe from stdin
cat app.log | logslice 'status >= 500 AND service == "api"'

# Select specific fields
logslice -f app.log --fields time,level,msg 'level == "warn"'

# Pretty-print output
logslice -f app.log --pretty 'latency > 200'
```

### Query DSL

| Operator | Example |
|----------|---------|
| `==`     | `level == "error"` |
| `!=`     | `env != "dev"` |
| `>` `<` `>=` `<=` | `latency > 500` |
| `AND` / `OR` | `level == "error" AND service == "auth"` |

### Flags

| Flag | Description |
|------|-------------|
| `-f, --file` | Input log file (defaults to stdin) |
| `--fields` | Comma-separated list of fields to output |
| `--pretty` | Pretty-print JSON output |
| `--count` | Print only the number of matching lines |

---

## License

MIT © yourname