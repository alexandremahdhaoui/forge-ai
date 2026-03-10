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
	"fmt"

	"github.com/alexandremahdhaoui/forge-ai/internal/adapter"
	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
)

var _ PlanManager = (*planManager)(nil)

type planManager struct {
	client adapter.TrackerClient
}

// NewPlanManager creates a PlanManager with the given tracker client.
func NewPlanManager(client adapter.TrackerClient) PlanManager {
	return &planManager{client: client}
}

func (m *planManager) ListMetaPlans(ctx context.Context, ts string) ([]tc.MetaPlan, error) {
	return m.client.ListMetaPlans(ctx, ts)
}

func (m *planManager) GetMetaPlan(ctx context.Context, ts, id string) (tc.MetaPlan, error) {
	return m.client.GetMetaPlan(ctx, ts, id)
}

func (m *planManager) ListPlans(ctx context.Context, ts string) ([]tc.Plan, error) {
	return m.client.ListPlans(ctx, ts)
}

func (m *planManager) GetPlanState(ctx context.Context, ts, id string) (tc.Plan, error) {
	return m.client.GetPlan(ctx, ts, id)
}

func (m *planManager) ListTasks(ctx context.Context, ts string, filter adapter.TicketFilter) ([]tc.Ticket, error) {
	return m.client.ListTickets(ctx, ts, filter)
}

func (m *planManager) GetTask(ctx context.Context, ts, id string) (tc.Ticket, error) {
	return m.client.GetTicket(ctx, ts, id)
}

func (m *planManager) AssignTask(ctx context.Context, ts, ticketID, agentID string) (tc.Ticket, error) {
	ticket, err := m.client.GetTicket(ctx, ts, ticketID)
	if err != nil {
		return tc.Ticket{}, fmt.Errorf("getting ticket for assignment: %w", err)
	}
	title := ""
	if ticket.Title != nil {
		title = *ticket.Title
	}
	status := "in_progress"
	return m.client.UpdateTicket(ctx, ts, ticketID, tc.UpdateTicketRequest{
		Title:    title,
		Status:   &status,
		Assignee: &agentID,
	})
}

func (m *planManager) CompleteTask(ctx context.Context, ts, ticketID string) (tc.Ticket, error) {
	ticket, err := m.client.GetTicket(ctx, ts, ticketID)
	if err != nil {
		return tc.Ticket{}, fmt.Errorf("getting ticket for completion: %w", err)
	}
	title := ""
	if ticket.Title != nil {
		title = *ticket.Title
	}
	status := "completed"
	return m.client.UpdateTicket(ctx, ts, ticketID, tc.UpdateTicketRequest{
		Title:  title,
		Status: &status,
	})
}

// --- MetaPlan CRUD ---

func (m *planManager) CreateMetaPlan(ctx context.Context, ts string, req tc.CreateMetaPlanRequest) (tc.MetaPlan, error) {
	return m.client.CreateMetaPlan(ctx, ts, req)
}

func (m *planManager) UpdateMetaPlan(ctx context.Context, ts, id string, req tc.UpdateMetaPlanRequest) (tc.MetaPlan, error) {
	return m.client.UpdateMetaPlan(ctx, ts, id, req)
}

func (m *planManager) DeleteMetaPlan(ctx context.Context, ts, id string) error {
	return m.client.DeleteMetaPlan(ctx, ts, id)
}

// --- Plan CRUD ---

func (m *planManager) CreatePlan(ctx context.Context, ts string, req tc.CreatePlanRequest) (tc.Plan, error) {
	return m.client.CreatePlan(ctx, ts, req)
}

func (m *planManager) UpdatePlan(ctx context.Context, ts, id string, req tc.UpdatePlanRequest) (tc.Plan, error) {
	return m.client.UpdatePlan(ctx, ts, id, req)
}

func (m *planManager) DeletePlan(ctx context.Context, ts, id string) error {
	return m.client.DeletePlan(ctx, ts, id)
}

// --- Task/Ticket CRUD ---

func (m *planManager) CreateTask(ctx context.Context, ts string, req tc.CreateTicketRequest) (tc.Ticket, error) {
	return m.client.CreateTicket(ctx, ts, req)
}

func (m *planManager) UpdateTask(ctx context.Context, ts, id string, req tc.UpdateTicketRequest) (tc.Ticket, error) {
	return m.client.UpdateTicket(ctx, ts, id, req)
}

func (m *planManager) DeleteTask(ctx context.Context, ts, id string) error {
	return m.client.DeleteTicket(ctx, ts, id)
}

// --- Graph queries ---

func (m *planManager) ListChildren(ctx context.Context, ts, ticketID string) ([]tc.Ticket, error) {
	return m.client.GetChildren(ctx, ts, ticketID)
}

func (m *planManager) ListBlocking(ctx context.Context, ts, ticketID string) ([]tc.Ticket, error) {
	return m.client.GetBlocking(ctx, ts, ticketID)
}
