# forge-ai Design

**forge-ai is an MCP server that gives AI agents task management, plan tracking, and shared memory through forge-tracker.**

## Problem Statement

AI agents that execute multi-step development plans need to coordinate. Each agent must know which
tasks exist, which are claimed, and what other agents discovered. Without a coordination layer,
agents duplicate work, lose inter-session context, and require manual task routing.

forge-ai solves this by exposing 28 MCP tools that map directly to forge-tracker REST endpoints.
Agents connect over stdio, call tools to manage plans and tasks, and share findings through a
tag-based memory system. The server adds no persistent state of its own -- forge-tracker is the
single source of truth.

## Tenets

Ordered by priority. When tenets conflict, higher-ranked wins.

1. **Stateless server.** forge-ai holds zero persistent state. All data lives in forge-tracker.
   This eliminates consistency bugs and simplifies deployment.
2. **Thin translation layer.** Each MCP tool maps to 1-2 forge-tracker API calls. No aggregation,
   caching, or business rules beyond what controllers require.
3. **Testable by construction.** Every layer depends on an interface. Mock any boundary in 1 line.
4. **Explicit over implicit.** Tool parameters are required, not inferred. Agents state exactly
   what tracking set and ticket they operate on.
5. **Composable tools.** Each tool does one thing. Agents compose workflows by calling tools in
   sequence. The server does not encode workflow assumptions.

## Requirements

From the perspective of an AI agent using forge-ai:

1. List and inspect meta-plans to understand project scope.
2. List and inspect plans to understand task breakdown.
3. List tasks with filters (status, assignee) to find available work.
4. Get task details including description and comments.
5. Assign a task to itself, setting status to `in_progress`.
6. Mark a task as `completed` when done.
7. Add tagged comments to tasks as shared memory.
8. List comments on a task, optionally filtered by author (agent ID).

## Out of Scope

- **Agent scheduling.** forge-ai does not decide which agent gets which task.
- **Conflict resolution.** No locking or optimistic concurrency on task assignment.
- **Persistent connections.** The server is stdio-only. No WebSocket or SSE transport.
- **Authentication.** forge-tracker handles access control. forge-ai passes requests through.

## Success Criteria

| Criteria                            | Target                                    |
|-------------------------------------|-------------------------------------------|
| MCP tools registered                | 28                                        |
| Test coverage (unit)                | All controller and driver functions tested |
| Lint stages                         | lint-tags, lint-licenses, lint all pass    |
| Build artifacts                     | 3 (generate-tracker-client, generate-mocks, binary) |
| External runtime dependencies       | 1 (forge-tracker)                         |
| Configuration parameters            | 1 (`FORGE_TRACKER_URL`)                   |

## Proposed Design

### Architecture

```
+---------------------------------------------------------------------+
|  AI Agent (Claude, GPT, etc.)                                       |
+----------------------------------+----------------------------------+
                                   | stdio (JSON-RPC 2.0 / MCP)
                                   v
+---------------------------------------------------------------------+
|  forge-ai MCP Server                                                |
|                                                                     |
|  +-------------------+    +------------------+    +---------------+ |
|  | MCP Driver        |    | Controllers        |    | Adapter       | |
|  | (28 tool regs)    |--->| PlanManager        |--->| TrackerClient | |
|  |                   |    | MemoryManager      |    | (interface)   | |
|  |                   |    | TrackingSetManager |    |               | |
|  |                   |    | EdgeManager        |    |               | |
|  +-------------------+    +------------------+    +------+--------+ |
|                                                          |          |
+----------------------------------------------------------+----------+
                                                           | HTTP
                                                           v
                                              +------------------------+
                                              |  forge-tracker         |
                                              |  REST API (:8080)      |
                                              +------------------------+
```

### Startup Sequence

```
main()
  |
  +--> enginecli.Bootstrap(config)
         |
         +--> runMCPServer()
                |
                +--> Read FORGE_TRACKER_URL (default: http://localhost:8080)
                +--> adapter.NewHTTPTrackerClient(url)
                +--> controller.NewPlanManager(client)
                +--> controller.NewMemoryManager(client)
                +--> controller.NewTrackingSetManager(client)
                +--> controller.NewEdgeManager(client)
                +--> mcpserver.New("forge-ai", version)
                +--> mcpdriver.RegisterTools(server, planMgr, memMgr, tsMgr, edgeMgr)
                +--> server.RunDefault()  // stdio loop
```

### Tool Call Flow (assign-task example)

```
Agent                  MCP Driver           PlanManager          Adapter             Tracker
  |                       |                     |                   |                   |
  |-- assign-task ------->|                     |                   |                   |
  |   {ts, ticketId,      |                     |                   |                   |
  |    agentId}           |                     |                   |                   |
  |                       |-- AssignTask ------>|                   |                   |
  |                       |                     |-- GetTicket ----->|                   |
  |                       |                     |                   |-- GET /ticket --->|
  |                       |                     |                   |<-- 200 OK --------|
  |                       |                     |<-- ticket --------|                   |
  |                       |                     |                   |                   |
  |                       |                     |-- UpdateTicket -->|                   |
  |                       |                     |   {status:        |-- PUT /ticket --->|
  |                       |                     |    in_progress,   |<-- 200 OK --------|
  |                       |                     |    assignee:      |                   |
  |                       |                     |    agentId}       |                   |
  |                       |                     |<-- ticket --------|                   |
  |                       |<-- ticket ----------|                   |                   |
  |<-- JSON result -------|                     |                   |                   |
```

### Memory Flow (add-comment + list-memories)

```
Agent A                MCP Driver           MemoryManager        Adapter
  |                       |                     |                   |
  |-- add-comment ------->|                     |                   |
  |   {ts, ticketId,      |                     |                   |
  |    author: "agent-a", |                     |                   |
  |    text: "found bug", |                     |                   |
  |    tags: ["blocker"]} |                     |                   |
  |                       |-- AddComment ------>|                   |
  |                       |                     |  format text:     |
  |                       |                     |  "[blocker] ..."  |
  |                       |                     |-- AddComment ---->|
  |                       |                     |<-- comment -------|
  |                       |<-- comment ---------|                   |
  |<-- JSON result -------|                     |                   |

Agent B                MCP Driver           MemoryManager        Adapter
  |                       |                     |                   |
  |-- list-memories ----->|                     |                   |
  |   {ts, ticketId,      |                     |                   |
  |    agentId: ""}       |                     |                   |
  |                       |-- ListMemories ---->|                   |
  |                       |                     |-- GetTicket ----->|
  |                       |                     |<-- ticket --------|
  |                       |                     |  parse tags from  |
  |                       |                     |  "[tag] text"     |
  |                       |                     |  format           |
  |                       |<-- []Memory --------|                   |
  |<-- JSON result -------|                     |                   |
```

## Technical Design

### Data Model

```go
// internal/types/types.go

type AgentContext struct {
    AgentID      string
    SessionID    string
    AssignedTask string
}

type TaskAssignment struct {
    TicketID string
    AgentID  string
}

type Memory struct {
    TicketID  string
    Timestamp time.Time
    Author    string
    Text      string
    Tags      []string    // parsed from "[tag1,tag2] text" format
}
```

The `Memory` type is the only domain type forge-ai defines. All other types (`MetaPlan`, `Plan`,
`Ticket`, `Comment`) come from the generated tracker client.

### Tag Format

Comments store tags inline: `[blocker,progress] actual comment text`. The `parseTags` function
extracts tags by finding the first `[...]` prefix. This avoids schema changes in forge-tracker.

### MCP Tool Catalog

| #  | Tool                | Controller         | Method            | Tracker API calls      |
|----|---------------------|--------------------|-------------------|------------------------|
| 1  | list-metaplans      | PlanManager        | ListMetaPlans     | 1 (ListMetaPlans)      |
| 2  | get-metaplan        | PlanManager        | GetMetaPlan       | 1 (GetMetaPlan)        |
| 3  | create-metaplan     | PlanManager        | CreateMetaPlan    | 1 (CreateMetaPlan)     |
| 4  | update-metaplan     | PlanManager        | UpdateMetaPlan    | 1 (UpdateMetaPlan)     |
| 5  | delete-metaplan     | PlanManager        | DeleteMetaPlan    | 1 (DeleteMetaPlan)     |
| 6  | list-plans          | PlanManager        | ListPlans         | 1 (ListPlans)          |
| 7  | get-plan-state      | PlanManager        | GetPlanState      | 1 (GetPlan)            |
| 8  | create-plan         | PlanManager        | CreatePlan        | 1 (CreatePlan)         |
| 9  | update-plan         | PlanManager        | UpdatePlan        | 1 (UpdatePlan)         |
| 10 | delete-plan         | PlanManager        | DeletePlan        | 1 (DeletePlan)         |
| 11 | list-tasks          | PlanManager        | ListTasks         | 1 (ListTickets)        |
| 12 | get-task            | PlanManager        | GetTask           | 1 (GetTicket)          |
| 13 | create-task         | PlanManager        | CreateTask        | 1 (CreateTicket)       |
| 14 | update-task         | PlanManager        | UpdateTask        | 1 (UpdateTicket)       |
| 15 | delete-task         | PlanManager        | DeleteTask        | 1 (DeleteTicket)       |
| 16 | assign-task         | PlanManager        | AssignTask        | 2 (Get + Update)       |
| 17 | complete-task       | PlanManager        | CompleteTask      | 2 (Get + Update)       |
| 18 | list-children       | PlanManager        | ListChildren      | 1 (GetChildren)        |
| 19 | list-blocking       | PlanManager        | ListBlocking      | 1 (GetBlocking)        |
| 20 | create-tracking-set | TrackingSetManager | CreateTrackingSet | 1 (CreateTrackingSet)  |
| 21 | list-tracking-sets  | TrackingSetManager | ListTrackingSets  | 1 (ListTrackingSets)   |
| 22 | get-tracking-set    | TrackingSetManager | GetTrackingSet    | 1 (GetTrackingSet)     |
| 23 | delete-tracking-set | TrackingSetManager | DeleteTrackingSet | 1 (DeleteTrackingSet)  |
| 24 | list-edges          | EdgeManager        | ListEdges         | 1 (ListEdges)          |
| 25 | create-edge         | EdgeManager        | AddEdge           | 1 (AddEdge)            |
| 26 | delete-edge         | EdgeManager        | RemoveEdge        | 1 (RemoveEdge)         |
| 27 | add-comment         | MemoryManager      | AddComment        | 1 (AddComment)         |
| 28 | list-memories       | MemoryManager      | ListMemories      | 1 (GetTicket)          |

### Package Catalog

#### Public packages

| Package                            | Description                              |
|------------------------------------|------------------------------------------|
| `pkg/generated/trackerclient`      | Auto-generated oapi-codegen HTTP client  |

#### Internal packages

| Package                            | Description                              |
|------------------------------------|------------------------------------------|
| `internal/adapter`                 | TrackerClient interface + HTTP impl      |
| `internal/controller`              | PlanManager (19) + MemoryManager (2) + TrackingSetManager (4) + EdgeManager (3) |
| `internal/driver/mcp`              | 28 MCP tool registrations + input types  |
| `internal/types`                   | AgentContext, TaskAssignment, Memory      |
| `internal/util/mocks/mockadapter`  | Generated mock for TrackerClient         |
| `internal/util/mocks/mockcontroller` | Generated mocks for PlanManager, MemoryManager, TrackingSetManager, EdgeManager |

### Adapter Interface

```go
// internal/adapter/tracker.go

type TrackerClient interface {
    // Tracking sets (4 methods)
    CreateTrackingSet(ctx, req)         (TrackingSet, error)
    ListTrackingSets(ctx, ...)          ([]TrackingSet, error)
    GetTrackingSet(ctx, name)           (TrackingSet, error)
    DeleteTrackingSet(ctx, name)        error

    // Meta-plans (5 methods)
    ListMetaPlans(ctx, ts)             ([]MetaPlan, error)
    GetMetaPlan(ctx, ts, id)           (MetaPlan, error)
    CreateMetaPlan(ctx, ts, req)       (MetaPlan, error)
    UpdateMetaPlan(ctx, ts, id, req)   (MetaPlan, error)
    DeleteMetaPlan(ctx, ts, id)        error

    // Plans (5 methods)
    ListPlans(ctx, ts)                 ([]Plan, error)
    GetPlan(ctx, ts, id)               (Plan, error)
    CreatePlan(ctx, ts, req)           (Plan, error)
    UpdatePlan(ctx, ts, id, req)       (Plan, error)
    DeletePlan(ctx, ts, id)            error

    // Tickets (5 methods)
    ListTickets(ctx, ts, filter)       ([]Ticket, error)
    GetTicket(ctx, ts, id)             (Ticket, error)
    CreateTicket(ctx, ts, req)         (Ticket, error)
    UpdateTicket(ctx, ts, id, req)     (Ticket, error)
    DeleteTicket(ctx, ts, id)          error

    // Graph queries (2 methods)
    GetChildren(ctx, ts, id)           ([]Ticket, error)
    GetBlocking(ctx, ts, id)           ([]Ticket, error)

    // Edges (3 methods)
    ListEdges(ctx, ts, filter)         ([]Edge, error)
    AddEdge(ctx, ts, req)              (Edge, error)
    RemoveEdge(ctx, ts, req)           error

    // Comments (1 method)
    AddComment(ctx, ts, id, req)       (Comment, error)
}
```

The interface has 25 methods. HTTPTrackerClient implements all 25 using the generated
oapi-codegen client with response validation.

## Design Patterns

**Interface-driven layers.** Each layer (adapter, controller, driver) depends on the layer below
through an interface. Tests inject mocks generated by mockery. No test touches the network.

**Stateless request handling.** Every tool call is independent. The server holds no request-scoped
or session-scoped state. This makes the server safe to restart at any time.

**Tag-in-text encoding.** Memory tags are encoded as a `[tag1,tag2]` prefix in comment text. This
avoids extending the forge-tracker schema while giving agents a structured search mechanism.

## Alternatives Considered

### Do nothing (agents call forge-tracker directly)

Agents could call forge-tracker REST API without an MCP intermediary. This was rejected because:
MCP is the standard protocol for AI tool use. Direct HTTP calls require each agent framework to
implement its own HTTP client, handle errors, and parse responses. forge-ai centralizes this in
28 typed tools with consistent error handling.

### Embed state in forge-ai (local database)

forge-ai could maintain a local SQLite database for agent assignments and memories. This was
rejected because it introduces a second source of truth, creates consistency bugs with
forge-tracker, and complicates deployment (now 2 stateful services instead of 1).

### WebSocket / SSE transport

A persistent connection transport would allow forge-tracker to push updates to agents. This was
rejected because MCP stdio transport is the standard for subprocess-based tool servers. Push
notifications add complexity with limited benefit -- agents poll when they need data.

## Risks and Mitigations

| Risk                                    | Mitigation                                |
|-----------------------------------------|-------------------------------------------|
| forge-tracker unavailable at startup    | Clear error message with URL in log       |
| Concurrent task assignment race         | Document at orchestrator level; no server-side lock |
| Generated client drift from API spec    | `forge build generate-tracker-client` in CI |
| Tag parsing breaks on malformed text    | `parseTags` returns nil for non-prefixed text |

## Testing Strategy

All tests use the `//go:build unit` tag. Tests run via `forge test run unit`.

| Layer       | Test approach                        | Mock target        |
|-------------|--------------------------------------|--------------------|
| Driver/MCP  | Call registered handler with mock controllers | PlanManager, MemoryManager |
| Controller  | Call controller methods with mock adapter     | TrackerClient      |
| Adapter     | Validate HTTP request construction (unit)     | httptest server    |

Mock generation uses mockery. Run `forge build generate-mocks` to regenerate.

4 test stages run in sequence via `forge test-all`:

1. `lint-tags` -- verifies all test files have build tags
2. `lint-licenses` -- verifies Apache 2.0 headers on all `.go` files
3. `lint` -- golangci-lint static analysis
4. `unit` -- `go test -tags unit ./...`

## FAQ

**Why does assign-task make 2 API calls (Get + Update)?**
The UpdateTicket endpoint requires the ticket title in the request body. AssignTask fetches the
current ticket first to preserve the existing title, then updates status and assignee.

**Why are tags stored inline in comment text?**
forge-tracker comments have `author` and `text` fields but no `tags` field. Encoding tags as a
`[tag1,tag2]` prefix avoids an API schema change while giving agents a structured way to
categorize memories.

**Why stdio instead of HTTP for the MCP transport?**
MCP clients (Claude Code, Cursor, etc.) spawn tool servers as subprocesses. Stdio is the standard
transport for this model. HTTP transport would require agents to manage a separate server process.

**Can forge-ai scale horizontally?**
Yes. Because forge-ai is stateless, you can run 1 instance per agent. Each instance connects
independently to forge-tracker. No shared state exists between instances.

**Why use oapi-codegen for the tracker client?**
The forge-tracker API publishes an OpenAPI v3 spec. oapi-codegen generates a type-safe Go client
with response validation. This eliminates hand-written HTTP code and keeps the client in sync
with the API spec.
