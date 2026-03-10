# forge-ai

**AI orchestration MCP server that connects AI agents to forge-tracker for task management, memory, and plan coordination.**

> "I run 4 AI agents in parallel across 3 repositories. Each agent needs to claim tasks, report
> progress, and read what other agents discovered. forge-ai gives them a shared task board and
> memory system through 28 MCP tools, backed by forge-tracker."
> -- An AI platform engineer

## What problem does forge-ai solve?

AI agents that work on multi-step plans need a coordination layer. Without one, agents duplicate
work, lose context between sessions, and cannot share findings. forge-ai exposes 28 MCP tools
that let agents list plans, claim tasks, record progress, and read memories. It connects to
forge-tracker over HTTP, so agents coordinate through a single source of truth. The result: agents
stay aligned without custom glue code.

## Quick Start

```bash
# 1. Build
forge build forge-ai

# 2. Start forge-tracker (required dependency)
forge-tracker serve &

# 3. Configure and run
export FORGE_TRACKER_URL=http://localhost:8080
./build/bin/forge-ai

# 4. Add to your MCP client config
cat <<EOF
{
  "mcpServers": {
    "forge-ai": {
      "command": "./build/bin/forge-ai",
      "env": { "FORGE_TRACKER_URL": "http://localhost:8080" }
    }
  }
}
EOF
```

## How does it work?

```
+-------------+       +-------------------+       +-----------------+
|  AI Agent   | stdio |    forge-ai       |  HTTP |  forge-tracker  |
|  (Claude,   |------>|  MCP Server       |------>|  REST API       |
|   GPT, ...) |<------| 28 tools          |<------| :8080           |
+-------------+       +-------------------+       +-----------------+
                       |                   |
                       | mcp driver        |
                       |   -> controllers  |
                       |     -> adapter    |
                       +-------------------+
```

AI agents connect to forge-ai over stdio using the Model Context Protocol (MCP). forge-ai
translates each tool call into HTTP requests against the forge-tracker REST API. Four controllers
(PlanManager, MemoryManager, TrackingSetManager, EdgeManager) encapsulate business logic. The adapter layer abstracts the
generated HTTP client. See [DESIGN.md](DESIGN.md) for full architecture details.

## Table of Contents

- [How do I configure forge-ai?](#how-do-i-configure-forge-ai)
- [What MCP tools are available?](#what-mcp-tools-are-available)
- [How do I build and test?](#how-do-i-build-and-test)
- [FAQ](#faq)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## How do I configure forge-ai?

| Variable           | Default                  | Description                      |
|--------------------|--------------------------|----------------------------------|
| `FORGE_TRACKER_URL`| `http://localhost:8080`  | Base URL of forge-tracker server |

forge-ai reads configuration from environment variables. Set `FORGE_TRACKER_URL` to point at your
forge-tracker instance. No config file is required.

## What MCP tools are available?

forge-ai exposes 28 MCP tools in 7 categories.

### Meta-plan tools (5)

| Tool              | Description                                      | Required params  |
|-------------------|--------------------------------------------------|------------------|
| `list-metaplans`  | List all meta-plans in a tracking set            | `ts`             |
| `get-metaplan`    | Get meta-plan with stages, repos, checkpoints    | `ts`, `id`       |
| `create-metaplan` | Create a new meta-plan                           | `ts`, `id`, `title` |
| `update-metaplan` | Update an existing meta-plan                     | `ts`, `id`, `title` |
| `delete-metaplan` | Delete a meta-plan by ID                         | `ts`, `id`       |

### Plan tools (5)

| Tool              | Description                                      | Required params  |
|-------------------|--------------------------------------------------|------------------|
| `list-plans`      | List all plans in a tracking set                 | `ts`             |
| `get-plan-state`  | Get plan with resolved task statuses             | `ts`, `id`       |
| `create-plan`     | Create a new plan                                | `ts`, `id`, `title` |
| `update-plan`     | Update an existing plan                          | `ts`, `id`, `title` |
| `delete-plan`     | Delete a plan by ID                              | `ts`, `id`       |

### Task tools (7)

| Tool              | Description                                      | Required params             |
|-------------------|--------------------------------------------------|-----------------------------|
| `list-tasks`      | List tasks (optional: status, assignee filters)  | `ts`                        |
| `get-task`        | Get task details with description and comments   | `ts`, `id`                  |
| `create-task`     | Create a new task                                | `ts`, `id`, `title`         |
| `update-task`     | Update an existing task                          | `ts`, `id`, `title`         |
| `delete-task`     | Delete a task by ID                              | `ts`, `id`                  |
| `assign-task`     | Assign task to agent, set status to `in_progress`| `ts`, `ticketId`, `agentId` |
| `complete-task`   | Mark task as completed                           | `ts`, `ticketId`            |

### Tracking set tools (4)

| Tool                  | Description                      | Required params |
|-----------------------|----------------------------------|-----------------|
| `create-tracking-set` | Create a new tracking set        | `name`          |
| `list-tracking-sets`  | List all tracking sets           | (none)          |
| `get-tracking-set`    | Get a tracking set by name       | `ts`            |
| `delete-tracking-set` | Delete a tracking set by name    | `ts`            |

### Edge tools (3)

| Tool              | Description                                      | Required params              |
|-------------------|--------------------------------------------------|------------------------------|
| `list-edges`      | List edges, optionally filtered by ticket or type| `ts`                         |
| `create-edge`     | Create a relationship edge between two tickets   | `ts`, `from`, `to`, `type`   |
| `delete-edge`     | Delete a relationship edge                       | `ts`, `from`, `to`, `type`   |

### Graph query tools (2)

| Tool              | Description                                      | Required params  |
|-------------------|--------------------------------------------------|------------------|
| `list-children`   | List child tickets of a given ticket             | `ts`, `id`       |
| `list-blocking`   | List tickets that block a given ticket           | `ts`, `id`       |

### Memory tools (2)

| Tool              | Description                                      | Required params                    |
|-------------------|--------------------------------------------------|------------------------------------|
| `add-comment`     | Add tagged comment to a task                     | `ts`, `ticketId`, `text`, `author` |
| `list-memories`   | List comments, optionally filtered by agent      | `ts`, `ticketId`                   |

All tools require the `ts` (tracking set) parameter except `list-tracking-sets` and
`create-tracking-set`. Tags in comments use the format `[tag1,tag2] comment text`.

## How do I build and test?

```bash
# Build all artifacts (generates client, mocks, binary)
forge build

# Build only the binary
forge build forge-ai

# Run all test stages (lint-tags, lint-licenses, lint, unit)
forge test-all

# Run a specific test stage
forge test run unit
forge test run lint
```

Build targets: `generate-tracker-client`, `generate-mocks`, `forge-ai`.
Test stages: `lint-tags`, `lint-licenses`, `lint`, `unit`.

## FAQ

**Does forge-ai require forge-tracker to be running?**
Yes. forge-ai makes HTTP calls to forge-tracker on every tool invocation. Start forge-tracker
before running forge-ai.

**What transport does the MCP server use?**
Stdio. The AI agent spawns forge-ai as a subprocess and communicates over stdin/stdout using
JSON-RPC 2.0.

**Can I run forge-ai without forge (the build tool)?**
Yes. Run `go build -o build/bin/forge-ai ./cmd/forge-ai` directly. forge simplifies build and
test orchestration but is not required at runtime.

**What Go version is required?**
Go 1.25.7 or later.

**How do agents share context?**
Agents write comments (memories) to tasks via `add-comment` with tags. Other agents read those
memories via `list-memories`, optionally filtering by agent ID.

**What happens if two agents assign the same task?**
The last `assign-task` call wins. forge-tracker does not enforce locking. Coordinate assignment
at the orchestrator level.

**How do I regenerate the tracker client?**
Run `forge build generate-tracker-client`. The OpenAPI spec lives at `api/forge-tracker.v1.yaml`.

## Documentation

| Document                              | Audience    | Description              |
|---------------------------------------|-------------|--------------------------|
| [README.md](README.md)               | Users       | Quick start and usage    |
| [DESIGN.md](DESIGN.md)               | Developers  | Architecture and design  |
| [CONTRIBUTING.md](CONTRIBUTING.md)    | Contributors| Build, test, contribute  |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for build instructions, commit conventions, and project
structure.

## License

Apache License 2.0. See [LICENSE](LICENSE).
