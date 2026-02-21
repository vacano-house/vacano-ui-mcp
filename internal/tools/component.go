package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/vacano-house/vacano-ui-mcp/internal/docs"
)

type GetComponentParams struct {
	Name string `json:"name" jsonschema:"Exact component name (e.g. Button, Modal, DatePicker)"`
}

func NewGetComponentHandler(store *docs.Store) func(context.Context, *mcp.CallToolRequest, *GetComponentParams) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, params *GetComponentParams) (*mcp.CallToolResult, any, error) {
		if params.Name == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "name parameter is required"}},
				IsError: true,
			}, nil, nil
		}

		entry := store.GetByName(params.Name)

		if entry == nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Component not found: %s", params.Name)}},
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: entry.Content}},
		}, nil, nil
	}
}
