# Contributing to forge-ai

**Build, test, and contribute to the AI orchestration MCP server.**

## Quick Start

```bash
# Clone and enter
git clone https://github.com/alexandremahdhaoui/forge-ai.git
cd forge-ai

# Generate dependencies (tracker client + mocks)
forge build generate-tracker-client
forge build generate-mocks

# Build the binary
forge build forge-ai

# Run all tests (lint-tags, lint-licenses, lint, unit)
forge test-all
```

## How do I structure commits?

Each commit uses an emoji prefix and a structured body.

| Emoji | Meaning                          |
|-------|----------------------------------|
| `✨`  | New feature (feat:)              |
| `🐛`  | Bug fix (fix:)                   |
| `📖`  | Documentation (docs:)            |
| `🌱`  | Misc (chore:, test:, refactor:)  |

**Commit body format:**

```
✨ Short imperative summary (50 chars or less)

Why: Explain the motivation. What problem exists?

How: Describe the approach. What strategy did you choose?

What:

- internal/controller/plan_manager.go: add timeout handling
- internal/driver/mcp/tools.go: register new tool

How changes were verified:

- Unit tests for new logic
- forge test-all: all stages passed

Signed-off-by: Your Name <your@email.com>
```

Every commit requires `Signed-off-by`. Use `git commit -s` to add it automatically.

## How do I submit a pull request?

1. Create a feature branch from `main`.
2. Make changes. Run `forge test-all` before pushing.
3. Push the branch and open a PR.
4. PR title: short imperative summary (under 70 characters).
5. PR body: include "Why", "How", "What", and test verification.

## How do I run tests?

forge-ai has 4 test stages. Run them individually or all at once.

```bash
# Run all stages (build + lint-tags + lint-licenses + lint + unit)
forge test-all

# Run individual stages
forge test run lint-tags        # verify //go:build tags on test files
forge test run lint-licenses    # verify Apache 2.0 headers
forge test run lint             # golangci-lint static analysis
forge test run unit             # go test -tags unit ./...
```

To regenerate mocks after changing interfaces:

```bash
forge build generate-mocks
```

To regenerate the tracker client after API spec changes:

```bash
forge build generate-tracker-client
```

## How is the project structured?

```
forge-ai/
  cmd/
    forge-ai/
      main.go                 # Entry point: enginecli.Bootstrap
  internal/
    adapter/
      tracker.go              # TrackerClient interface (11 methods)
      tracker_http.go         # HTTPTrackerClient implementation
      tracker_http_test.go    # Adapter unit tests
    controller/
      plan_manager.go         # PlanManager (8 methods: plans, tasks)
      plan_manager_test.go    # PlanManager unit tests
      memory_manager.go       # MemoryManager (2 methods: comments)
      memory_manager_test.go  # MemoryManager unit tests
    driver/
      mcp/
        tools.go              # 10 MCP tool registrations
        tools_test.go         # Driver unit tests
        input.go              # Input structs for tool parameters
    types/
      types.go                # AgentContext, TaskAssignment, Memory
    util/
      mocks/                  # Generated mockery mocks
        mockadapter/          # TrackerClient mock
        mockcontroller/       # PlanManager + MemoryManager mocks
  pkg/
    generated/
      trackerclient/          # Auto-generated oapi-codegen client
  api/
    forge-tracker.v1.yaml     # OpenAPI spec for forge-tracker
  forge.yaml                  # Build and test configuration
  go.mod                      # Go 1.25.7
```

## What does each package do?

### Public packages

| Package                        | Description                              |
|--------------------------------|------------------------------------------|
| `pkg/generated/trackerclient`  | Auto-generated HTTP client for forge-tracker (oapi-codegen) |

### Internal packages

| Package                | Description                                          |
|------------------------|------------------------------------------------------|
| `internal/adapter`     | TrackerClient interface + HTTPTrackerClient impl     |
| `internal/controller`  | PlanManager (8 methods) + MemoryManager (2 methods)  |
| `internal/driver/mcp`  | 10 MCP tool registrations and input type definitions |
| `internal/types`       | Domain types: AgentContext, TaskAssignment, Memory    |
| `internal/util/mocks`  | Generated mocks (mockery) for all interfaces         |

## What conventions must I follow?

### Build tags

All test files require a build tag. Add this as the first line (after license header):

```go
//go:build unit
```

### License headers

Every `.go` file requires the Apache 2.0 license header:

```go
// Copyright 2024 Alexandre Mahdhaoui
//
// Licensed under the Apache License, Version 2.0 (the "License");
// ...
```

The `lint-licenses` stage enforces this. Run `bash ./hack/add-license-headers.sh` to add
missing headers.

### Generated files

Files prefixed with `zz_generated.` are auto-generated. Do not edit them manually.

- `pkg/generated/trackerclient/zz_generated.oapi-codegen.go` -- regenerate with
  `forge build generate-tracker-client`
- `internal/util/mocks/**/zz_generated.*.go` -- regenerate with `forge build generate-mocks`

### Dependencies

| Dependency              | Version  | Purpose                        |
|-------------------------|----------|--------------------------------|
| forge                   | v0.38.0  | enginecli bootstrap + mcpserver|
| go-sdk (MCP)            | v1.4.0   | MCP protocol types             |
| kin-openapi             | v0.133.0 | OpenAPI spec parsing           |
| oapi-codegen/runtime    | v1.2.0   | Generated client runtime       |

### Linting

Run `forge test run lint` before submitting. The linter config lives in the project root. Fix
all warnings -- the CI pipeline treats warnings as errors.
