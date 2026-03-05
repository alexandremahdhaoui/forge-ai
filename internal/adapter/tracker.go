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

package adapter

import (
	"context"

	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
)

// TrackerClient abstracts all forge-tracker REST API operations that forge-ai needs.
type TrackerClient interface {
	ListTrackingSets(ctx context.Context) ([]tc.TrackingSet, error)
	ListMetaPlans(ctx context.Context, ts string) ([]tc.MetaPlan, error)
	GetMetaPlan(ctx context.Context, ts, id string) (tc.MetaPlan, error)
	ListPlans(ctx context.Context, ts string) ([]tc.Plan, error)
	GetPlan(ctx context.Context, ts, id string) (tc.Plan, error)
	ListTickets(ctx context.Context, ts string, filter TicketFilter) ([]tc.Ticket, error)
	GetTicket(ctx context.Context, ts, id string) (tc.Ticket, error)
	UpdateTicket(ctx context.Context, ts, id string, req tc.UpdateTicketRequest) (tc.Ticket, error)
	AddComment(ctx context.Context, ts, ticketID string, req tc.AddCommentRequest) (tc.Comment, error)
	GetChildren(ctx context.Context, ts, ticketID string) ([]tc.Ticket, error)
	GetBlocking(ctx context.Context, ts, ticketID string) ([]tc.Ticket, error)
}

// TicketFilter holds optional query parameters for ListTickets.
type TicketFilter struct {
	Status   string
	Assignee string
	Labels   []string
	Priority *int
}
