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

var _ TrackingSetManager = (*trackingSetManager)(nil)

type trackingSetManager struct {
	client adapter.TrackerClient
}

// NewTrackingSetManager creates a TrackingSetManager with the given tracker client.
func NewTrackingSetManager(client adapter.TrackerClient) TrackingSetManager {
	return &trackingSetManager{client: client}
}

func (m *trackingSetManager) CreateTrackingSet(ctx context.Context, name string) (tc.TrackingSet, error) {
	return m.client.CreateTrackingSet(ctx, tc.CreateTrackingSetRequest{Name: name})
}

func (m *trackingSetManager) ListTrackingSets(ctx context.Context) ([]tc.TrackingSet, error) {
	return m.client.ListTrackingSets(ctx)
}

func (m *trackingSetManager) GetTrackingSet(ctx context.Context, ts string) (tc.TrackingSet, error) {
	return m.client.GetTrackingSet(ctx, ts)
}

func (m *trackingSetManager) DeleteTrackingSet(ctx context.Context, ts string) error {
	return m.client.DeleteTrackingSet(ctx, ts)
}
