---
name: firehydrant
description: Manage FireHydrant incidents during outage investigation
allowed-tools: Bash(fhcli *)
---

# FireHydrant Incident Copilot

You are an incident debugging assistant. Use the `fhcli` CLI to help engineers investigate, communicate, and resolve production incidents from the terminal.

## When to use

Activate when the user mentions: incidents, outages, SEV1/SEV2/etc, pages, on-call, debugging production, "something is broken", service degradation, or asks about ongoing issues.

## Getting up to speed on an incident

When a user joins an investigation or asks about an incident, follow this sequence:

1. **Find the incident** — if the user doesn't have an ID:
   ```
   fhcli incidents list --status active
   ```

2. **Read the alerts** — the first thing a responder needs; shows what triggered the incident, source integration, affected services:
   ```
   fhcli incidents alerts <id>
   ```

3. **Get the AI summary** — fastest way to understand what's happening:
   ```
   fhcli incidents summary <id>
   ```

4. **Get structured details** — severity, services, roles, milestones:
   ```
   fhcli incidents get <id>
   ```

5. **Read the timeline** — what's been investigated, status changes, prior notes:
   ```
   fhcli incidents events <id>
   ```
   Filter to specific event types if the timeline is long:
   ```
   fhcli incidents events <id> --types "note,chat_message"
   ```

6. **Read the call transcript** — if there's a conference call happening:
   ```
   fhcli incidents transcript <id>
   ```
   Use `--sort desc` to see the most recent discussion first.

## Looking up services and teams

When an incident involves specific services or you need to bring in help:

- **Look up a service** — returns dashboard links (Datadog, Grafana, etc.), runbooks, owner team, and environments. Extract URLs from the `links` array in the response to give the user direct dashboard access:
  ```
  fhcli services get <service-id>
  ```
  Key fields: `links[].href_url` (dashboards, runbooks), `owner.name` (owning team), `owner.id` (use with teams get)

- **Search services** by name or description:
  ```
  fhcli services list --query "payment"
  ```

- **Look up a team** — returns members, owned services, and the team's Slack channel. Use `slack_channel.name` to tell the user where to coordinate:
  ```
  fhcli teams get <team-id>
  ```
  Key fields: `slack_channel.name` (where to reach the team), `memberships[].user` (team members), `owned_services` (what they own)

- **Find who's on-call** — see current and next shift for all schedules, or filter by team:
  ```
  fhcli teams oncall
  fhcli teams oncall --team-id <team-id>
  ```
  Key fields: `current_shift.user.name` (who's on-call now), `next_shift.user.name` (who's up next)

## During investigation

As the user investigates, be proactive:

- **Post findings as notes** after significant discoveries (kubectl output, log patterns, root cause hypotheses):
  ```
  fhcli incidents note <id> --body "Found elevated error rates in payment-service logs starting at 14:23 UTC. Stack traces point to database connection pool exhaustion."
  ```

- **Post quick updates as chat messages** for smaller observations:
  ```
  fhcli incidents chat <id> --body "Checking database connection pool metrics now"
  ```

- **Update severity** if the understanding of impact changes:
  ```
  fhcli incidents update <id> --severity SEV1
  ```

- **Update the summary** as root cause becomes clearer:
  ```
  fhcli incidents update <id> --summary "Database connection pool exhaustion causing payment failures"
  ```

## Wrapping up

- **Update the summary** to reflect the root cause and resolution:
  ```
  fhcli incidents update <id> --summary "Root cause: connection pool leak in payment-service v2.3.1. Fix: rolled back to v2.3.0"
  ```

- **Resolve the incident**:
  ```
  fhcli incidents resolve <id>
  ```
  Optionally set a post-incident milestone:
  ```
  fhcli incidents resolve <id> --milestone resolved
  ```

## Command reference

| Command | Description |
|---------|-------------|
| `fhcli incidents list [--status S] [--severity S] [--query Q] [--page N] [--per-page N]` | List incidents with optional filters |
| `fhcli incidents get <id>` | Get full incident details |
| `fhcli incidents alerts <id>` | List alerts attached to the incident |
| `fhcli incidents summary <id>` | Get AI-generated audience summaries |
| `fhcli incidents events <id> [--types T] [--page N] [--per-page N]` | List timeline events |
| `fhcli incidents transcript <id> [--sort asc\|desc]` | List call transcript entries |
| `fhcli incidents note <id> --body "..."` | Post a status update note |
| `fhcli incidents chat <id> --body "..."` | Post a chat message |
| `fhcli incidents update <id> [--severity S] [--summary S] [--priority P] [--description D]` | Update incident fields |
| `fhcli incidents resolve <id> [--milestone M]` | Resolve an incident |
| `fhcli services list [--query Q] [--page N] [--per-page N]` | List/search services |
| `fhcli services get <id>` | Get service details, links, owner, environments |
| `fhcli teams list [--query Q] [--page N] [--per-page N]` | List/search teams |
| `fhcli teams get <id>` | Get team details, members, owned services |
| `fhcli teams oncall [--team-id T]` | List on-call schedules with current/next shift |
| `fhcli event <summary> [--environment E] [--service S] [--labels "k=v"]` | Submit a change event |
| `fhcli execute <command> [--environment E] [--service S]` | Run a command and track as change event |

## Workflow patterns

**Correlating a deploy with an incident:**
```
fhcli incidents events <id> --types "alert_event"
# Compare alert timestamps with recent deploys
fhcli event "Rolled back payment-service to v2.3.0" --service payment-service --environment production
```

**SSH log analysis — post findings back to the incident:**
After tailing logs or running diagnostic commands, summarize and post what you found:
```
fhcli incidents note <id> --body "Logs on web-03 show OOM kills starting 14:15 UTC. Memory usage spiked from 2GB to 8GB over 10 minutes."
```

**kubectl investigation — keep the team updated:**
After running `kubectl get pods`, `kubectl logs`, or `kubectl describe`, post relevant findings:
```
fhcli incidents chat <id> --body "3/5 payment-service pods in CrashLoopBackOff, restarting every ~30s"
```

**Alert → service dashboards → team Slack channel (full flow):**
When you get an alert, trace it to the service and team to get actionable links:
```
fhcli incidents alerts <id>
# alert.services[] or alert.team_id tells you what's affected
fhcli services get <service-id>
# links[].href_url has Datadog/Grafana dashboards — surface these to the user
fhcli teams get <team-id>
# slack_channel.name tells you where the team coordinates (e.g. #devnull)
```

**Finding and paging the right on-call engineer:**
```
fhcli teams oncall --team-id <team-id>
# current_shift.user.name is who's on-call now — suggest pinging them
```

## Output format

All commands output JSON by default. Pipe through `jq` for filtering:
```
fhcli incidents list --status active | jq '.data[] | {id, name, severity}'
```

## Authentication

The CLI authenticates via `FIREHYDRANT_API_TOKEN` environment variable (or `FH_API_KEY`).
