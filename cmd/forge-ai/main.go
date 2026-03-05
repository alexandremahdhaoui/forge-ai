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

package main

import (
	"log"
	"os"

	"github.com/alexandremahdhaoui/forge-ai/internal/adapter"
	"github.com/alexandremahdhaoui/forge-ai/internal/controller"
	mcpdriver "github.com/alexandremahdhaoui/forge-ai/internal/driver/mcp"
	"github.com/alexandremahdhaoui/forge/pkg/enginecli"
	"github.com/alexandremahdhaoui/forge/pkg/mcpserver"
)

// Version is set via ldflags at build time.
var Version = "dev"

func main() {
	enginecli.Bootstrap(enginecli.Config{
		Name:    "forge-ai",
		Version: Version,
		RunMCP:  runMCPServer(),
		RunCLI:  nil,
	})
}

func runMCPServer() func() error {
	return func() error {
		trackerURL := os.Getenv("FORGE_TRACKER_URL")
		if trackerURL == "" {
			trackerURL = "http://localhost:8080"
		}

		client, err := adapter.NewHTTPTrackerClient(trackerURL)
		if err != nil {
			log.Fatalf("creating tracker client: %v", err)
		}

		planMgr := controller.NewPlanManager(client)
		memMgr := controller.NewMemoryManager(client)

		server := mcpserver.New("forge-ai", Version)
		mcpdriver.RegisterTools(server, planMgr, memMgr)

		return server.RunDefault()
	}
}
