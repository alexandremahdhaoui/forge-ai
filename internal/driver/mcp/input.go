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
