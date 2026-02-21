package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/vacano-house/vacano-ui-mcp/internal/docs"
)

type SearchIconsParams struct {
	Query string `json:"query" jsonschema:"Search query to find icons by name, description, or category (e.g. 'arrow', 'close', 'navigation', 'chart')"`
}

func NewSearchIconsHandler(store *docs.Store) func(context.Context, *mcp.CallToolRequest, *SearchIconsParams) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, params *SearchIconsParams) (*mcp.CallToolResult, any, error) {
		if params.Query == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "query parameter is required"}},
				IsError: true,
			}, nil, nil
		}

		results := store.SearchIcons(params.Query)

		if len(results) == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No icons found for: %s", params.Query)}},
			}, nil, nil
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Found %d icon(s) for \"%s\".\n", len(results), params.Query))
		sb.WriteString("Import from `@vacano/ui/icons`.\n\n")

		currentCategory := ""
		for _, icon := range results {
			if icon.Category != currentCategory {
				sb.WriteString(fmt.Sprintf("### %s\n\n", icon.Category))
				currentCategory = icon.Category
			}
			sb.WriteString(fmt.Sprintf("- `%s` — %s\n", icon.Name, icon.Description))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: sb.String()}},
		}, nil, nil
	}
}
