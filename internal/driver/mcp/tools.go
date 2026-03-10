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
	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
	"github.com/alexandremahdhaoui/forge/pkg/mcpserver"
)

// RegisterTools registers all forge-ai MCP tools on the given server.
func RegisterTools(
	server *mcpserver.Server,
	planMgr controller.PlanManager,
	memMgr controller.MemoryManager,
	tsMgr controller.TrackingSetManager,
	edgeMgr controller.EdgeManager,
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
	registerCreatePlan(server, planMgr)
	registerCreateTask(server, planMgr)
	registerCreateMetaPlan(server, planMgr)
	registerUpdatePlan(server, planMgr)
	registerUpdateMetaPlan(server, planMgr)
	registerUpdateTask(server, planMgr)
	registerCreateEdge(server, edgeMgr)
	registerListEdges(server, edgeMgr)
	// Delete tools
	registerDeletePlan(server, planMgr)
	registerDeleteTask(server, planMgr)
	registerDeleteMetaPlan(server, planMgr)
	registerDeleteEdge(server, edgeMgr)
	registerCreateTrackingSet(server, tsMgr)
	registerListTrackingSets(server, tsMgr)
	registerGetTrackingSet(server, tsMgr)
	registerDeleteTrackingSet(server, tsMgr)
	// Graph query tools
	registerListChildren(server, planMgr)
	registerListBlocking(server, planMgr)
}

func handleListMetaPlans(ctx context.Context, input listMetaPlansInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.ListMetaPlans(ctx, input.TS)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerListMetaPlans(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "list-metaplans",
		Description: "List all meta-plans in a tracking set.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input listMetaPlansInput) (*gomcp.CallToolResult, any, error) {
		return handleListMetaPlans(ctx, input, planMgr)
	})
}

func handleGetMetaPlan(ctx context.Context, input getMetaPlanInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.GetMetaPlan(ctx, input.TS, input.ID)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerGetMetaPlan(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "get-metaplan",
		Description: "Get a meta-plan by ID with stages, repos, and checkpoints.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input getMetaPlanInput) (*gomcp.CallToolResult, any, error) {
		return handleGetMetaPlan(ctx, input, planMgr)
	})
}

func handleListPlans(ctx context.Context, input listPlansInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.ListPlans(ctx, input.TS)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerListPlans(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "list-plans",
		Description: "List all plans in a tracking set.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input listPlansInput) (*gomcp.CallToolResult, any, error) {
		return handleListPlans(ctx, input, planMgr)
	})
}

func handleGetPlanState(ctx context.Context, input getPlanStateInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.GetPlanState(ctx, input.TS, input.ID)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerGetPlanState(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "get-plan-state",
		Description: "Get a plan by ID with resolved task statuses.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input getPlanStateInput) (*gomcp.CallToolResult, any, error) {
		return handleGetPlanState(ctx, input, planMgr)
	})
}

func handleListTasks(ctx context.Context, input listTasksInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
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
}

func registerListTasks(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "list-tasks",
		Description: "List tasks with optional status and assignee filters.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input listTasksInput) (*gomcp.CallToolResult, any, error) {
		return handleListTasks(ctx, input, planMgr)
	})
}

func handleGetTask(ctx context.Context, input getTaskInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.GetTask(ctx, input.TS, input.ID)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerGetTask(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "get-task",
		Description: "Get a task by ID with description and comments.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input getTaskInput) (*gomcp.CallToolResult, any, error) {
		return handleGetTask(ctx, input, planMgr)
	})
}

func handleAssignTask(ctx context.Context, input assignTaskInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.AssignTask(ctx, input.TS, input.TicketID, input.AgentID)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerAssignTask(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "assign-task",
		Description: "Assign a task to an agent. Sets assignee and status to in_progress.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input assignTaskInput) (*gomcp.CallToolResult, any, error) {
		return handleAssignTask(ctx, input, planMgr)
	})
}

func handleCompleteTask(ctx context.Context, input completeTaskInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.CompleteTask(ctx, input.TS, input.TicketID)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerCompleteTask(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "complete-task",
		Description: "Mark a task as completed.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input completeTaskInput) (*gomcp.CallToolResult, any, error) {
		return handleCompleteTask(ctx, input, planMgr)
	})
}

func handleAddComment(ctx context.Context, input addCommentInput, memMgr controller.MemoryManager) (*gomcp.CallToolResult, any, error) {
	result, err := memMgr.AddComment(ctx, input.TS, input.TicketID, input.Author, input.Text, input.Tags)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerAddComment(server *mcpserver.Server, memMgr controller.MemoryManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "add-comment",
		Description: "Add a comment/memory to a task.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input addCommentInput) (*gomcp.CallToolResult, any, error) {
		return handleAddComment(ctx, input, memMgr)
	})
}

func handleListMemories(ctx context.Context, input listMemoriesInput, memMgr controller.MemoryManager) (*gomcp.CallToolResult, any, error) {
	result, err := memMgr.ListMemories(ctx, input.TS, input.TicketID, input.AgentID)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerListMemories(server *mcpserver.Server, memMgr controller.MemoryManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "list-memories",
		Description: "List memories/comments for a task, optionally filtered by agent.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input listMemoriesInput) (*gomcp.CallToolResult, any, error) {
		return handleListMemories(ctx, input, memMgr)
	})
}

func handleCreatePlan(ctx context.Context, input createPlanInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.CreatePlan(ctx, input.TS, tc.CreatePlanRequest{
		Id:          input.ID,
		Title:       input.Title,
		Status:      input.Status,
		Priority:    input.Priority,
		Labels:      input.Labels,
		Annotations: input.Annotations,
		Assignee:    input.Assignee,
		Description: input.Description,
		Tasks:       input.Tasks,
	})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerCreatePlan(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "create-plan",
		Description: "Create a new plan in a tracking set.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input createPlanInput) (*gomcp.CallToolResult, any, error) {
		return handleCreatePlan(ctx, input, planMgr)
	})
}

func handleCreateTask(ctx context.Context, input createTaskInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.CreateTask(ctx, input.TS, tc.CreateTicketRequest{
		Id:          input.ID,
		Title:       input.Title,
		Status:      input.Status,
		Priority:    input.Priority,
		Labels:      input.Labels,
		Annotations: input.Annotations,
		Assignee:    input.Assignee,
		Description: input.Description,
	})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerCreateTask(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "create-task",
		Description: "Create a new task in a tracking set.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input createTaskInput) (*gomcp.CallToolResult, any, error) {
		return handleCreateTask(ctx, input, planMgr)
	})
}

func handleCreateMetaPlan(ctx context.Context, input createMetaPlanInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.CreateMetaPlan(ctx, input.TS, tc.CreateMetaPlanRequest{
		Id:          input.ID,
		Title:       input.Title,
		Status:      input.Status,
		Priority:    input.Priority,
		Labels:      input.Labels,
		Annotations: input.Annotations,
		Assignee:    input.Assignee,
		Description: input.Description,
	})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerCreateMetaPlan(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "create-metaplan",
		Description: "Create a new meta-plan in a tracking set.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input createMetaPlanInput) (*gomcp.CallToolResult, any, error) {
		return handleCreateMetaPlan(ctx, input, planMgr)
	})
}

func handleUpdatePlan(ctx context.Context, input updatePlanInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.UpdatePlan(ctx, input.TS, input.ID, tc.UpdatePlanRequest{
		Title:       input.Title,
		Status:      input.Status,
		Priority:    input.Priority,
		Labels:      input.Labels,
		Annotations: input.Annotations,
		Assignee:    input.Assignee,
		Description: input.Description,
		Tasks:       input.Tasks,
	})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerUpdatePlan(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "update-plan",
		Description: "Update an existing plan.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input updatePlanInput) (*gomcp.CallToolResult, any, error) {
		return handleUpdatePlan(ctx, input, planMgr)
	})
}

func handleUpdateMetaPlan(ctx context.Context, input updateMetaPlanInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.UpdateMetaPlan(ctx, input.TS, input.ID, tc.UpdateMetaPlanRequest{
		Title:       input.Title,
		Status:      input.Status,
		Priority:    input.Priority,
		Labels:      input.Labels,
		Annotations: input.Annotations,
		Assignee:    input.Assignee,
		Description: input.Description,
	})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerUpdateMetaPlan(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "update-metaplan",
		Description: "Update an existing meta-plan.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input updateMetaPlanInput) (*gomcp.CallToolResult, any, error) {
		return handleUpdateMetaPlan(ctx, input, planMgr)
	})
}

func handleUpdateTask(ctx context.Context, input updateTaskInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.UpdateTask(ctx, input.TS, input.ID, tc.UpdateTicketRequest{
		Title:       input.Title,
		Status:      input.Status,
		Priority:    input.Priority,
		Labels:      input.Labels,
		Annotations: input.Annotations,
		Assignee:    input.Assignee,
		Description: input.Description,
	})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerUpdateTask(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "update-task",
		Description: "Update an existing task.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input updateTaskInput) (*gomcp.CallToolResult, any, error) {
		return handleUpdateTask(ctx, input, planMgr)
	})
}

func handleCreateEdge(ctx context.Context, input createEdgeInput, edgeMgr controller.EdgeManager) (*gomcp.CallToolResult, any, error) {
	result, err := edgeMgr.AddEdge(ctx, input.TS, tc.EdgeRequest{
		From: input.From,
		To:   input.To,
		Type: input.Type,
	})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerCreateEdge(server *mcpserver.Server, edgeMgr controller.EdgeManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "create-edge",
		Description: "Create a relationship edge between two tickets.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input createEdgeInput) (*gomcp.CallToolResult, any, error) {
		return handleCreateEdge(ctx, input, edgeMgr)
	})
}

func handleListEdges(ctx context.Context, input listEdgesInput, edgeMgr controller.EdgeManager) (*gomcp.CallToolResult, any, error) {
	params := &tc.ListEdgesParams{
		Ticket: input.Ticket,
		Type:   input.Type,
	}
	result, err := edgeMgr.ListEdges(ctx, input.TS, params)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerListEdges(server *mcpserver.Server, edgeMgr controller.EdgeManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "list-edges",
		Description: "List edges in a tracking set, optionally filtered by ticket or type.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input listEdgesInput) (*gomcp.CallToolResult, any, error) {
		return handleListEdges(ctx, input, edgeMgr)
	})
}

// Delete tool registrations

func handleDeletePlan(ctx context.Context, input deletePlanInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	if err := planMgr.DeletePlan(ctx, input.TS, input.ID); err != nil {
		return errResult(err)
	}
	return textResult(fmt.Sprintf("deleted plan %s", input.ID))
}

func registerDeletePlan(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "delete-plan",
		Description: "Delete a plan by ID.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input deletePlanInput) (*gomcp.CallToolResult, any, error) {
		return handleDeletePlan(ctx, input, planMgr)
	})
}

func handleDeleteTask(ctx context.Context, input deleteTaskInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	if err := planMgr.DeleteTask(ctx, input.TS, input.ID); err != nil {
		return errResult(err)
	}
	return textResult(fmt.Sprintf("deleted task %s", input.ID))
}

func registerDeleteTask(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "delete-task",
		Description: "Delete a task by ID.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input deleteTaskInput) (*gomcp.CallToolResult, any, error) {
		return handleDeleteTask(ctx, input, planMgr)
	})
}

func handleDeleteMetaPlan(ctx context.Context, input deleteMetaPlanInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	if err := planMgr.DeleteMetaPlan(ctx, input.TS, input.ID); err != nil {
		return errResult(err)
	}
	return textResult(fmt.Sprintf("deleted meta-plan %s", input.ID))
}

func registerDeleteMetaPlan(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "delete-metaplan",
		Description: "Delete a meta-plan by ID.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input deleteMetaPlanInput) (*gomcp.CallToolResult, any, error) {
		return handleDeleteMetaPlan(ctx, input, planMgr)
	})
}

func handleDeleteEdge(ctx context.Context, input deleteEdgeInput, edgeMgr controller.EdgeManager) (*gomcp.CallToolResult, any, error) {
	if err := edgeMgr.RemoveEdge(ctx, input.TS, tc.EdgeRequest{
		From: input.From,
		To:   input.To,
		Type: input.Type,
	}); err != nil {
		return errResult(err)
	}
	return textResult(fmt.Sprintf("deleted edge %s -> %s", input.From, input.To))
}

func registerDeleteEdge(server *mcpserver.Server, edgeMgr controller.EdgeManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "delete-edge",
		Description: "Delete a relationship edge between two tickets.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input deleteEdgeInput) (*gomcp.CallToolResult, any, error) {
		return handleDeleteEdge(ctx, input, edgeMgr)
	})
}

func handleCreateTrackingSet(ctx context.Context, input createTrackingSetInput, tsMgr controller.TrackingSetManager) (*gomcp.CallToolResult, any, error) {
	result, err := tsMgr.CreateTrackingSet(ctx, input.Name)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerCreateTrackingSet(server *mcpserver.Server, tsMgr controller.TrackingSetManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "create-tracking-set",
		Description: "Create a new tracking set.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input createTrackingSetInput) (*gomcp.CallToolResult, any, error) {
		return handleCreateTrackingSet(ctx, input, tsMgr)
	})
}

func handleListTrackingSets(ctx context.Context, _ listTrackingSetsInput, tsMgr controller.TrackingSetManager) (*gomcp.CallToolResult, any, error) {
	result, err := tsMgr.ListTrackingSets(ctx)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerListTrackingSets(server *mcpserver.Server, tsMgr controller.TrackingSetManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "list-tracking-sets",
		Description: "List all tracking sets.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input listTrackingSetsInput) (*gomcp.CallToolResult, any, error) {
		return handleListTrackingSets(ctx, input, tsMgr)
	})
}

func handleGetTrackingSet(ctx context.Context, input getTrackingSetInput, tsMgr controller.TrackingSetManager) (*gomcp.CallToolResult, any, error) {
	result, err := tsMgr.GetTrackingSet(ctx, input.TS)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerGetTrackingSet(server *mcpserver.Server, tsMgr controller.TrackingSetManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "get-tracking-set",
		Description: "Get a tracking set by name.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input getTrackingSetInput) (*gomcp.CallToolResult, any, error) {
		return handleGetTrackingSet(ctx, input, tsMgr)
	})
}

func handleDeleteTrackingSet(ctx context.Context, input deleteTrackingSetInput, tsMgr controller.TrackingSetManager) (*gomcp.CallToolResult, any, error) {
	if err := tsMgr.DeleteTrackingSet(ctx, input.TS); err != nil {
		return errResult(err)
	}
	return textResult(fmt.Sprintf("deleted tracking set %s", input.TS))
}

func registerDeleteTrackingSet(server *mcpserver.Server, tsMgr controller.TrackingSetManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "delete-tracking-set",
		Description: "Delete a tracking set by name.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input deleteTrackingSetInput) (*gomcp.CallToolResult, any, error) {
		return handleDeleteTrackingSet(ctx, input, tsMgr)
	})
}

// Graph query tool registrations

func handleListChildren(ctx context.Context, input listChildrenInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.ListChildren(ctx, input.TS, input.ID)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerListChildren(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "list-children",
		Description: "List child tickets of a given ticket.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input listChildrenInput) (*gomcp.CallToolResult, any, error) {
		return handleListChildren(ctx, input, planMgr)
	})
}

func handleListBlocking(ctx context.Context, input listBlockingInput, planMgr controller.PlanManager) (*gomcp.CallToolResult, any, error) {
	result, err := planMgr.ListBlocking(ctx, input.TS, input.ID)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(result)
}

func registerListBlocking(server *mcpserver.Server, planMgr controller.PlanManager) {
	mcpserver.RegisterTool(server, &gomcp.Tool{
		Name:        "list-blocking",
		Description: "List tickets that block a given ticket.",
	}, func(ctx context.Context, req *gomcp.CallToolRequest, input listBlockingInput) (*gomcp.CallToolResult, any, error) {
		return handleListBlocking(ctx, input, planMgr)
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

func textResult(msg string) (*gomcp.CallToolResult, any, error) {
	return &gomcp.CallToolResult{
		Content: []gomcp.Content{
			&gomcp.TextContent{Text: msg},
		},
	}, nil, nil
}
