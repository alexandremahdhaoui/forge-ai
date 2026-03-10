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
	"fmt"
	"net/http"
	"strings"

	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
)

var _ TrackerClient = (*HTTPTrackerClient)(nil)

// HTTPTrackerClient implements TrackerClient using the generated oapi-codegen client.
type HTTPTrackerClient struct {
	client *tc.ClientWithResponses
}

// NewHTTPTrackerClient creates a TrackerClient backed by forge-tracker REST API.
func NewHTTPTrackerClient(baseURL string) (*HTTPTrackerClient, error) {
	client, err := tc.NewClientWithResponses(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return nil, fmt.Errorf("creating tracker client: %w", err)
	}
	return &HTTPTrackerClient{client: client}, nil
}

func (c *HTTPTrackerClient) ListTrackingSets(ctx context.Context) ([]tc.TrackingSet, error) {
	resp, err := c.client.ListTrackingSetsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing tracking sets: %w", err)
	}
	if resp.JSON200 == nil {
		return nil, apiError("listing tracking sets", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) ListMetaPlans(ctx context.Context, ts string) ([]tc.MetaPlan, error) {
	resp, err := c.client.ListMetaPlansWithResponse(ctx, ts)
	if err != nil {
		return nil, fmt.Errorf("listing meta-plans: %w", err)
	}
	if resp.JSON200 == nil {
		return nil, apiError("listing meta-plans", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) GetMetaPlan(ctx context.Context, ts, id string) (tc.MetaPlan, error) {
	resp, err := c.client.GetMetaPlanWithResponse(ctx, ts, id)
	if err != nil {
		return tc.MetaPlan{}, fmt.Errorf("getting meta-plan: %w", err)
	}
	if resp.JSON200 == nil {
		return tc.MetaPlan{}, apiError("getting meta-plan", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) ListPlans(ctx context.Context, ts string) ([]tc.Plan, error) {
	resp, err := c.client.ListPlansWithResponse(ctx, ts)
	if err != nil {
		return nil, fmt.Errorf("listing plans: %w", err)
	}
	if resp.JSON200 == nil {
		return nil, apiError("listing plans", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) GetPlan(ctx context.Context, ts, id string) (tc.Plan, error) {
	resp, err := c.client.GetPlanWithResponse(ctx, ts, id)
	if err != nil {
		return tc.Plan{}, fmt.Errorf("getting plan: %w", err)
	}
	if resp.JSON200 == nil {
		return tc.Plan{}, apiError("getting plan", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) ListTickets(ctx context.Context, ts string, filter TicketFilter) ([]tc.Ticket, error) {
	params := &tc.ListTicketsParams{}
	if filter.Status != "" {
		params.Status = &filter.Status
	}
	if filter.Assignee != "" {
		params.Assignee = &filter.Assignee
	}
	if len(filter.Labels) > 0 {
		joined := strings.Join(filter.Labels, ",")
		params.Labels = &joined
	}
	if filter.Priority != nil {
		params.Priority = filter.Priority
	}
	resp, err := c.client.ListTicketsWithResponse(ctx, ts, params)
	if err != nil {
		return nil, fmt.Errorf("listing tickets: %w", err)
	}
	if resp.JSON200 == nil {
		return nil, apiError("listing tickets", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) GetTicket(ctx context.Context, ts, id string) (tc.Ticket, error) {
	resp, err := c.client.GetTicketWithResponse(ctx, ts, id)
	if err != nil {
		return tc.Ticket{}, fmt.Errorf("getting ticket: %w", err)
	}
	if resp.JSON200 == nil {
		return tc.Ticket{}, apiError("getting ticket", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) UpdateTicket(ctx context.Context, ts, id string, req tc.UpdateTicketRequest) (tc.Ticket, error) {
	resp, err := c.client.UpdateTicketWithResponse(ctx, ts, id, tc.UpdateTicketJSONRequestBody(req))
	if err != nil {
		return tc.Ticket{}, fmt.Errorf("updating ticket: %w", err)
	}
	if resp.JSON200 == nil {
		return tc.Ticket{}, apiError("updating ticket", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) AddComment(ctx context.Context, ts, ticketID string, req tc.AddCommentRequest) (tc.Comment, error) {
	resp, err := c.client.AddCommentWithResponse(ctx, ts, ticketID, tc.AddCommentJSONRequestBody(req))
	if err != nil {
		return tc.Comment{}, fmt.Errorf("adding comment: %w", err)
	}
	if resp.JSON201 == nil {
		return tc.Comment{}, apiError("adding comment", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON201, nil
}

func (c *HTTPTrackerClient) GetChildren(ctx context.Context, ts, ticketID string) ([]tc.Ticket, error) {
	resp, err := c.client.GetChildrenWithResponse(ctx, ts, ticketID)
	if err != nil {
		return nil, fmt.Errorf("getting children: %w", err)
	}
	if resp.JSON200 == nil {
		return nil, apiError("getting children", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) GetBlocking(ctx context.Context, ts, ticketID string) ([]tc.Ticket, error) {
	resp, err := c.client.GetBlockingWithResponse(ctx, ts, ticketID)
	if err != nil {
		return nil, fmt.Errorf("getting blocking: %w", err)
	}
	if resp.JSON200 == nil {
		return nil, apiError("getting blocking", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) CreateTrackingSet(ctx context.Context, req tc.CreateTrackingSetRequest) (tc.TrackingSet, error) {
	resp, err := c.client.CreateTrackingSetWithResponse(ctx, tc.CreateTrackingSetJSONRequestBody(req))
	if err != nil {
		return tc.TrackingSet{}, fmt.Errorf("creating tracking set: %w", err)
	}
	if resp.JSON201 == nil {
		return tc.TrackingSet{}, apiError("creating tracking set", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON201, nil
}

func (c *HTTPTrackerClient) GetTrackingSet(ctx context.Context, ts string) (tc.TrackingSet, error) {
	resp, err := c.client.GetTrackingSetWithResponse(ctx, ts)
	if err != nil {
		return tc.TrackingSet{}, fmt.Errorf("getting tracking set: %w", err)
	}
	if resp.JSON200 == nil {
		return tc.TrackingSet{}, apiError("getting tracking set", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) DeleteTrackingSet(ctx context.Context, ts string) error {
	resp, err := c.client.DeleteTrackingSetWithResponse(ctx, ts)
	if err != nil {
		return fmt.Errorf("deleting tracking set: %w", err)
	}
	if resp.StatusCode() != http.StatusNoContent {
		return apiError("deleting tracking set", resp.StatusCode(), resp.Body)
	}
	return nil
}

func (c *HTTPTrackerClient) CreateMetaPlan(ctx context.Context, ts string, req tc.CreateMetaPlanRequest) (tc.MetaPlan, error) {
	resp, err := c.client.CreateMetaPlanWithResponse(ctx, ts, tc.CreateMetaPlanJSONRequestBody(req))
	if err != nil {
		return tc.MetaPlan{}, fmt.Errorf("creating meta-plan: %w", err)
	}
	if resp.JSON201 == nil {
		return tc.MetaPlan{}, apiError("creating meta-plan", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON201, nil
}

func (c *HTTPTrackerClient) UpdateMetaPlan(ctx context.Context, ts, id string, req tc.UpdateMetaPlanRequest) (tc.MetaPlan, error) {
	resp, err := c.client.UpdateMetaPlanWithResponse(ctx, ts, id, tc.UpdateMetaPlanJSONRequestBody(req))
	if err != nil {
		return tc.MetaPlan{}, fmt.Errorf("updating meta-plan: %w", err)
	}
	if resp.JSON200 == nil {
		return tc.MetaPlan{}, apiError("updating meta-plan", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) DeleteMetaPlan(ctx context.Context, ts, id string) error {
	resp, err := c.client.DeleteMetaPlanWithResponse(ctx, ts, id)
	if err != nil {
		return fmt.Errorf("deleting meta-plan: %w", err)
	}
	if resp.StatusCode() != http.StatusNoContent {
		return apiError("deleting meta-plan", resp.StatusCode(), resp.Body)
	}
	return nil
}

func (c *HTTPTrackerClient) CreatePlan(ctx context.Context, ts string, req tc.CreatePlanRequest) (tc.Plan, error) {
	resp, err := c.client.CreatePlanWithResponse(ctx, ts, tc.CreatePlanJSONRequestBody(req))
	if err != nil {
		return tc.Plan{}, fmt.Errorf("creating plan: %w", err)
	}
	if resp.JSON201 == nil {
		return tc.Plan{}, apiError("creating plan", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON201, nil
}

func (c *HTTPTrackerClient) UpdatePlan(ctx context.Context, ts, id string, req tc.UpdatePlanRequest) (tc.Plan, error) {
	resp, err := c.client.UpdatePlanWithResponse(ctx, ts, id, tc.UpdatePlanJSONRequestBody(req))
	if err != nil {
		return tc.Plan{}, fmt.Errorf("updating plan: %w", err)
	}
	if resp.JSON200 == nil {
		return tc.Plan{}, apiError("updating plan", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) DeletePlan(ctx context.Context, ts, id string) error {
	resp, err := c.client.DeletePlanWithResponse(ctx, ts, id)
	if err != nil {
		return fmt.Errorf("deleting plan: %w", err)
	}
	if resp.StatusCode() != http.StatusNoContent {
		return apiError("deleting plan", resp.StatusCode(), resp.Body)
	}
	return nil
}

func (c *HTTPTrackerClient) CreateTicket(ctx context.Context, ts string, req tc.CreateTicketRequest) (tc.Ticket, error) {
	resp, err := c.client.CreateTicketWithResponse(ctx, ts, tc.CreateTicketJSONRequestBody(req))
	if err != nil {
		return tc.Ticket{}, fmt.Errorf("creating ticket: %w", err)
	}
	if resp.JSON201 == nil {
		return tc.Ticket{}, apiError("creating ticket", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON201, nil
}

func (c *HTTPTrackerClient) DeleteTicket(ctx context.Context, ts, id string) error {
	resp, err := c.client.DeleteTicketWithResponse(ctx, ts, id)
	if err != nil {
		return fmt.Errorf("deleting ticket: %w", err)
	}
	if resp.StatusCode() != http.StatusNoContent {
		return apiError("deleting ticket", resp.StatusCode(), resp.Body)
	}
	return nil
}

func (c *HTTPTrackerClient) ListEdges(ctx context.Context, ts string, params *tc.ListEdgesParams) ([]tc.Edge, error) {
	resp, err := c.client.ListEdgesWithResponse(ctx, ts, params)
	if err != nil {
		return nil, fmt.Errorf("listing edges: %w", err)
	}
	if resp.JSON200 == nil {
		return nil, apiError("listing edges", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON200, nil
}

func (c *HTTPTrackerClient) AddEdge(ctx context.Context, ts string, req tc.EdgeRequest) (tc.Edge, error) {
	resp, err := c.client.AddEdgeWithResponse(ctx, ts, tc.AddEdgeJSONRequestBody(req))
	if err != nil {
		return tc.Edge{}, fmt.Errorf("adding edge: %w", err)
	}
	if resp.JSON201 == nil {
		return tc.Edge{}, apiError("adding edge", resp.StatusCode(), resp.Body)
	}
	return *resp.JSON201, nil
}

func (c *HTTPTrackerClient) RemoveEdge(ctx context.Context, ts string, req tc.EdgeRequest) error {
	resp, err := c.client.RemoveEdgeWithResponse(ctx, ts, tc.RemoveEdgeJSONRequestBody(req))
	if err != nil {
		return fmt.Errorf("removing edge: %w", err)
	}
	if resp.StatusCode() != http.StatusNoContent {
		return apiError("removing edge", resp.StatusCode(), resp.Body)
	}
	return nil
}

func apiError(operation string, statusCode int, body []byte) error {
	return fmt.Errorf("%s: unexpected status %d: %s", operation, statusCode, string(body))
}
