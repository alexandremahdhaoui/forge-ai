//go:build integration

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
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"

	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
)

func TestGeneratedClientMatchesSpec(t *testing.T) {
	specPath := "../../api/forge-tracker.v1.yaml"
	data, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("reading OpenAPI spec: %v", err)
	}

	loader := openapi3.NewLoader()
	spec, err := loader.LoadFromData(data)
	if err != nil {
		t.Fatalf("parsing OpenAPI spec: %v", err)
	}

	// Extract operation IDs from spec.
	var operationIDs []string
	for _, pathItem := range spec.Paths.Map() {
		for _, op := range pathItem.Operations() {
			if op.OperationID != "" {
				operationIDs = append(operationIDs, op.OperationID)
			}
		}
	}

	if len(operationIDs) == 0 {
		t.Fatal("no operation IDs found in spec")
	}

	// Verify each operation ID has a corresponding method on ClientWithResponses.
	// oapi-codegen generates methods named "<PascalCaseOperationID>WithResponse".
	clientType := reflect.TypeOf(&tc.ClientWithResponses{})
	for _, opID := range operationIDs {
		methodName := strings.ToUpper(opID[:1]) + opID[1:] + "WithResponse"
		_, found := clientType.MethodByName(methodName)
		if !found {
			t.Errorf("operation %q: expected method %q on ClientWithResponses, not found", opID, methodName)
		}
	}

	t.Logf("validated %d operations against generated client", len(operationIDs))
}
