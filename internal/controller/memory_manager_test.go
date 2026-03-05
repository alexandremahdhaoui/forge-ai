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
	"time"

	"github.com/alexandremahdhaoui/forge-ai/internal/util/mocks/mockadapter"
	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddComment(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	now := time.Now()
	mockClient.EXPECT().AddComment(mock.Anything, "ts1", "task-1", mock.MatchedBy(func(req tc.AddCommentRequest) bool {
		return req.Author == "agent-1" && req.Text == "[blocker,progress] Found issue"
	})).Return(tc.Comment{
		Author:    strPtr("agent-1"),
		Text:      strPtr("[blocker,progress] Found issue"),
		Timestamp: &now,
	}, nil)

	mgr := NewMemoryManager(mockClient)
	result, err := mgr.AddComment(context.Background(), "ts1", "task-1", "agent-1", "Found issue", []string{"blocker", "progress"})
	require.NoError(t, err)
	assert.Equal(t, "agent-1", *result.Author)
	assert.Equal(t, "[blocker,progress] Found issue", *result.Text)
}

func TestAddComment_NoTags(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	now := time.Now()
	mockClient.EXPECT().AddComment(mock.Anything, "ts1", "task-1", mock.MatchedBy(func(req tc.AddCommentRequest) bool {
		return req.Author == "agent-1" && req.Text == "Just a plain comment"
	})).Return(tc.Comment{
		Author:    strPtr("agent-1"),
		Text:      strPtr("Just a plain comment"),
		Timestamp: &now,
	}, nil)

	mgr := NewMemoryManager(mockClient)
	result, err := mgr.AddComment(context.Background(), "ts1", "task-1", "agent-1", "Just a plain comment", nil)
	require.NoError(t, err)
	assert.Equal(t, "Just a plain comment", *result.Text)
}

func TestListMemories(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	now := time.Now()
	comments := []tc.Comment{
		{Author: strPtr("agent-1"), Text: strPtr("[blocker] Issue found"), Timestamp: &now},
		{Author: strPtr("agent-2"), Text: strPtr("Progress update"), Timestamp: &now},
	}
	mockClient.EXPECT().GetTicket(mock.Anything, "ts1", "task-1").
		Return(tc.Ticket{Id: strPtr("task-1"), Comments: &comments}, nil)

	mgr := NewMemoryManager(mockClient)
	memories, err := mgr.ListMemories(context.Background(), "ts1", "task-1", "")
	require.NoError(t, err)
	assert.Len(t, memories, 2)
	assert.Equal(t, "agent-1", memories[0].Author)
	assert.Equal(t, []string{"blocker"}, memories[0].Tags)
	assert.Equal(t, "agent-2", memories[1].Author)
	assert.Nil(t, memories[1].Tags)
}

func TestListMemories_FilterByAgent(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	now := time.Now()
	comments := []tc.Comment{
		{Author: strPtr("agent-1"), Text: strPtr("[blocker] Issue found"), Timestamp: &now},
		{Author: strPtr("agent-2"), Text: strPtr("Progress update"), Timestamp: &now},
		{Author: strPtr("agent-1"), Text: strPtr("Another note"), Timestamp: &now},
	}
	mockClient.EXPECT().GetTicket(mock.Anything, "ts1", "task-1").
		Return(tc.Ticket{Id: strPtr("task-1"), Comments: &comments}, nil)

	mgr := NewMemoryManager(mockClient)
	memories, err := mgr.ListMemories(context.Background(), "ts1", "task-1", "agent-1")
	require.NoError(t, err)
	assert.Len(t, memories, 2)
	for _, mem := range memories {
		assert.Equal(t, "agent-1", mem.Author)
	}
}

func TestListMemories_NoComments(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	mockClient.EXPECT().GetTicket(mock.Anything, "ts1", "task-1").
		Return(tc.Ticket{Id: strPtr("task-1"), Comments: nil}, nil)

	mgr := NewMemoryManager(mockClient)
	memories, err := mgr.ListMemories(context.Background(), "ts1", "task-1", "")
	require.NoError(t, err)
	assert.Empty(t, memories)
}

func TestAddComment_Error(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	mockClient.EXPECT().AddComment(mock.Anything, "ts1", "task-1", mock.Anything).
		Return(tc.Comment{}, errors.New("server error"))

	mgr := NewMemoryManager(mockClient)
	_, err := mgr.AddComment(context.Background(), "ts1", "task-1", "agent-1", "test", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "server error")
}
