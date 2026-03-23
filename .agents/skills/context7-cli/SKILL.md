---
name: context7-cli
description: "Manage Context7 via CLI - search libraries, get documentation context. Use when user mentions 'context7', 'library docs', 'documentation context', 'fetch up-to-date library documentation', or needs current API references for coding."
category: documentation
---

# context7-cli

CLI for fetching up-to-date, version-specific documentation from Context7. The CLI is already installed on this system.

**Always use `--json` flag when calling commands programmatically.**

## Authentication

API key is already configured. If you need to change it:

```bash
context7-cli auth set "ctx7sk-your-api-key"
context7-cli auth test
```

Get your API key at https://context7.com/dashboard (keys start with `ctx7sk_`).

## Resources

### libs

Search for libraries in the Context7 database.

| Command | Description |
|---------|-------------|
| `context7-cli libs --name <lib> --query <q> --json` | Search libraries by name with query context |
| `context7-cli libs --name react --query "state management" --json` | Search React docs for state management |
| `context7-cli libs --name nextjs --query "app router" --fields id,title,stars` | Show specific fields |

**Flags:**
- `--name <name>` (required): Library name to search for (e.g. react, nextjs, vue)
- `--query <query>` (required): Your question or task for relevance ranking
- `--fields <cols>`: Comma-separated columns to display (e.g. id,title,description,stars)
- `--json`: Output as JSON
- `--format <fmt>`: Output format: text, json, csv, yaml

**Output format (JSON):**
```json
{
  "ok": true,
  "data": {
    "results": [...],
    "searchFilterApplied": false
  }
}
```

### context

Retrieve documentation context for a specific library.

| Command | Description |
|---------|-------------|
| `context7-cli context get --library <id> --query <q> --json` | Get documentation snippets as JSON |
| `context7-cli context get --library /facebook/react --query "useEffect" --json` | Get React useEffect docs |
| `context7-cli context get --library /vercel/next.js --query "app router" --type txt --raw` | Get docs as plain text |
| `context7-cli context get --library /vercel/next.js/v15.1.8 --query "middleware"` | Pin to specific version |

**Flags:**
- `--library <id>` (required): Library ID from search results (e.g. /facebook/react)
- `--query <query>` (required): Your question or task
- `--type <type>`: Response format: json or txt (default: "json")
- `--raw`: Output raw text (when --type txt)
- `--fields <cols>`: Comma-separated columns to display
- `--json`: Output as JSON
- `--format <fmt>`: Output format: text, json, csv, yaml

## Global Flags

All commands support:
- `--json`: Output as JSON (recommended for programmatic use)
- `--format <text|json|csv|yaml>`: Output format
- `--verbose`: Enable debug logging
- `--no-color`: Disable colored output
- `--no-header`: Omit table/csv headers (for piping)

## Typical Workflow

1. **Search for the library:**
   ```bash
   context7-cli libs --name react --query "state management" --json
   ```

2. **Extract the library ID from results** (e.g. `/facebook/react`)

3. **Get documentation context:**
   ```bash
   context7-cli context get --library /facebook/react --query "useState hook" --json
   ```

## Error Handling

Exit codes:
- `0`: Success
- `1`: API error (network, server, rate limit)
- `2`: Usage error (invalid arguments, missing token)

Common errors:
- **401**: Invalid or missing API token → Run `context7-cli auth set <token>`
- **404**: Resource not found → Verify library ID with `context7-cli libs --name <library>`
- **429**: Rate limited → Wait a moment and try again

## Examples

```bash
# Search for Vue.js composables documentation
context7-cli libs --name vue --query "composables" --json

# Get React Router documentation
context7-cli context get --library /remix-run/react-router --query "useNavigate" --json

# Get TypeScript docs as plain text
context7-cli context get --library /microsoft/typescript --query "generics" --type txt --raw

# Search with specific fields
context7-cli libs --name express --query "middleware" --fields id,title,stars,trustScore --json
```
