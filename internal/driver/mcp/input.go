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

// Input structs for MCP tools. JSON tags must match the MCP tool parameter names.

type listMetaPlansInput struct {
	TS string `json:"ts" jsonschema:"Tracking set name"`
}

type getMetaPlanInput struct {
	TS string `json:"ts" jsonschema:"Tracking set name"`
	ID string `json:"id" jsonschema:"Meta-plan ID"`
}

type listPlansInput struct {
	TS string `json:"ts" jsonschema:"Tracking set name"`
}

type getPlanStateInput struct {
	TS string `json:"ts" jsonschema:"Tracking set name"`
	ID string `json:"id" jsonschema:"Plan ID"`
}

type listTasksInput struct {
	TS       string  `json:"ts" jsonschema:"Tracking set name"`
	Status   *string `json:"status,omitempty" jsonschema:"Filter by status"`
	Assignee *string `json:"assignee,omitempty" jsonschema:"Filter by assignee"`
}

type getTaskInput struct {
	TS string `json:"ts" jsonschema:"Tracking set name"`
	ID string `json:"id" jsonschema:"Task/ticket ID"`
}

type assignTaskInput struct {
	TS       string `json:"ts" jsonschema:"Tracking set name"`
	TicketID string `json:"ticketId" jsonschema:"Ticket ID to assign"`
	AgentID  string `json:"agentId" jsonschema:"Agent ID to assign to"`
}

type completeTaskInput struct {
	TS       string `json:"ts" jsonschema:"Tracking set name"`
	TicketID string `json:"ticketId" jsonschema:"Ticket ID to complete"`
}

type addCommentInput struct {
	TS       string   `json:"ts" jsonschema:"Tracking set name"`
	TicketID string   `json:"ticketId" jsonschema:"Ticket ID"`
	Text     string   `json:"text" jsonschema:"Comment text"`
	Author   string   `json:"author" jsonschema:"Comment author"`
	Tags     []string `json:"tags,omitempty" jsonschema:"Optional tags"`
}

type listMemoriesInput struct {
	TS       string `json:"ts" jsonschema:"Tracking set name"`
	TicketID string `json:"ticketId" jsonschema:"Ticket ID"`
	AgentID  string `json:"agentId,omitempty" jsonschema:"Filter by agent ID"`
}

type createPlanInput struct {
	TS          string             `json:"ts"          jsonschema:"Tracking set name"`
	ID          string             `json:"id"          jsonschema:"Plan ID"`
	Title       string             `json:"title"       jsonschema:"Plan title"`
	Status      *string            `json:"status,omitempty"      jsonschema:"Plan status"`
	Priority    *int               `json:"priority,omitempty"    jsonschema:"Priority (0=lowest)"`
	Labels      *[]string          `json:"labels,omitempty"      jsonschema:"Labels"`
	Annotations *map[string]string `json:"annotations,omitempty" jsonschema:"Key-value annotations"`
	Assignee    *string            `json:"assignee,omitempty"    jsonschema:"Assignee"`
	Description *string            `json:"description,omitempty" jsonschema:"Plan description"`
	Tasks       *[]string          `json:"tasks,omitempty"       jsonschema:"Task IDs in this plan"`
}

type createTaskInput struct {
	TS          string             `json:"ts"          jsonschema:"Tracking set name"`
	ID          string             `json:"id"          jsonschema:"Task ID"`
	Title       string             `json:"title"       jsonschema:"Task title"`
	Status      *string            `json:"status,omitempty"      jsonschema:"Task status"`
	Priority    *int               `json:"priority,omitempty"    jsonschema:"Priority (0=lowest)"`
	Labels      *[]string          `json:"labels,omitempty"      jsonschema:"Labels"`
	Annotations *map[string]string `json:"annotations,omitempty" jsonschema:"Key-value annotations"`
	Assignee    *string            `json:"assignee,omitempty"    jsonschema:"Assignee"`
	Description *string            `json:"description,omitempty" jsonschema:"Task description"`
}

type createMetaPlanInput struct {
	TS          string             `json:"ts"          jsonschema:"Tracking set name"`
	ID          string             `json:"id"          jsonschema:"Meta-plan ID"`
	Title       string             `json:"title"       jsonschema:"Meta-plan title"`
	Status      *string            `json:"status,omitempty"      jsonschema:"Meta-plan status"`
	Priority    *int               `json:"priority,omitempty"    jsonschema:"Priority (0=lowest)"`
	Labels      *[]string          `json:"labels,omitempty"      jsonschema:"Labels"`
	Annotations *map[string]string `json:"annotations,omitempty" jsonschema:"Key-value annotations"`
	Assignee    *string            `json:"assignee,omitempty"    jsonschema:"Assignee"`
	Description *string            `json:"description,omitempty" jsonschema:"Meta-plan description"`
}

type updatePlanInput struct {
	TS          string             `json:"ts"          jsonschema:"Tracking set name"`
	ID          string             `json:"id"          jsonschema:"Plan ID"`
	Title       string             `json:"title"       jsonschema:"Plan title (required)"`
	Status      *string            `json:"status,omitempty"      jsonschema:"Plan status"`
	Priority    *int               `json:"priority,omitempty"    jsonschema:"Priority (0=lowest)"`
	Labels      *[]string          `json:"labels,omitempty"      jsonschema:"Labels"`
	Annotations *map[string]string `json:"annotations,omitempty" jsonschema:"Key-value annotations"`
	Assignee    *string            `json:"assignee,omitempty"    jsonschema:"Assignee"`
	Description *string            `json:"description,omitempty" jsonschema:"Plan description"`
	Tasks       *[]string          `json:"tasks,omitempty"       jsonschema:"Task IDs in this plan"`
}

type updateMetaPlanInput struct {
	TS          string             `json:"ts"          jsonschema:"Tracking set name"`
	ID          string             `json:"id"          jsonschema:"Meta-plan ID"`
	Title       string             `json:"title"       jsonschema:"Meta-plan title (required)"`
	Status      *string            `json:"status,omitempty"      jsonschema:"Meta-plan status"`
	Priority    *int               `json:"priority,omitempty"    jsonschema:"Priority (0=lowest)"`
	Labels      *[]string          `json:"labels,omitempty"      jsonschema:"Labels"`
	Annotations *map[string]string `json:"annotations,omitempty" jsonschema:"Key-value annotations"`
	Assignee    *string            `json:"assignee,omitempty"    jsonschema:"Assignee"`
	Description *string            `json:"description,omitempty" jsonschema:"Meta-plan description"`
}

type updateTaskInput struct {
	TS          string             `json:"ts"          jsonschema:"Tracking set name"`
	ID          string             `json:"id"          jsonschema:"Task ID"`
	Title       string             `json:"title"       jsonschema:"Task title (required)"`
	Status      *string            `json:"status,omitempty"      jsonschema:"Task status"`
	Priority    *int               `json:"priority,omitempty"    jsonschema:"Priority (0=lowest)"`
	Labels      *[]string          `json:"labels,omitempty"      jsonschema:"Labels"`
	Annotations *map[string]string `json:"annotations,omitempty" jsonschema:"Key-value annotations"`
	Assignee    *string            `json:"assignee,omitempty"    jsonschema:"Assignee"`
	Description *string            `json:"description,omitempty" jsonschema:"Task description"`
}

type createEdgeInput struct {
	TS   string `json:"ts"   jsonschema:"Tracking set name"`
	From string `json:"from" jsonschema:"Source ticket ID"`
	To   string `json:"to"   jsonschema:"Target ticket ID"`
	Type string `json:"type" jsonschema:"Edge type: parent, blocks, or relates-to"`
}

type listEdgesInput struct {
	TS     string  `json:"ts"     jsonschema:"Tracking set name"`
	Ticket *string `json:"ticket,omitempty" jsonschema:"Filter by ticket ID (from or to)"`
	Type   *string `json:"type,omitempty"   jsonschema:"Filter by edge type"`
}

// P2 input structs

type deletePlanInput struct {
	TS string `json:"ts" jsonschema:"Tracking set name"`
	ID string `json:"id" jsonschema:"Plan ID"`
}

type deleteTaskInput struct {
	TS string `json:"ts" jsonschema:"Tracking set name"`
	ID string `json:"id" jsonschema:"Task ID"`
}

type deleteMetaPlanInput struct {
	TS string `json:"ts" jsonschema:"Tracking set name"`
	ID string `json:"id" jsonschema:"Meta-plan ID"`
}

type deleteEdgeInput struct {
	TS   string `json:"ts"   jsonschema:"Tracking set name"`
	From string `json:"from" jsonschema:"Source ticket ID"`
	To   string `json:"to"   jsonschema:"Target ticket ID"`
	Type string `json:"type" jsonschema:"Edge type: parent, blocks, or relates-to"`
}

type createTrackingSetInput struct {
	Name string `json:"name" jsonschema:"Tracking set name"`
}

type listTrackingSetsInput struct{}

type getTrackingSetInput struct {
	TS string `json:"ts" jsonschema:"Tracking set name"`
}

type deleteTrackingSetInput struct {
	TS string `json:"ts" jsonschema:"Tracking set name"`
}

// P3 input structs

type listChildrenInput struct {
	TS string `json:"ts" jsonschema:"Tracking set name"`
	ID string `json:"id" jsonschema:"Ticket ID"`
}

type listBlockingInput struct {
	TS string `json:"ts" jsonschema:"Tracking set name"`
	ID string `json:"id" jsonschema:"Ticket ID"`
}
