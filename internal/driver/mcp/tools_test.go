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

package mcp

import (
	"context"
	"errors"
	"testing"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/alexandremahdhaoui/forge-ai/internal/adapter"
	"github.com/alexandremahdhaoui/forge-ai/internal/types"
	"github.com/alexandremahdhaoui/forge-ai/internal/util/mocks/mockcontroller"
	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
	"github.com/alexandremahdhaoui/forge/pkg/mcpserver"
)

func strPtr(s string) *string { return &s }

func TestRegisterTools(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockMemMgr := mockcontroller.NewMockMemoryManager(t)
	mockTsMgr := mockcontroller.NewMockTrackingSetManager(t)
	mockEdgeMgr := mockcontroller.NewMockEdgeManager(t)
	server := mcpserver.New("test-forge-ai", "test")

	require.NotPanics(t, func() {
		RegisterTools(server, mockPlanMgr, mockMemMgr, mockTsMgr, mockEdgeMgr)
	})
}

func TestJsonResult(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		input := map[string]string{"key": "value"}
		result, extra, err := jsonResult(input)
		require.NoError(t, err)
		require.Nil(t, extra)
		require.NotNil(t, result)
		require.Len(t, result.Content, 1)

		content, ok := result.Content[0].(*gomcp.TextContent)
		require.True(t, ok, "expected *gomcp.TextContent, got %T", result.Content[0])
		assert.Contains(t, content.Text, `"key"`)
		assert.Contains(t, content.Text, `"value"`)
	})

	t.Run("struct", func(t *testing.T) {
		type sample struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		result, _, err := jsonResult(sample{Name: "test", Age: 42})
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Content, 1)

		content, ok := result.Content[0].(*gomcp.TextContent)
		require.True(t, ok)
		assert.Contains(t, content.Text, `"name":"test"`)
		assert.Contains(t, content.Text, `"age":42`)
	})

	t.Run("marshal_error", func(t *testing.T) {
		_, _, err := jsonResult(make(chan int))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "marshaling result")
	})

	t.Run("nil_input", func(t *testing.T) {
		result, _, err := jsonResult(nil)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Content, 1)

		content, ok := result.Content[0].(*gomcp.TextContent)
		require.True(t, ok)
		assert.Equal(t, "null", content.Text)
	})
}

func TestErrResult(t *testing.T) {
	testErr := errors.New("test error")
	result, extra, err := errResult(testErr)
	require.Error(t, err)
	assert.Equal(t, "test error", err.Error())
	assert.Nil(t, result)
	assert.Nil(t, extra)
}

func TestTextResult(t *testing.T) {
	result, extra, err := textResult("deleted plan p1")
	require.NoError(t, err)
	require.Nil(t, extra)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)

	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok, "expected *gomcp.TextContent, got %T", result.Content[0])
	assert.Equal(t, "deleted plan p1", content.Text)
}

// --- PlanManager handler tests (19) ---

func TestHandleListMetaPlans(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := []tc.MetaPlan{{Id: strPtr("mp-1")}}
	mockPlanMgr.EXPECT().ListMetaPlans(mock.Anything, "ts1").Return(expected, nil)

	result, _, err := handleListMetaPlans(context.Background(), listMetaPlansInput{TS: "ts1"}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"mp-1"`)
}

func TestHandleListMetaPlans_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().ListMetaPlans(mock.Anything, "ts1").Return(nil, errors.New("not found"))

	_, _, err := handleListMetaPlans(context.Background(), listMetaPlansInput{TS: "ts1"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandleGetMetaPlan(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := tc.MetaPlan{Id: strPtr("mp-1")}
	mockPlanMgr.EXPECT().GetMetaPlan(mock.Anything, "ts1", "mp-1").Return(expected, nil)

	result, _, err := handleGetMetaPlan(context.Background(), getMetaPlanInput{TS: "ts1", ID: "mp-1"}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"mp-1"`)
}

func TestHandleGetMetaPlan_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().GetMetaPlan(mock.Anything, "ts1", "mp-1").Return(tc.MetaPlan{}, errors.New("not found"))

	_, _, err := handleGetMetaPlan(context.Background(), getMetaPlanInput{TS: "ts1", ID: "mp-1"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandleListPlans(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := []tc.Plan{{Id: strPtr("p1")}}
	mockPlanMgr.EXPECT().ListPlans(mock.Anything, "ts1").Return(expected, nil)

	result, _, err := handleListPlans(context.Background(), listPlansInput{TS: "ts1"}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"p1"`)
}

func TestHandleListPlans_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().ListPlans(mock.Anything, "ts1").Return(nil, errors.New("fail"))

	_, _, err := handleListPlans(context.Background(), listPlansInput{TS: "ts1"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail")
}

func TestHandleGetPlanState(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := tc.Plan{Id: strPtr("p1"), Title: strPtr("Plan 1")}
	mockPlanMgr.EXPECT().GetPlanState(mock.Anything, "ts1", "p1").Return(expected, nil)

	result, _, err := handleGetPlanState(context.Background(), getPlanStateInput{TS: "ts1", ID: "p1"}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"p1"`)
	assert.Contains(t, content.Text, `"Plan 1"`)
}

func TestHandleGetPlanState_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().GetPlanState(mock.Anything, "ts1", "p1").Return(tc.Plan{}, errors.New("not found"))

	_, _, err := handleGetPlanState(context.Background(), getPlanStateInput{TS: "ts1", ID: "p1"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandleListTasks(t *testing.T) {
	t.Run("both filters set", func(t *testing.T) {
		mockPlanMgr := mockcontroller.NewMockPlanManager(t)
		status := "pending"
		assignee := "agent-1"
		expected := []tc.Ticket{{Id: strPtr("t1")}}
		mockPlanMgr.EXPECT().ListTasks(mock.Anything, "ts1", adapter.TicketFilter{
			Status: "pending", Assignee: "agent-1",
		}).Return(expected, nil)

		result, _, err := handleListTasks(context.Background(), listTasksInput{
			TS: "ts1", Status: &status, Assignee: &assignee,
		}, mockPlanMgr)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Content, 1)
		content, ok := result.Content[0].(*gomcp.TextContent)
		require.True(t, ok)
		assert.Contains(t, content.Text, `"t1"`)
	})

	t.Run("no filters", func(t *testing.T) {
		mockPlanMgr := mockcontroller.NewMockPlanManager(t)
		mockPlanMgr.EXPECT().ListTasks(mock.Anything, "ts1", adapter.TicketFilter{}).Return(nil, nil)

		result, _, err := handleListTasks(context.Background(), listTasksInput{TS: "ts1"}, mockPlanMgr)
		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("status only", func(t *testing.T) {
		mockPlanMgr := mockcontroller.NewMockPlanManager(t)
		status := "done"
		mockPlanMgr.EXPECT().ListTasks(mock.Anything, "ts1", adapter.TicketFilter{
			Status: "done",
		}).Return([]tc.Ticket{}, nil)

		result, _, err := handleListTasks(context.Background(), listTasksInput{
			TS: "ts1", Status: &status,
		}, mockPlanMgr)
		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("assignee only", func(t *testing.T) {
		mockPlanMgr := mockcontroller.NewMockPlanManager(t)
		assignee := "agent-2"
		mockPlanMgr.EXPECT().ListTasks(mock.Anything, "ts1", adapter.TicketFilter{
			Assignee: "agent-2",
		}).Return([]tc.Ticket{}, nil)

		result, _, err := handleListTasks(context.Background(), listTasksInput{
			TS: "ts1", Assignee: &assignee,
		}, mockPlanMgr)
		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("error", func(t *testing.T) {
		mockPlanMgr := mockcontroller.NewMockPlanManager(t)
		mockPlanMgr.EXPECT().ListTasks(mock.Anything, "ts1", adapter.TicketFilter{}).Return(nil, errors.New("fail"))

		_, _, err := handleListTasks(context.Background(), listTasksInput{TS: "ts1"}, mockPlanMgr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "fail")
	})
}

func TestHandleGetTask(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := tc.Ticket{Id: strPtr("t1"), Title: strPtr("Task 1")}
	mockPlanMgr.EXPECT().GetTask(mock.Anything, "ts1", "t1").Return(expected, nil)

	result, _, err := handleGetTask(context.Background(), getTaskInput{TS: "ts1", ID: "t1"}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"t1"`)
}

func TestHandleGetTask_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().GetTask(mock.Anything, "ts1", "t1").Return(tc.Ticket{}, errors.New("not found"))

	_, _, err := handleGetTask(context.Background(), getTaskInput{TS: "ts1", ID: "t1"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandleAssignTask(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := tc.Ticket{Id: strPtr("t1"), Assignee: strPtr("agent-1")}
	mockPlanMgr.EXPECT().AssignTask(mock.Anything, "ts1", "t1", "agent-1").Return(expected, nil)

	result, _, err := handleAssignTask(context.Background(), assignTaskInput{
		TS: "ts1", TicketID: "t1", AgentID: "agent-1",
	}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"agent-1"`)
}

func TestHandleAssignTask_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().AssignTask(mock.Anything, "ts1", "t1", "agent-1").Return(tc.Ticket{}, errors.New("conflict"))

	_, _, err := handleAssignTask(context.Background(), assignTaskInput{
		TS: "ts1", TicketID: "t1", AgentID: "agent-1",
	}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "conflict")
}

func TestHandleCompleteTask(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := tc.Ticket{Id: strPtr("t1"), Status: strPtr("done")}
	mockPlanMgr.EXPECT().CompleteTask(mock.Anything, "ts1", "t1").Return(expected, nil)

	result, _, err := handleCompleteTask(context.Background(), completeTaskInput{
		TS: "ts1", TicketID: "t1",
	}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"done"`)
}

func TestHandleCompleteTask_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().CompleteTask(mock.Anything, "ts1", "t1").Return(tc.Ticket{}, errors.New("already completed"))

	_, _, err := handleCompleteTask(context.Background(), completeTaskInput{
		TS: "ts1", TicketID: "t1",
	}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already completed")
}

func TestHandleCreatePlan(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	status := "active"
	expected := tc.Plan{Id: strPtr("p1"), Title: strPtr("Plan 1")}
	mockPlanMgr.EXPECT().CreatePlan(mock.Anything, "ts1", tc.CreatePlanRequest{
		Id:     "p1",
		Title:  "Plan 1",
		Status: &status,
	}).Return(expected, nil)

	result, _, err := handleCreatePlan(context.Background(), createPlanInput{
		TS:     "ts1",
		ID:     "p1",
		Title:  "Plan 1",
		Status: &status,
	}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"p1"`)
}

func TestHandleCreatePlan_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().CreatePlan(mock.Anything, mock.Anything, mock.Anything).
		Return(tc.Plan{}, errors.New("conflict"))

	_, _, err := handleCreatePlan(context.Background(), createPlanInput{TS: "ts1", ID: "p1", Title: "Plan"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "conflict")
}

func TestHandleCreateTask(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := tc.Ticket{Id: strPtr("t1"), Title: strPtr("Task 1")}
	mockPlanMgr.EXPECT().CreateTask(mock.Anything, "ts1", tc.CreateTicketRequest{
		Id:    "t1",
		Title: "Task 1",
	}).Return(expected, nil)

	result, _, err := handleCreateTask(context.Background(), createTaskInput{
		TS:    "ts1",
		ID:    "t1",
		Title: "Task 1",
	}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"t1"`)
}

func TestHandleCreateTask_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().CreateTask(mock.Anything, mock.Anything, mock.Anything).
		Return(tc.Ticket{}, errors.New("conflict"))

	_, _, err := handleCreateTask(context.Background(), createTaskInput{TS: "ts1", ID: "t1", Title: "Task"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "conflict")
}

func TestHandleCreateMetaPlan(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := tc.MetaPlan{Id: strPtr("mp-1")}
	mockPlanMgr.EXPECT().CreateMetaPlan(mock.Anything, "ts1", tc.CreateMetaPlanRequest{
		Id:    "mp-1",
		Title: "Meta Plan 1",
	}).Return(expected, nil)

	result, _, err := handleCreateMetaPlan(context.Background(), createMetaPlanInput{
		TS:    "ts1",
		ID:    "mp-1",
		Title: "Meta Plan 1",
	}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"mp-1"`)
}

func TestHandleCreateMetaPlan_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().CreateMetaPlan(mock.Anything, mock.Anything, mock.Anything).
		Return(tc.MetaPlan{}, errors.New("conflict"))

	_, _, err := handleCreateMetaPlan(context.Background(), createMetaPlanInput{TS: "ts1", ID: "mp-1", Title: "MP"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "conflict")
}

func TestHandleUpdatePlan(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	status := "active"
	expected := tc.Plan{Id: strPtr("p1"), Title: strPtr("Updated Plan")}
	mockPlanMgr.EXPECT().UpdatePlan(mock.Anything, "ts1", "p1", tc.UpdatePlanRequest{
		Title:  "Updated Plan",
		Status: &status,
	}).Return(expected, nil)

	result, _, err := handleUpdatePlan(context.Background(), updatePlanInput{
		TS:     "ts1",
		ID:     "p1",
		Title:  "Updated Plan",
		Status: &status,
	}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"Updated Plan"`)
}

func TestHandleUpdatePlan_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().UpdatePlan(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(tc.Plan{}, errors.New("not found"))

	_, _, err := handleUpdatePlan(context.Background(), updatePlanInput{TS: "ts1", ID: "p1", Title: "T"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandleUpdateMetaPlan(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := tc.MetaPlan{Id: strPtr("mp-1"), Title: strPtr("Updated MP")}
	mockPlanMgr.EXPECT().UpdateMetaPlan(mock.Anything, "ts1", "mp-1", tc.UpdateMetaPlanRequest{
		Title: "Updated MP",
	}).Return(expected, nil)

	result, _, err := handleUpdateMetaPlan(context.Background(), updateMetaPlanInput{
		TS:    "ts1",
		ID:    "mp-1",
		Title: "Updated MP",
	}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"Updated MP"`)
}

func TestHandleUpdateMetaPlan_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().UpdateMetaPlan(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(tc.MetaPlan{}, errors.New("not found"))

	_, _, err := handleUpdateMetaPlan(context.Background(), updateMetaPlanInput{TS: "ts1", ID: "mp-1", Title: "T"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandleUpdateTask(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := tc.Ticket{Id: strPtr("t1"), Title: strPtr("Updated Task")}
	mockPlanMgr.EXPECT().UpdateTask(mock.Anything, "ts1", "t1", tc.UpdateTicketRequest{
		Title: "Updated Task",
	}).Return(expected, nil)

	result, _, err := handleUpdateTask(context.Background(), updateTaskInput{
		TS:    "ts1",
		ID:    "t1",
		Title: "Updated Task",
	}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"Updated Task"`)
}

func TestHandleUpdateTask_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().UpdateTask(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(tc.Ticket{}, errors.New("not found"))

	_, _, err := handleUpdateTask(context.Background(), updateTaskInput{TS: "ts1", ID: "t1", Title: "T"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandleDeletePlan(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().DeletePlan(mock.Anything, "ts1", "p1").Return(nil)

	result, _, err := handleDeletePlan(context.Background(), deletePlanInput{TS: "ts1", ID: "p1"}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Equal(t, "deleted plan p1", content.Text)
}

func TestHandleDeletePlan_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().DeletePlan(mock.Anything, "ts1", "p1").Return(errors.New("not found"))

	_, _, err := handleDeletePlan(context.Background(), deletePlanInput{TS: "ts1", ID: "p1"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandleDeleteTask(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().DeleteTask(mock.Anything, "ts1", "t1").Return(nil)

	result, _, err := handleDeleteTask(context.Background(), deleteTaskInput{TS: "ts1", ID: "t1"}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Equal(t, "deleted task t1", content.Text)
}

func TestHandleDeleteTask_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().DeleteTask(mock.Anything, "ts1", "t1").Return(errors.New("not found"))

	_, _, err := handleDeleteTask(context.Background(), deleteTaskInput{TS: "ts1", ID: "t1"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandleDeleteMetaPlan(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().DeleteMetaPlan(mock.Anything, "ts1", "mp-1").Return(nil)

	result, _, err := handleDeleteMetaPlan(context.Background(), deleteMetaPlanInput{TS: "ts1", ID: "mp-1"}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Equal(t, "deleted meta-plan mp-1", content.Text)
}

func TestHandleDeleteMetaPlan_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().DeleteMetaPlan(mock.Anything, "ts1", "mp-1").Return(errors.New("not found"))

	_, _, err := handleDeleteMetaPlan(context.Background(), deleteMetaPlanInput{TS: "ts1", ID: "mp-1"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandleListChildren(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := []tc.Ticket{{Id: strPtr("child-1")}, {Id: strPtr("child-2")}}
	mockPlanMgr.EXPECT().ListChildren(mock.Anything, "ts1", "p1").Return(expected, nil)

	result, _, err := handleListChildren(context.Background(), listChildrenInput{TS: "ts1", ID: "p1"}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"child-1"`)
	assert.Contains(t, content.Text, `"child-2"`)
}

func TestHandleListChildren_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().ListChildren(mock.Anything, "ts1", "p1").Return(nil, errors.New("not found"))

	_, _, err := handleListChildren(context.Background(), listChildrenInput{TS: "ts1", ID: "p1"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandleListBlocking(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	expected := []tc.Ticket{{Id: strPtr("blocker-1")}}
	mockPlanMgr.EXPECT().ListBlocking(mock.Anything, "ts1", "t1").Return(expected, nil)

	result, _, err := handleListBlocking(context.Background(), listBlockingInput{TS: "ts1", ID: "t1"}, mockPlanMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"blocker-1"`)
}

func TestHandleListBlocking_Error(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockPlanMgr.EXPECT().ListBlocking(mock.Anything, "ts1", "t1").Return(nil, errors.New("not found"))

	_, _, err := handleListBlocking(context.Background(), listBlockingInput{TS: "ts1", ID: "t1"}, mockPlanMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// --- MemoryManager handler tests (2) ---

func TestHandleAddComment(t *testing.T) {
	mockMemMgr := mockcontroller.NewMockMemoryManager(t)
	expected := tc.Comment{Author: strPtr("agent-1"), Text: strPtr("hello")}
	mockMemMgr.EXPECT().AddComment(mock.Anything, "ts1", "t1", "agent-1", "hello", []string{"tag1"}).Return(expected, nil)

	result, _, err := handleAddComment(context.Background(), addCommentInput{
		TS:       "ts1",
		TicketID: "t1",
		Author:   "agent-1",
		Text:     "hello",
		Tags:     []string{"tag1"},
	}, mockMemMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"agent-1"`)
	assert.Contains(t, content.Text, `"hello"`)
}

func TestHandleAddComment_Error(t *testing.T) {
	mockMemMgr := mockcontroller.NewMockMemoryManager(t)
	mockMemMgr.EXPECT().AddComment(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(tc.Comment{}, errors.New("fail"))

	_, _, err := handleAddComment(context.Background(), addCommentInput{
		TS: "ts1", TicketID: "t1", Author: "a", Text: "t",
	}, mockMemMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail")
}

func TestHandleListMemories(t *testing.T) {
	mockMemMgr := mockcontroller.NewMockMemoryManager(t)
	expected := []types.Memory{{TicketID: "t1", Author: "agent-1", Text: "memory text"}}
	mockMemMgr.EXPECT().ListMemories(mock.Anything, "ts1", "t1", "agent-1").Return(expected, nil)

	result, _, err := handleListMemories(context.Background(), listMemoriesInput{
		TS: "ts1", TicketID: "t1", AgentID: "agent-1",
	}, mockMemMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"memory text"`)
}

func TestHandleListMemories_Error(t *testing.T) {
	mockMemMgr := mockcontroller.NewMockMemoryManager(t)
	mockMemMgr.EXPECT().ListMemories(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("fail"))

	_, _, err := handleListMemories(context.Background(), listMemoriesInput{
		TS: "ts1", TicketID: "t1", AgentID: "agent-1",
	}, mockMemMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail")
}

// --- TrackingSetManager handler tests (4) ---

func TestHandleCreateTrackingSet(t *testing.T) {
	mockTsMgr := mockcontroller.NewMockTrackingSetManager(t)
	expected := tc.TrackingSet{Name: strPtr("ts-new")}
	mockTsMgr.EXPECT().CreateTrackingSet(mock.Anything, "ts-new").Return(expected, nil)

	result, _, err := handleCreateTrackingSet(context.Background(), createTrackingSetInput{Name: "ts-new"}, mockTsMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"ts-new"`)
}

func TestHandleCreateTrackingSet_Error(t *testing.T) {
	mockTsMgr := mockcontroller.NewMockTrackingSetManager(t)
	mockTsMgr.EXPECT().CreateTrackingSet(mock.Anything, "ts-new").Return(tc.TrackingSet{}, errors.New("exists"))

	_, _, err := handleCreateTrackingSet(context.Background(), createTrackingSetInput{Name: "ts-new"}, mockTsMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exists")
}

func TestHandleListTrackingSets(t *testing.T) {
	mockTsMgr := mockcontroller.NewMockTrackingSetManager(t)
	expected := []tc.TrackingSet{{Name: strPtr("ts1")}, {Name: strPtr("ts2")}}
	mockTsMgr.EXPECT().ListTrackingSets(mock.Anything).Return(expected, nil)

	result, _, err := handleListTrackingSets(context.Background(), listTrackingSetsInput{}, mockTsMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"ts1"`)
	assert.Contains(t, content.Text, `"ts2"`)
}

func TestHandleListTrackingSets_Error(t *testing.T) {
	mockTsMgr := mockcontroller.NewMockTrackingSetManager(t)
	mockTsMgr.EXPECT().ListTrackingSets(mock.Anything).Return(nil, errors.New("fail"))

	_, _, err := handleListTrackingSets(context.Background(), listTrackingSetsInput{}, mockTsMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail")
}

func TestHandleGetTrackingSet(t *testing.T) {
	mockTsMgr := mockcontroller.NewMockTrackingSetManager(t)
	expected := tc.TrackingSet{Name: strPtr("ts1")}
	mockTsMgr.EXPECT().GetTrackingSet(mock.Anything, "ts1").Return(expected, nil)

	result, _, err := handleGetTrackingSet(context.Background(), getTrackingSetInput{TS: "ts1"}, mockTsMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"ts1"`)
}

func TestHandleGetTrackingSet_Error(t *testing.T) {
	mockTsMgr := mockcontroller.NewMockTrackingSetManager(t)
	mockTsMgr.EXPECT().GetTrackingSet(mock.Anything, "ts1").Return(tc.TrackingSet{}, errors.New("not found"))

	_, _, err := handleGetTrackingSet(context.Background(), getTrackingSetInput{TS: "ts1"}, mockTsMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandleDeleteTrackingSet(t *testing.T) {
	mockTsMgr := mockcontroller.NewMockTrackingSetManager(t)
	mockTsMgr.EXPECT().DeleteTrackingSet(mock.Anything, "ts1").Return(nil)

	result, _, err := handleDeleteTrackingSet(context.Background(), deleteTrackingSetInput{TS: "ts1"}, mockTsMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Equal(t, "deleted tracking set ts1", content.Text)
}

func TestHandleDeleteTrackingSet_Error(t *testing.T) {
	mockTsMgr := mockcontroller.NewMockTrackingSetManager(t)
	mockTsMgr.EXPECT().DeleteTrackingSet(mock.Anything, "ts1").Return(errors.New("not found"))

	_, _, err := handleDeleteTrackingSet(context.Background(), deleteTrackingSetInput{TS: "ts1"}, mockTsMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// --- EdgeManager handler tests (3) ---

func TestHandleCreateEdge(t *testing.T) {
	mockEdgeMgr := mockcontroller.NewMockEdgeManager(t)
	expected := tc.Edge{From: strPtr("t1"), To: strPtr("t2"), Type: strPtr("blocks")}
	mockEdgeMgr.EXPECT().AddEdge(mock.Anything, "ts1", tc.EdgeRequest{
		From: "t1", To: "t2", Type: "blocks",
	}).Return(expected, nil)

	result, _, err := handleCreateEdge(context.Background(), createEdgeInput{
		TS: "ts1", From: "t1", To: "t2", Type: "blocks",
	}, mockEdgeMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, content.Text, `"t1"`)
	assert.Contains(t, content.Text, `"t2"`)
	assert.Contains(t, content.Text, `"blocks"`)
}

func TestHandleCreateEdge_Error(t *testing.T) {
	mockEdgeMgr := mockcontroller.NewMockEdgeManager(t)
	mockEdgeMgr.EXPECT().AddEdge(mock.Anything, mock.Anything, mock.Anything).
		Return(tc.Edge{}, errors.New("conflict"))

	_, _, err := handleCreateEdge(context.Background(), createEdgeInput{
		TS: "ts1", From: "t1", To: "t2", Type: "blocks",
	}, mockEdgeMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "conflict")
}

func TestHandleListEdges(t *testing.T) {
	t.Run("with filters", func(t *testing.T) {
		mockEdgeMgr := mockcontroller.NewMockEdgeManager(t)
		ticket := "t1"
		edgeType := "blocks"
		expected := []tc.Edge{{From: strPtr("t1"), To: strPtr("t2"), Type: strPtr("blocks")}}
		mockEdgeMgr.EXPECT().ListEdges(mock.Anything, "ts1", &tc.ListEdgesParams{
			Ticket: &ticket,
			Type:   &edgeType,
		}).Return(expected, nil)

		result, _, err := handleListEdges(context.Background(), listEdgesInput{
			TS: "ts1", Ticket: &ticket, Type: &edgeType,
		}, mockEdgeMgr)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Content, 1)
		content, ok := result.Content[0].(*gomcp.TextContent)
		require.True(t, ok)
		assert.Contains(t, content.Text, `"blocks"`)
	})

	t.Run("no filters", func(t *testing.T) {
		mockEdgeMgr := mockcontroller.NewMockEdgeManager(t)
		mockEdgeMgr.EXPECT().ListEdges(mock.Anything, "ts1", &tc.ListEdgesParams{}).Return([]tc.Edge{}, nil)

		result, _, err := handleListEdges(context.Background(), listEdgesInput{TS: "ts1"}, mockEdgeMgr)
		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("error", func(t *testing.T) {
		mockEdgeMgr := mockcontroller.NewMockEdgeManager(t)
		mockEdgeMgr.EXPECT().ListEdges(mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("fail"))

		_, _, err := handleListEdges(context.Background(), listEdgesInput{TS: "ts1"}, mockEdgeMgr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "fail")
	})
}

func TestHandleDeleteEdge(t *testing.T) {
	mockEdgeMgr := mockcontroller.NewMockEdgeManager(t)
	mockEdgeMgr.EXPECT().RemoveEdge(mock.Anything, "ts1", tc.EdgeRequest{
		From: "t1", To: "t2", Type: "blocks",
	}).Return(nil)

	result, _, err := handleDeleteEdge(context.Background(), deleteEdgeInput{
		TS: "ts1", From: "t1", To: "t2", Type: "blocks",
	}, mockEdgeMgr)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(*gomcp.TextContent)
	require.True(t, ok)
	assert.Equal(t, "deleted edge t1 -> t2", content.Text)
}

func TestHandleDeleteEdge_Error(t *testing.T) {
	mockEdgeMgr := mockcontroller.NewMockEdgeManager(t)
	mockEdgeMgr.EXPECT().RemoveEdge(mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("not found"))

	_, _, err := handleDeleteEdge(context.Background(), deleteEdgeInput{
		TS: "ts1", From: "t1", To: "t2", Type: "blocks",
	}, mockEdgeMgr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
