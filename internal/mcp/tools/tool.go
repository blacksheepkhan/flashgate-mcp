package tools

import (
	"context"
	"encoding/json"

	"github.com/blacksheepkhan/flashgate-mcp/internal/protocol"
)

// Tool defines an executable MCP tool.
type Tool interface {
	Name() string
	Description() string
	InputSchema() any
	Execute(ctx context.Context, arguments json.RawMessage) (any, *protocol.Error)
}
