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

package controller

import (
	"context"

	"github.com/alexandremahdhaoui/forge-ai/internal/adapter"
	"github.com/alexandremahdhaoui/forge-ai/internal/types"
	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
)

// PlanManager handles plan, meta-plan, and task operations.
// TODO: Consider decomposing into smaller interfaces (e.g., TaskManager, MetaPlanManager)
// when the method count exceeds maintainability thresholds. Currently 19 methods spanning
// plans, meta-plans, tasks, and graph queries.
type PlanManager interface {
	ListMetaPlans(ctx context.Context, ts string) ([]tc.MetaPlan, error)
	GetMetaPlan(ctx context.Context, ts, id string) (tc.MetaPlan, error)
	ListPlans(ctx context.Context, ts string) ([]tc.Plan, error)
	GetPlanState(ctx context.Context, ts, id string) (tc.Plan, error)
	ListTasks(ctx context.Context, ts string, filter adapter.TicketFilter) ([]tc.Ticket, error)
	GetTask(ctx context.Context, ts, id string) (tc.Ticket, error)
	AssignTask(ctx context.Context, ts, ticketID, agentID string) (tc.Ticket, error)
	CompleteTask(ctx context.Context, ts, ticketID string) (tc.Ticket, error)

	// MetaPlan CRUD
	CreateMetaPlan(ctx context.Context, ts string, req tc.CreateMetaPlanRequest) (tc.MetaPlan, error)
	UpdateMetaPlan(ctx context.Context, ts, id string, req tc.UpdateMetaPlanRequest) (tc.MetaPlan, error)
	DeleteMetaPlan(ctx context.Context, ts, id string) error

	// Plan CRUD
	CreatePlan(ctx context.Context, ts string, req tc.CreatePlanRequest) (tc.Plan, error)
	UpdatePlan(ctx context.Context, ts, id string, req tc.UpdatePlanRequest) (tc.Plan, error)
	DeletePlan(ctx context.Context, ts, id string) error

	// Task/Ticket CRUD
	CreateTask(ctx context.Context, ts string, req tc.CreateTicketRequest) (tc.Ticket, error)
	UpdateTask(ctx context.Context, ts, id string, req tc.UpdateTicketRequest) (tc.Ticket, error)
	DeleteTask(ctx context.Context, ts, id string) error

	// Graph queries
	ListChildren(ctx context.Context, ts, ticketID string) ([]tc.Ticket, error)
	ListBlocking(ctx context.Context, ts, ticketID string) ([]tc.Ticket, error)
}

// MemoryManager handles comment/memory operations.
type MemoryManager interface {
	AddComment(ctx context.Context, ts, ticketID, author, text string, tags []string) (tc.Comment, error)
	ListMemories(ctx context.Context, ts, ticketID, agentID string) ([]types.Memory, error)
}

// TrackingSetManager handles tracking set operations.
type TrackingSetManager interface {
	CreateTrackingSet(ctx context.Context, name string) (tc.TrackingSet, error)
	ListTrackingSets(ctx context.Context) ([]tc.TrackingSet, error)
	GetTrackingSet(ctx context.Context, ts string) (tc.TrackingSet, error)
	DeleteTrackingSet(ctx context.Context, ts string) error
}

// EdgeManager handles edge (relationship) operations between tickets.
type EdgeManager interface {
	ListEdges(ctx context.Context, ts string, params *tc.ListEdgesParams) ([]tc.Edge, error)
	AddEdge(ctx context.Context, ts string, req tc.EdgeRequest) (tc.Edge, error)
	RemoveEdge(ctx context.Context, ts string, req tc.EdgeRequest) error
}
