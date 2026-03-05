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
	"errors"
	"testing"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexandremahdhaoui/forge-ai/internal/util/mocks/mockcontroller"
	"github.com/alexandremahdhaoui/forge/pkg/mcpserver"
)

func TestRegisterTools(t *testing.T) {
	mockPlanMgr := mockcontroller.NewMockPlanManager(t)
	mockMemMgr := mockcontroller.NewMockMemoryManager(t)
	server := mcpserver.New("test-forge-ai", "test")

	require.NotPanics(t, func() {
		RegisterTools(server, mockPlanMgr, mockMemMgr)
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

		tc, ok := result.Content[0].(*gomcp.TextContent)
		require.True(t, ok, "expected *gomcp.TextContent, got %T", result.Content[0])
		assert.Contains(t, tc.Text, `"key"`)
		assert.Contains(t, tc.Text, `"value"`)
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

		tc, ok := result.Content[0].(*gomcp.TextContent)
		require.True(t, ok)
		assert.Contains(t, tc.Text, `"name":"test"`)
		assert.Contains(t, tc.Text, `"age":42`)
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

		tc, ok := result.Content[0].(*gomcp.TextContent)
		require.True(t, ok)
		assert.Equal(t, "null", tc.Text)
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
