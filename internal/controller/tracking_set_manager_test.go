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

func TestCreateTrackingSet(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	expected := tc.TrackingSet{Name: strPtr("ts1")}
	mockClient.EXPECT().CreateTrackingSet(mock.Anything, tc.CreateTrackingSetRequest{Name: "ts1"}).Return(expected, nil)

	mgr := NewTrackingSetManager(mockClient)
	result, err := mgr.CreateTrackingSet(context.Background(), "ts1")
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestListTrackingSets(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	expected := []tc.TrackingSet{{Name: strPtr("ts1")}, {Name: strPtr("ts2")}}
	mockClient.EXPECT().ListTrackingSets(mock.Anything).Return(expected, nil)

	mgr := NewTrackingSetManager(mockClient)
	result, err := mgr.ListTrackingSets(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetTrackingSet(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	expected := tc.TrackingSet{Name: strPtr("ts1")}
	mockClient.EXPECT().GetTrackingSet(mock.Anything, "ts1").Return(expected, nil)

	mgr := NewTrackingSetManager(mockClient)
	result, err := mgr.GetTrackingSet(context.Background(), "ts1")
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestDeleteTrackingSet(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	mockClient.EXPECT().DeleteTrackingSet(mock.Anything, "ts1").Return(nil)

	mgr := NewTrackingSetManager(mockClient)
	err := mgr.DeleteTrackingSet(context.Background(), "ts1")
	require.NoError(t, err)
}

func TestCreateTrackingSet_Error(t *testing.T) {
	mockClient := mockadapter.NewMockTrackerClient(t)
	mockClient.EXPECT().CreateTrackingSet(mock.Anything, tc.CreateTrackingSetRequest{Name: "ts1"}).
		Return(tc.TrackingSet{}, errors.New("connection refused"))

	mgr := NewTrackingSetManager(mockClient)
	_, err := mgr.CreateTrackingSet(context.Background(), "ts1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
}
