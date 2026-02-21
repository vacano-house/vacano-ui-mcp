package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/vacano-house/vacano-ui-mcp/internal/docs"
)

type ListParams struct {
	Category string `json:"category,omitempty" jsonschema:"Optional filter by category: form, data-display, feedback, layout, navigation, utility, guide"`
}

func NewListHandler(store *docs.Store) func(context.Context, *mcp.CallToolRequest, *ListParams) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, params *ListParams) (*mcp.CallToolResult, any, error) {
		results := store.List(params.Category)

		if len(results) == 0 {
			msg := "No components found"
			if params.Category != "" {
				msg = fmt.Sprintf("No components found in category: %s", params.Category)
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: msg}},
			}, nil, nil
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Found %d component(s):\n\n", len(results)))

		currentCategory := ""
		for _, entry := range results {
			cat := string(entry.Category)
			if cat != currentCategory {
				sb.WriteString(fmt.Sprintf("### %s\n\n", cat))
				currentCategory = cat
			}
			sb.WriteString(fmt.Sprintf("- **%s** — %s\n", entry.Name, entry.Description))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: sb.String()}},
		}, nil, nil
	}
}
