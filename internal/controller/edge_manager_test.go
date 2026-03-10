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

	"github.com/alexandremahdhaoui/forge-ai/internal/util/mocks/mockadapter"
	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListEdges(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	params := &tc.ListEdgesParams{Ticket: strPtr("t1")}
	expected := []tc.Edge{{From: strPtr("t1"), To: strPtr("t2"), Type: strPtr("blocks")}}
	mockClient.EXPECT().ListEdges(mock.Anything, "ts1", params).Return(expected, nil)

	mgr := NewEdgeManager(mockClient)
	result, err := mgr.ListEdges(context.Background(), "ts1", params)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestAddEdge(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	req := tc.EdgeRequest{From: "t1", To: "t2", Type: "blocks"}
	expected := tc.Edge{From: strPtr("t1"), To: strPtr("t2"), Type: strPtr("blocks")}
	mockClient.EXPECT().AddEdge(mock.Anything, "ts1", req).Return(expected, nil)

	mgr := NewEdgeManager(mockClient)
	result, err := mgr.AddEdge(context.Background(), "ts1", req)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestRemoveEdge(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	req := tc.EdgeRequest{From: "t1", To: "t2", Type: "blocks"}
	mockClient.EXPECT().RemoveEdge(mock.Anything, "ts1", req).Return(nil)

	mgr := NewEdgeManager(mockClient)
	err := mgr.RemoveEdge(context.Background(), "ts1", req)
	require.NoError(t, err)
}

func TestAddEdge_Error(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	req := tc.EdgeRequest{From: "t1", To: "t2", Type: "blocks"}
	mockClient.EXPECT().AddEdge(mock.Anything, "ts1", req).
		Return(tc.Edge{}, errors.New("internal server error"))

	mgr := NewEdgeManager(mockClient)
	_, err := mgr.AddEdge(context.Background(), "ts1", req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "internal server error")
}
