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

package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/alexandremahdhaoui/forge-ai/internal/adapter"
	"github.com/alexandremahdhaoui/forge-ai/internal/util/mocks/mockadapter"
	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func strPtr(s string) *string { return &s }

func TestListMetaPlans(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	expected := []tc.MetaPlan{{Id: strPtr("mp-1")}}
	mockClient.EXPECT().ListMetaPlans(mock.Anything, "ts1").Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.ListMetaPlans(context.Background(), "ts1")
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetMetaPlan(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	expected := tc.MetaPlan{Id: strPtr("mp-1")}
	mockClient.EXPECT().GetMetaPlan(mock.Anything, "ts1", "mp-1").Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.GetMetaPlan(context.Background(), "ts1", "mp-1")
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestListPlans(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	expected := []tc.Plan{{Id: strPtr("plan-1")}}
	mockClient.EXPECT().ListPlans(mock.Anything, "ts1").Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.ListPlans(context.Background(), "ts1")
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetPlanState(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	expected := tc.Plan{Id: strPtr("plan-1")}
	mockClient.EXPECT().GetPlan(mock.Anything, "ts1", "plan-1").Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.GetPlanState(context.Background(), "ts1", "plan-1")
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestListTasks(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	filter := adapter.TicketFilter{Status: "open"}
	expected := []tc.Ticket{{Id: strPtr("task-1")}}
	mockClient.EXPECT().ListTickets(mock.Anything, "ts1", filter).Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.ListTasks(context.Background(), "ts1", filter)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetTask(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	expected := tc.Ticket{Id: strPtr("task-1"), Title: strPtr("My Task")}
	mockClient.EXPECT().GetTicket(mock.Anything, "ts1", "task-1").Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.GetTask(context.Background(), "ts1", "task-1")
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestAssignTask(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)

	// First call: GetTicket to retrieve current title
	mockClient.EXPECT().GetTicket(mock.Anything, "ts1", "task-1").
		Return(tc.Ticket{Title: strPtr("My Task"), Id: strPtr("task-1")}, nil)

	// Second call: UpdateTicket with assignee and status
	mockClient.EXPECT().UpdateTicket(mock.Anything, "ts1", "task-1", mock.MatchedBy(func(req tc.UpdateTicketRequest) bool {
		return req.Title == "My Task" &&
			req.Assignee != nil && *req.Assignee == "agent-1" &&
			req.Status != nil && *req.Status == "in_progress"
	})).Return(tc.Ticket{
		Id:       strPtr("task-1"),
		Title:    strPtr("My Task"),
		Assignee: strPtr("agent-1"),
		Status:   strPtr("in_progress"),
	}, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.AssignTask(context.Background(), "ts1", "task-1", "agent-1")
	require.NoError(t, err)
	assert.Equal(t, "agent-1", *result.Assignee)
	assert.Equal(t, "in_progress", *result.Status)
}

func TestCompleteTask(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)

	// First call: GetTicket to retrieve current title
	mockClient.EXPECT().GetTicket(mock.Anything, "ts1", "task-1").
		Return(tc.Ticket{Title: strPtr("My Task"), Id: strPtr("task-1")}, nil)

	// Second call: UpdateTicket with status=completed
	mockClient.EXPECT().UpdateTicket(mock.Anything, "ts1", "task-1", mock.MatchedBy(func(req tc.UpdateTicketRequest) bool {
		return req.Title == "My Task" &&
			req.Status != nil && *req.Status == "completed" &&
			req.Assignee == nil
	})).Return(tc.Ticket{
		Id:     strPtr("task-1"),
		Title:  strPtr("My Task"),
		Status: strPtr("completed"),
	}, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.CompleteTask(context.Background(), "ts1", "task-1")
	require.NoError(t, err)
	assert.Equal(t, "completed", *result.Status)
}

func TestAssignTask_Error(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	mockClient.EXPECT().GetTicket(mock.Anything, "ts1", "task-1").
		Return(tc.Ticket{}, errors.New("not found"))

	mgr := NewPlanManager(mockClient)
	_, err := mgr.AssignTask(context.Background(), "ts1", "task-1", "agent-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestListMetaPlans_Error(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	mockClient.EXPECT().ListMetaPlans(mock.Anything, "ts1").
		Return(nil, errors.New("connection refused"))

	mgr := NewPlanManager(mockClient)
	_, err := mgr.ListMetaPlans(context.Background(), "ts1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
}

// --- Create methods ---

func TestCreateMetaPlan(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	req := tc.CreateMetaPlanRequest{Id: "mp1", Title: "Meta Plan 1"}
	expected := tc.MetaPlan{Id: strPtr("mp1"), Title: strPtr("Meta Plan 1")}
	mockClient.EXPECT().CreateMetaPlan(mock.Anything, "ts1", req).Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.CreateMetaPlan(context.Background(), "ts1", req)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreatePlan(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	req := tc.CreatePlanRequest{Id: "p1", Title: "Plan 1"}
	expected := tc.Plan{Id: strPtr("p1"), Title: strPtr("Plan 1")}
	mockClient.EXPECT().CreatePlan(mock.Anything, "ts1", req).Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.CreatePlan(context.Background(), "ts1", req)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateTask(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	req := tc.CreateTicketRequest{Id: "task-1", Title: "Task 1"}
	expected := tc.Ticket{Id: strPtr("task-1"), Title: strPtr("Task 1")}
	mockClient.EXPECT().CreateTicket(mock.Anything, "ts1", req).Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.CreateTask(context.Background(), "ts1", req)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

// --- Update methods ---

func TestUpdateMetaPlan(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	req := tc.UpdateMetaPlanRequest{Title: "Updated"}
	expected := tc.MetaPlan{Id: strPtr("mp1"), Title: strPtr("Updated")}
	mockClient.EXPECT().UpdateMetaPlan(mock.Anything, "ts1", "mp1", req).Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.UpdateMetaPlan(context.Background(), "ts1", "mp1", req)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestUpdatePlan(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	req := tc.UpdatePlanRequest{Title: "Updated"}
	expected := tc.Plan{Id: strPtr("p1"), Title: strPtr("Updated")}
	mockClient.EXPECT().UpdatePlan(mock.Anything, "ts1", "p1", req).Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.UpdatePlan(context.Background(), "ts1", "p1", req)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestUpdateTask(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	req := tc.UpdateTicketRequest{Title: "Updated Task"}
	expected := tc.Ticket{Id: strPtr("task-1"), Title: strPtr("Updated Task")}
	mockClient.EXPECT().UpdateTicket(mock.Anything, "ts1", "task-1", req).Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.UpdateTask(context.Background(), "ts1", "task-1", req)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

// --- Delete methods ---

func TestDeleteMetaPlan(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	mockClient.EXPECT().DeleteMetaPlan(mock.Anything, "ts1", "mp1").Return(nil)

	mgr := NewPlanManager(mockClient)
	err := mgr.DeleteMetaPlan(context.Background(), "ts1", "mp1")
	require.NoError(t, err)
}

func TestDeletePlan(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	mockClient.EXPECT().DeletePlan(mock.Anything, "ts1", "p1").Return(nil)

	mgr := NewPlanManager(mockClient)
	err := mgr.DeletePlan(context.Background(), "ts1", "p1")
	require.NoError(t, err)
}

func TestDeleteTask(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	mockClient.EXPECT().DeleteTicket(mock.Anything, "ts1", "task-1").Return(nil)

	mgr := NewPlanManager(mockClient)
	err := mgr.DeleteTask(context.Background(), "ts1", "task-1")
	require.NoError(t, err)
}

// --- Graph query methods ---

func TestListChildren(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	expected := []tc.Ticket{{Id: strPtr("child-1")}}
	mockClient.EXPECT().GetChildren(mock.Anything, "ts1", "t1").Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.ListChildren(context.Background(), "ts1", "t1")
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestListBlocking(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	expected := []tc.Ticket{{Id: strPtr("blocker-1")}}
	mockClient.EXPECT().GetBlocking(mock.Anything, "ts1", "t1").Return(expected, nil)

	mgr := NewPlanManager(mockClient)
	result, err := mgr.ListBlocking(context.Background(), "ts1", "t1")
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
