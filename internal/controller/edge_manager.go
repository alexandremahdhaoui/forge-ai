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
	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
)

// EdgeManager handles edge (relationship) operations between tickets.
type EdgeManager interface {
	ListEdges(ctx context.Context, ts string, params *tc.ListEdgesParams) ([]tc.Edge, error)
	AddEdge(ctx context.Context, ts string, req tc.EdgeRequest) (tc.Edge, error)
	RemoveEdge(ctx context.Context, ts string, req tc.EdgeRequest) error
}

var _ EdgeManager = (*edgeManager)(nil)

type edgeManager struct {
	client adapter.TrackerClient
}

// NewEdgeManager creates an EdgeManager with the given tracker client.
func NewEdgeManager(client adapter.TrackerClient) EdgeManager {
	return &edgeManager{client: client}
}

func (m *edgeManager) ListEdges(ctx context.Context, ts string, params *tc.ListEdgesParams) ([]tc.Edge, error) {
	return m.client.ListEdges(ctx, ts, params)
}

func (m *edgeManager) AddEdge(ctx context.Context, ts string, req tc.EdgeRequest) (tc.Edge, error) {
	return m.client.AddEdge(ctx, ts, req)
}

func (m *edgeManager) RemoveEdge(ctx context.Context, ts string, req tc.EdgeRequest) error {
	return m.client.RemoveEdge(ctx, ts, req)
}
