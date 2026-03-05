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
	"fmt"
	"strings"

	"github.com/alexandremahdhaoui/forge-ai/internal/adapter"
	"github.com/alexandremahdhaoui/forge-ai/internal/types"
	tc "github.com/alexandremahdhaoui/forge-ai/pkg/generated/trackerclient"
)

// MemoryManager handles comment/memory operations.
type MemoryManager interface {
	AddComment(ctx context.Context, ts, ticketID, author, text string, tags []string) (tc.Comment, error)
	ListMemories(ctx context.Context, ts, ticketID, agentID string) ([]types.Memory, error)
}

var _ MemoryManager = (*memoryManager)(nil)

type memoryManager struct {
	client adapter.TrackerClient
}

// NewMemoryManager creates a MemoryManager with the given tracker client.
func NewMemoryManager(client adapter.TrackerClient) MemoryManager {
	return &memoryManager{client: client}
}

func (m *memoryManager) AddComment(ctx context.Context, ts, ticketID, author, text string, tags []string) (tc.Comment, error) {
	formattedText := text
	if len(tags) > 0 {
		formattedText = fmt.Sprintf("[%s] %s", strings.Join(tags, ","), text)
	}
	return m.client.AddComment(ctx, ts, ticketID, tc.AddCommentRequest{
		Author: author,
		Text:   formattedText,
	})
}

func (m *memoryManager) ListMemories(ctx context.Context, ts, ticketID, agentID string) ([]types.Memory, error) {
	ticket, err := m.client.GetTicket(ctx, ts, ticketID)
	if err != nil {
		return nil, fmt.Errorf("getting ticket for memories: %w", err)
	}
	var memories []types.Memory
	if ticket.Comments == nil {
		return memories, nil
	}
	for _, c := range *ticket.Comments {
		author := ""
		if c.Author != nil {
			author = *c.Author
		}
		if agentID != "" && author != agentID {
			continue
		}
		mem := types.Memory{
			TicketID: ticketID,
			Author:   author,
		}
		if c.Timestamp != nil {
			mem.Timestamp = *c.Timestamp
		}
		if c.Text != nil {
			mem.Text = *c.Text
			mem.Tags = parseTags(*c.Text)
		}
		memories = append(memories, mem)
	}
	return memories, nil
}

// parseTags extracts tags from formatted text like "[blocker,progress] actual text".
func parseTags(text string) []string {
	if !strings.HasPrefix(text, "[") {
		return nil
	}
	end := strings.Index(text, "]")
	if end < 0 {
		return nil
	}
	tagStr := text[1:end]
	return strings.Split(tagStr, ",")
}
