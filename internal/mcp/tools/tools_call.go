package tools

import (
	"context"
	"encoding/json"

	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/handlers"
	"github.com/blacksheepkhan/fileserver-mcp/internal/protocol"
)

const toolsCallMethod = "tools/call"

// CallHandler handles MCP tools/call requests.
type CallHandler struct {
	registry *Registry
}

// NewCallHandler creates a new tools/call handler.
func NewCallHandler(registry *Registry) *CallHandler {
	return &CallHandler{
		registry: registry,
	}
}

// Method returns the JSON-RPC method handled by this handler.
func (h *CallHandler) Method() string {
	return toolsCallMethod
}

// Handle executes the requested tool.
func (h *CallHandler) Handle(ctx handlers.Context, rawParams json.RawMessage) (any, *protocol.Error) {
	var params callParams
	if err := json.Unmarshal(rawParams, &params); err != nil {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "invalid params",
		}
	}

	if params.Name == "" {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "missing tool name",
		}
	}

	tool, ok := h.registry.Get(params.Name)
	if !ok {
		return nil, &protocol.Error{
			Code:    protocol.ErrMethodNotFound,
			Message: "tool not found",
		}
	}

	execCtx := ctx.Context
	if execCtx == nil {
		execCtx = context.Background()
	}

	return tool.Execute(execCtx, params.Arguments)
}

type callParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
}
