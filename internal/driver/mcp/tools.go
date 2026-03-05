// Copyright 2024 Alexandre Mahdhaoui
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/alexandremahdhaoui/forge-ai/internal/adapter"
	"github.com/alexandremahdhaoui/forge-ai/internal/controller"
	"github.com/alexandremahdhaoui/forge/pkg/mcpserver"
)

// RegisterTools registers all forge-ai MCP tools on the given server.
func RegisterTools(
	server *mcpserver.Server,
	planMgr controller.PlanManager,
	memMgr controller.MemoryManager,
) {
	registerListMetaPlans(server, planMgr)
	registerGetMetaPlan(server, planMgr)
	registerListPlans(server, planMgr)
	registerGetPlanState(server, planMgr)
	registerListTasks(server, planMgr)
	registerGetTask(server, planMgr)
	registerAssignTask(server, planMgr)
	registerCompleteTask(server, planMgr)
	registerAddComment(server, memMgr)
	registerListMemories(server, memMgr)
}

func registerListMetaPlans(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "list-metaplans",
		Description: "List all meta-plans in a tracking set.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input listMetaPlansInput) (*gomcp.CallToolResult, any, error) {
		result, err := planMgr.ListMetaPlans(ctx, input.TS)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(result)
	})
}

func registerGetMetaPlan(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "get-metaplan",
		Description: "Get a meta-plan by ID with stages, repos, and checkpoints.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input getMetaPlanInput) (*gomcp.CallToolResult, any, error) {
		result, err := planMgr.GetMetaPlan(ctx, input.TS, input.ID)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(result)
	})
}

func registerListPlans(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "list-plans",
		Description: "List all plans in a tracking set.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input listPlansInput) (*gomcp.CallToolResult, any, error) {
		result, err := planMgr.ListPlans(ctx, input.TS)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(result)
	})
}

func registerGetPlanState(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "get-plan-state",
		Description: "Get a plan by ID with resolved task statuses.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input getPlanStateInput) (*gomcp.CallToolResult, any, error) {
		result, err := planMgr.GetPlanState(ctx, input.TS, input.ID)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(result)
	})
}

func registerListTasks(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "list-tasks",
		Description: "List tasks with optional status and assignee filters.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input listTasksInput) (*gomcp.CallToolResult, any, error) {
		filter := adapter.TicketFilter{}
		if input.Status != nil {
			filter.Status = *input.Status
		}
		if input.Assignee != nil {
			filter.Assignee = *input.Assignee
		}
		result, err := planMgr.ListTasks(ctx, input.TS, filter)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(result)
	})
}

func registerGetTask(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "get-task",
		Description: "Get a task by ID with description and comments.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input getTaskInput) (*gomcp.CallToolResult, any, error) {
		result, err := planMgr.GetTask(ctx, input.TS, input.ID)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(result)
	})
}

func registerAssignTask(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "assign-task",
		Description: "Assign a task to an agent. Sets assignee and status to in_progress.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input assignTaskInput) (*gomcp.CallToolResult, any, error) {
		result, err := planMgr.AssignTask(ctx, input.TS, input.TicketID, input.AgentID)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(result)
	})
}

func registerCompleteTask(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "complete-task",
		Description: "Mark a task as completed.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input completeTaskInput) (*gomcp.CallToolResult, any, error) {
		result, err := planMgr.CompleteTask(ctx, input.TS, input.TicketID)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(result)
	})
}

func registerAddComment(server *mcpserver.Server, memMgr controller.MemoryManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "add-comment",
		Description: "Add a comment/memory to a task.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input addCommentInput) (*gomcp.CallToolResult, any, error) {
		result, err := memMgr.AddComment(ctx, input.TS, input.TicketID, input.Author, input.Text, input.Tags)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(result)
	})
}

func registerListMemories(server *mcpserver.Server, memMgr controller.MemoryManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "list-memories",
		Description: "List memories/comments for a task, optionally filtered by agent.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input listMemoriesInput) (*gomcp.CallToolResult, any, error) {
		result, err := memMgr.ListMemories(ctx, input.TS, input.TicketID, input.AgentID)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(result)
	})
}

func jsonResult(v any) (*gomcp.CallToolResult, any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, nil, fmt.Errorf("marshaling result: %w", err)
	}
	return &gomcp.CallToolResult{
		Content: []gomcp.Content{
			&gomcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func errResult(err error) (*gomcp.CallToolResult, any, error) {
	return nil, nil, err
}
