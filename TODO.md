# TODOs

- [ ] This repo should consume forge tracker API client from the forge-tracker repo.
      OR: We should use the forge have a forge-spec repo where all oapi specs are defined and all code depend on the specs from this repo and we generate the clients from there -> or clients could even be generated from there -> this way we ensure we are all in sync

## Feature Request: Expose Missing CRUD MCP Tools

### Problem

forge-ai exposes 10 MCP tools that map to forge-tracker REST endpoints. forge-tracker has 26
REST endpoints with full CRUD for metaplans, plans, tickets, edges, and tracking sets. forge-ai
only exposes read and update operations. Agents cannot create plans, tasks, or metaplans through
MCP — they can only read existing ones.

This blocks the forge-skills orchestrator and meta-orchestrator skills, which need to:
- Create plans during the planning phase
- Create tasks from plan breakdowns
- Create metaplans for cross-repo coordination
- Create edges to express task dependencies
- Update plans and metaplans as execution progresses

### Current State (10 tools)

| MCP Tool | HTTP Method | forge-tracker Endpoint | Operation |
|----------|-------------|------------------------|-----------|
| `list-metaplans` | GET | `/api/v1/tracking-sets/{ts}/metaplans` | Read |
| `get-metaplan` | GET | `/api/v1/tracking-sets/{ts}/metaplans/{id}` | Read |
| `list-plans` | GET | `/api/v1/tracking-sets/{ts}/plans` | Read |
| `get-plan-state` | GET | `/api/v1/tracking-sets/{ts}/plans/{id}` | Read |
| `list-tasks` | GET | `/api/v1/tracking-sets/{ts}/tickets` | Read |
| `get-task` | GET | `/api/v1/tracking-sets/{ts}/tickets/{id}` | Read |
| `assign-task` | PUT | `/api/v1/tracking-sets/{ts}/tickets/{id}` | Update |
| `complete-task` | PUT | `/api/v1/tracking-sets/{ts}/tickets/{id}` | Update |
| `add-comment` | POST | `/api/v1/tracking-sets/{ts}/tickets/{id}/comments` | Create |
| `list-memories` | GET | `/api/v1/tracking-sets/{ts}/tickets/{id}/comments` | Read |

### Missing Tools (16 tools needed)

These forge-tracker endpoints have no corresponding MCP tool:

**Priority 1 — Required for orchestrator/meta-orchestrator workflows:**

| Proposed MCP Tool | HTTP Method | forge-tracker Endpoint | Purpose |
|-------------------|-------------|------------------------|---------|
| `create-plan` | POST | `/api/v1/tracking-sets/{ts}/plans` | Create a plan with title and task list |
| `create-task` | POST | `/api/v1/tracking-sets/{ts}/tickets` | Create a task/ticket with title, description |
| `create-metaplan` | POST | `/api/v1/tracking-sets/{ts}/metaplans` | Create a meta-plan with stages and checkpoints |
| `update-plan` | PUT | `/api/v1/tracking-sets/{ts}/plans/{id}` | Update plan (add/remove tasks, change title) |
| `update-metaplan` | PUT | `/api/v1/tracking-sets/{ts}/metaplans/{id}` | Update meta-plan (change stages, checkpoints) |
| `update-task` | PUT | `/api/v1/tracking-sets/{ts}/tickets/{id}` | Update task (change title, description, labels) |
| `create-edge` | POST | `/api/v1/tracking-sets/{ts}/edges` | Create dependency edge between tasks |
| `list-edges` | GET | `/api/v1/tracking-sets/{ts}/edges` | List dependency edges (filter by ticket, type) |

**Priority 2 — Useful for cleanup and management:**

| Proposed MCP Tool | HTTP Method | forge-tracker Endpoint | Purpose |
|-------------------|-------------|------------------------|---------|
| `delete-plan` | DELETE | `/api/v1/tracking-sets/{ts}/plans/{id}` | Remove a plan |
| `delete-task` | DELETE | `/api/v1/tracking-sets/{ts}/tickets/{id}` | Remove a task/ticket |
| `delete-metaplan` | DELETE | `/api/v1/tracking-sets/{ts}/metaplans/{id}` | Remove a meta-plan |
| `delete-edge` | DELETE | `/api/v1/tracking-sets/{ts}/edges` | Remove dependency edge |
| `create-tracking-set` | POST | `/api/v1/tracking-sets` | Create a tracking set |
| `list-tracking-sets` | GET | `/api/v1/tracking-sets` | List all tracking sets |
| `get-tracking-set` | GET | `/api/v1/tracking-sets/{ts}` | Get tracking set details |
| `delete-tracking-set` | DELETE | `/api/v1/tracking-sets/{ts}` | Remove a tracking set |

**Priority 3 — Graph relationship queries:**

| Proposed MCP Tool | HTTP Method | forge-tracker Endpoint | Purpose |
|-------------------|-------------|------------------------|---------|
| `list-children` | GET | `/api/v1/tracking-sets/{ts}/tickets/{id}/children` | Get child tickets |
| `list-blocking` | GET | `/api/v1/tracking-sets/{ts}/tickets/{id}/blocking` | Get blocking tickets |

### Implementation Notes

- Each new MCP tool follows the existing pattern: thin MCP driver handler -> controller method -> adapter HTTP call
- forge-tracker already implements all REST endpoints. No forge-tracker changes needed.
- The existing `assign-task` and `complete-task` tools are thin wrappers around PUT `/tickets/{id}` with specific field updates. The new `update-task` tool generalizes this.
- `create-plan` accepts a `tasks` array in the request body (CreatePlanRequest in forge-tracker). This allows creating a plan with tasks in a single call.
- `create-metaplan` accepts `stages` and `checkpoints` arrays (CreateMetaPlanRequest).
- `create-task` accepts `id`, `title`, and optional fields like `description`, `labels`, `priority`.
- The `create-edge` tool needs `source`, `target`, and `type` parameters.
- Update the DESIGN.md "Out of Scope" section to remove "Plan creation" since it will become in-scope.
- Update the success criteria from "10 MCP tools" to the new total.
