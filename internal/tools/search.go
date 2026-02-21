package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/vacano-house/vacano-ui-mcp/internal/docs"
)

type SearchParams struct {
	Query string `json:"query" jsonschema:"Search query to find in component names, descriptions, and documentation content"`
}

func NewSearchHandler(store *docs.Store) func(context.Context, *mcp.CallToolRequest, *SearchParams) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, params *SearchParams) (*mcp.CallToolResult, any, error) {
		if params.Query == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "query parameter is required"}},
				IsError: true,
			}, nil, nil
		}

		results := store.Search(params.Query)

		if len(results) == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No results found for: %s", params.Query)}},
			}, nil, nil
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Found %d result(s) for \"%s\":\n\n", len(results), params.Query))

		for _, entry := range results {
			sb.WriteString(fmt.Sprintf("## %s [%s]\n", entry.Name, entry.Category))
			sb.WriteString(entry.Description)
			sb.WriteString("\n\n---\n\n")
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: sb.String()}},
		}, nil, nil
	}
}
