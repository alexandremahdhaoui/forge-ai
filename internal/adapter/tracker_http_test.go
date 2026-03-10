//go:build unit

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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
)

func strPtr(s string) *string       { return &s }
func intPtr(i int) *int             { return &i }
func timePtr(t time.Time) *time.Time { return &t }

func newTestClient(t *testing.T, handler http.Handler) *HTTPTrackerClient {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client, err := NewHTTPTrackerClient(srv.URL)
	if err != nil {
		t.Fatalf("creating test client: %v", err)
	}
	return client
}

func TestListTrackingSets(t *testing.T) {
	expected := []tc.TrackingSet{
		{Name: strPtr("ts1")},
		{Name: strPtr("ts2")},
	}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.ListTrackingSets(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 tracking sets, got %d", len(result))
	}
	if *result[0].Name != "ts1" {
		t.Errorf("expected ts1, got %s", *result[0].Name)
	}
	if *result[1].Name != "ts2" {
		t.Errorf("expected ts2, got %s", *result[1].Name)
	}
}

func TestListMetaPlans(t *testing.T) {
	expected := []tc.MetaPlan{
		{Id: strPtr("mp1"), Title: strPtr("Meta Plan 1")},
		{Id: strPtr("mp2"), Title: strPtr("Meta Plan 2")},
	}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/metaplans" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.ListMetaPlans(context.Background(), "myts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 meta-plans, got %d", len(result))
	}
	if *result[0].Id != "mp1" {
		t.Errorf("expected mp1, got %s", *result[0].Id)
	}
}

func TestGetMetaPlan(t *testing.T) {
	expected := tc.MetaPlan{
		Id:          strPtr("mp1"),
		Title:       strPtr("Meta Plan 1"),
		Description: strPtr("A test meta-plan"),
		Status:      strPtr("active"),
	}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/metaplans/mp1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.GetMetaPlan(context.Background(), "myts", "mp1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.Id != "mp1" {
		t.Errorf("expected mp1, got %s", *result.Id)
	}
	if *result.Title != "Meta Plan 1" {
		t.Errorf("expected Meta Plan 1, got %s", *result.Title)
	}
}

func TestListPlans(t *testing.T) {
	expected := []tc.Plan{
		{Id: strPtr("p1"), Title: strPtr("Plan 1")},
		{Id: strPtr("p2"), Title: strPtr("Plan 2")},
	}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/plans" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.ListPlans(context.Background(), "myts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 plans, got %d", len(result))
	}
	if *result[0].Id != "p1" {
		t.Errorf("expected p1, got %s", *result[0].Id)
	}
}

func TestGetPlan(t *testing.T) {
	tasks := []string{"t1", "t2", "t3"}
	expected := tc.Plan{
		Id:    strPtr("p1"),
		Title: strPtr("Plan 1"),
		Tasks: &tasks,
	}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/plans/p1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.GetPlan(context.Background(), "myts", "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.Id != "p1" {
		t.Errorf("expected p1, got %s", *result.Id)
	}
	if result.Tasks == nil || len(*result.Tasks) != 3 {
		t.Errorf("expected 3 tasks, got %v", result.Tasks)
	}
}

func TestListTickets(t *testing.T) {
	expected := []tc.Ticket{
		{Id: strPtr("t1"), Title: strPtr("Task 1")},
	}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/tickets" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if got := r.URL.Query().Get("status"); got != "pending" {
			t.Errorf("expected status=pending, got %q", got)
		}
		if got := r.URL.Query().Get("assignee"); got != "agent-1" {
			t.Errorf("expected assignee=agent-1, got %q", got)
		}
		if got := r.URL.Query().Get("labels"); got != "bug,feature" {
			t.Errorf("expected labels=bug,feature, got %q", got)
		}
		if got := r.URL.Query().Get("priority"); got != "1" {
			t.Errorf("expected priority=1, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.ListTickets(context.Background(), "myts", TicketFilter{
		Status:   "pending",
		Assignee: "agent-1",
		Labels:   []string{"bug", "feature"},
		Priority: intPtr(1),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 ticket, got %d", len(result))
	}
	if *result[0].Id != "t1" {
		t.Errorf("expected t1, got %s", *result[0].Id)
	}
}

func TestGetTicket(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	comments := []tc.Comment{
		{Author: strPtr("user1"), Text: strPtr("A comment"), Timestamp: timePtr(now)},
	}
	expected := tc.Ticket{
		Id:       strPtr("t1"),
		Title:    strPtr("Task 1"),
		Status:   strPtr("in_progress"),
		Assignee: strPtr("agent-1"),
		Comments: &comments,
	}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/tickets/t1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.GetTicket(context.Background(), "myts", "t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.Id != "t1" {
		t.Errorf("expected t1, got %s", *result.Id)
	}
	if *result.Status != "in_progress" {
		t.Errorf("expected in_progress, got %s", *result.Status)
	}
	if result.Comments == nil || len(*result.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %v", result.Comments)
	}
	if *(*result.Comments)[0].Author != "user1" {
		t.Errorf("expected user1, got %s", *(*result.Comments)[0].Author)
	}
}

func TestUpdateTicket(t *testing.T) {
	expected := tc.Ticket{
		Id:       strPtr("t1"),
		Title:    strPtr("Updated Task"),
		Status:   strPtr("in_progress"),
		Assignee: strPtr("agent-2"),
	}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/tickets/t1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT method, got %s", r.Method)
		}
		var body tc.UpdateTicketRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		if body.Title != "Updated Task" {
			t.Errorf("expected title Updated Task, got %s", body.Title)
		}
		if body.Status == nil || *body.Status != "in_progress" {
			t.Errorf("expected status in_progress, got %v", body.Status)
		}
		if body.Assignee == nil || *body.Assignee != "agent-2" {
			t.Errorf("expected assignee agent-2, got %v", body.Assignee)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.UpdateTicket(context.Background(), "myts", "t1", tc.UpdateTicketRequest{
		Title:    "Updated Task",
		Status:   strPtr("in_progress"),
		Assignee: strPtr("agent-2"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.Id != "t1" {
		t.Errorf("expected t1, got %s", *result.Id)
	}
	if *result.Status != "in_progress" {
		t.Errorf("expected in_progress, got %s", *result.Status)
	}
}

func TestAddComment(t *testing.T) {
	expected := tc.Comment{
		Author:    strPtr("agent-1"),
		Text:      strPtr("Test comment"),
		Timestamp: timePtr(time.Now().Truncate(time.Second)),
	}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/tickets/t1/comments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		var body tc.AddCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		if body.Author != "agent-1" {
			t.Errorf("expected author agent-1, got %s", body.Author)
		}
		if body.Text != "Test comment" {
			t.Errorf("expected text Test comment, got %s", body.Text)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		writeJSON(t, w, expected)
	}))
	result, err := client.AddComment(context.Background(), "myts", "t1", tc.AddCommentRequest{
		Author: "agent-1",
		Text:   "Test comment",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.Author != "agent-1" {
		t.Errorf("expected agent-1, got %s", *result.Author)
	}
	if *result.Text != "Test comment" {
		t.Errorf("expected Test comment, got %s", *result.Text)
	}
}

func TestGetChildren(t *testing.T) {
	expected := []tc.Ticket{
		{Id: strPtr("child-1"), Title: strPtr("Child 1")},
		{Id: strPtr("child-2"), Title: strPtr("Child 2")},
	}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/tickets/t1/children" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.GetChildren(context.Background(), "myts", "t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 children, got %d", len(result))
	}
	if *result[0].Id != "child-1" {
		t.Errorf("expected child-1, got %s", *result[0].Id)
	}
}

func TestGetBlocking(t *testing.T) {
	expected := []tc.Ticket{
		{Id: strPtr("blocker-1"), Title: strPtr("Blocker 1")},
	}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/tickets/t1/blocking" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.GetBlocking(context.Background(), "myts", "t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 blocker, got %d", len(result))
	}
	if *result[0].Id != "blocker-1" {
		t.Errorf("expected blocker-1, got %s", *result[0].Id)
	}
}

func TestAPIError_404(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		writeJSON(t, w, map[string]string{"error": "not found"})
	}))
	_, err := client.GetTicket(context.Background(), "ts1", "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "404") {
		t.Errorf("expected error to contain 404, got: %s", got)
	}
}

func TestAPIError_500(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(t, w, map[string]string{"error": "internal server error"})
	}))
	_, err := client.ListTrackingSets(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "500") {
		t.Errorf("expected error to contain 500, got: %s", got)
	}
}

func TestCreateTrackingSet(t *testing.T) {
	expected := tc.TrackingSet{Name: strPtr("ts1")}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		var body tc.CreateTrackingSetRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		if body.Name != "ts1" {
			t.Errorf("expected name ts1, got %s", body.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		writeJSON(t, w, expected)
	}))
	result, err := client.CreateTrackingSet(context.Background(), tc.CreateTrackingSetRequest{Name: "ts1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.Name != "ts1" {
		t.Errorf("expected ts1, got %s", *result.Name)
	}
}

func TestGetTrackingSet(t *testing.T) {
	expected := tc.TrackingSet{Name: strPtr("myts")}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.GetTrackingSet(context.Background(), "myts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.Name != "myts" {
		t.Errorf("expected myts, got %s", *result.Name)
	}
}

func TestDeleteTrackingSet(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	err := client.DeleteTrackingSet(context.Background(), "myts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateMetaPlan(t *testing.T) {
	expected := tc.MetaPlan{Id: strPtr("mp1"), Title: strPtr("Meta Plan 1")}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/metaplans" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		var body tc.CreateMetaPlanRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		if body.Id != "mp1" {
			t.Errorf("expected id mp1, got %s", body.Id)
		}
		if body.Title != "Meta Plan 1" {
			t.Errorf("expected title Meta Plan 1, got %s", body.Title)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		writeJSON(t, w, expected)
	}))
	result, err := client.CreateMetaPlan(context.Background(), "myts", tc.CreateMetaPlanRequest{
		Id:    "mp1",
		Title: "Meta Plan 1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.Id != "mp1" {
		t.Errorf("expected mp1, got %s", *result.Id)
	}
	if *result.Title != "Meta Plan 1" {
		t.Errorf("expected Meta Plan 1, got %s", *result.Title)
	}
}

func TestUpdateMetaPlan(t *testing.T) {
	expected := tc.MetaPlan{Id: strPtr("mp1"), Title: strPtr("Updated")}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/metaplans/mp1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT method, got %s", r.Method)
		}
		var body tc.UpdateMetaPlanRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		if body.Title != "Updated" {
			t.Errorf("expected title Updated, got %s", body.Title)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.UpdateMetaPlan(context.Background(), "myts", "mp1", tc.UpdateMetaPlanRequest{
		Title: "Updated",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.Id != "mp1" {
		t.Errorf("expected mp1, got %s", *result.Id)
	}
	if *result.Title != "Updated" {
		t.Errorf("expected Updated, got %s", *result.Title)
	}
}

func TestDeleteMetaPlan(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/metaplans/mp1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	err := client.DeleteMetaPlan(context.Background(), "myts", "mp1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreatePlan(t *testing.T) {
	expected := tc.Plan{Id: strPtr("p1"), Title: strPtr("Plan 1")}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/plans" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		var body tc.CreatePlanRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		if body.Id != "p1" {
			t.Errorf("expected id p1, got %s", body.Id)
		}
		if body.Title != "Plan 1" {
			t.Errorf("expected title Plan 1, got %s", body.Title)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		writeJSON(t, w, expected)
	}))
	result, err := client.CreatePlan(context.Background(), "myts", tc.CreatePlanRequest{
		Id:    "p1",
		Title: "Plan 1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.Id != "p1" {
		t.Errorf("expected p1, got %s", *result.Id)
	}
	if *result.Title != "Plan 1" {
		t.Errorf("expected Plan 1, got %s", *result.Title)
	}
}

func TestUpdatePlan(t *testing.T) {
	expected := tc.Plan{Id: strPtr("p1"), Title: strPtr("Updated")}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/plans/p1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT method, got %s", r.Method)
		}
		var body tc.UpdatePlanRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		if body.Title != "Updated" {
			t.Errorf("expected title Updated, got %s", body.Title)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	result, err := client.UpdatePlan(context.Background(), "myts", "p1", tc.UpdatePlanRequest{
		Title: "Updated",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.Id != "p1" {
		t.Errorf("expected p1, got %s", *result.Id)
	}
	if *result.Title != "Updated" {
		t.Errorf("expected Updated, got %s", *result.Title)
	}
}

func TestDeletePlan(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/plans/p1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	err := client.DeletePlan(context.Background(), "myts", "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateTicket(t *testing.T) {
	expected := tc.Ticket{Id: strPtr("t1"), Title: strPtr("Task 1")}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/tickets" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		var body tc.CreateTicketRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		if body.Id != "t1" {
			t.Errorf("expected id t1, got %s", body.Id)
		}
		if body.Title != "Task 1" {
			t.Errorf("expected title Task 1, got %s", body.Title)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		writeJSON(t, w, expected)
	}))
	result, err := client.CreateTicket(context.Background(), "myts", tc.CreateTicketRequest{
		Id:    "t1",
		Title: "Task 1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.Id != "t1" {
		t.Errorf("expected t1, got %s", *result.Id)
	}
	if *result.Title != "Task 1" {
		t.Errorf("expected Task 1, got %s", *result.Title)
	}
}

func TestDeleteTicket(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/tickets/t1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	err := client.DeleteTicket(context.Background(), "myts", "t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListEdges(t *testing.T) {
	expected := []tc.Edge{
		{From: strPtr("t1"), To: strPtr("t2"), Type: strPtr("blocks")},
	}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/edges" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if got := r.URL.Query().Get("ticket"); got != "t1" {
			t.Errorf("expected ticket=t1, got %q", got)
		}
		if got := r.URL.Query().Get("type"); got != "blocks" {
			t.Errorf("expected type=blocks, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, expected)
	}))
	ticket := "t1"
	edgeType := "blocks"
	result, err := client.ListEdges(context.Background(), "myts", &tc.ListEdgesParams{
		Ticket: &ticket,
		Type:   &edgeType,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(result))
	}
	if *result[0].From != "t1" {
		t.Errorf("expected from t1, got %s", *result[0].From)
	}
	if *result[0].To != "t2" {
		t.Errorf("expected to t2, got %s", *result[0].To)
	}
	if *result[0].Type != "blocks" {
		t.Errorf("expected type blocks, got %s", *result[0].Type)
	}
}

func TestAddEdge(t *testing.T) {
	expected := tc.Edge{From: strPtr("t1"), To: strPtr("t2"), Type: strPtr("blocks")}
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/edges" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		var body tc.EdgeRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		if body.From != "t1" {
			t.Errorf("expected from t1, got %s", body.From)
		}
		if body.To != "t2" {
			t.Errorf("expected to t2, got %s", body.To)
		}
		if body.Type != "blocks" {
			t.Errorf("expected type blocks, got %s", body.Type)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		writeJSON(t, w, expected)
	}))
	result, err := client.AddEdge(context.Background(), "myts", tc.EdgeRequest{
		From: "t1",
		To:   "t2",
		Type: "blocks",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *result.From != "t1" {
		t.Errorf("expected from t1, got %s", *result.From)
	}
	if *result.To != "t2" {
		t.Errorf("expected to t2, got %s", *result.To)
	}
	if *result.Type != "blocks" {
		t.Errorf("expected type blocks, got %s", *result.Type)
	}
}

func TestRemoveEdge(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tracking-sets/myts/edges" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		var body tc.EdgeRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		if body.From != "t1" {
			t.Errorf("expected from t1, got %s", body.From)
		}
		if body.To != "t2" {
			t.Errorf("expected to t2, got %s", body.To)
		}
		if body.Type != "blocks" {
			t.Errorf("expected type blocks, got %s", body.Type)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	err := client.RemoveEdge(context.Background(), "myts", tc.EdgeRequest{
		From: "t1",
		To:   "t2",
		Type: "blocks",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func writeJSON(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Fatalf("encoding response: %v", err)
	}
}

// contains checks if s contains substr. Using a helper to avoid importing strings.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

