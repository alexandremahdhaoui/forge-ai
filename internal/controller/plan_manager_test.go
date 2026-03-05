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
