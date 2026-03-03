# FireHydrant CLI (`fh`)

You have access to the FireHydrant CLI for managing incidents and change events.

## Authentication

The CLI authenticates via `FIREHYDRANT_API_TOKEN` environment variable (or `FH_API_KEY`).

## Commands

### List incidents
```bash
fh incidents list [--page N] [--per-page N] [--query "search"] [--status "active"] [--severity "SEV1"]
```
Returns JSON with incident data and pagination metadata.

### Get incident details
```bash
fh incidents get <incident-id>
```
Returns full JSON details for a single incident.

### Post a status update note
```bash
fh incidents note <incident-id> --body "Status update text here"
```
Posts a status update note on the incident timeline.

### Post a chat message
```bash
fh incidents chat <incident-id> --body "Message text here"
```
Posts a chat message on the incident.

### Submit a change event
```bash
fh event <summary> [--environment production] [--service my-service] [--labels "key=value"]
```
Records a change event in FireHydrant.

### Execute and track a command
```bash
fh execute <command> [args...] [--environment production] [--service my-service]
```
Runs a shell command and submits timing/result as a change event.

## Output Format

All commands output JSON by default. Use `--interactive` or `-i` for human-readable output.

## Tips

- Pipe output through `jq` for filtering: `fh incidents list | jq '.data[].name'`
- Use `--query` to search incidents by name or description
- Use `--status` to filter by incident status (e.g., "active", "resolved")

allowed-tools: Bash(fh:*)
