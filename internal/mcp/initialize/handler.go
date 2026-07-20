package initialize

import (
	"encoding/json"

	"github.com/thomasweidner/flashgate-mcp/internal/mcp/handlers"
	"github.com/thomasweidner/flashgate-mcp/internal/protocol"
)

const method = "initialize"

// Handler handles MCP initialize requests.
type Handler struct {
	serverName    string
	serverVersion string
}

// NewHandler creates a new initialize handler.
func NewHandler(serverName string, serverVersion string) *Handler {
	return &Handler{
		serverName:    serverName,
		serverVersion: serverVersion,
	}
}

// Method returns the JSON-RPC method handled by this handler.
func (h *Handler) Method() string {
	return method
}

// Handle processes an initialize request.
func (h *Handler) Handle(_ handlers.Context, rawParams json.RawMessage) (any, *protocol.Error) {
	var params requestParams
	if err := json.Unmarshal(rawParams, &params); err != nil {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "invalid initialize params",
		}
	}

	if params.ProtocolVersion == "" {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "missing protocol version",
		}
	}

	return response{
		ProtocolVersion: protocol.ProtocolVersion,
		Capabilities: serverCapabilities{
			Tools: toolsCapability{},
		},
		ServerInfo: implementation{
			Name:    h.serverName,
			Version: h.serverVersion,
		},
	}, nil
}

type requestParams struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    json.RawMessage `json:"capabilities,omitempty"`
	ClientInfo      implementation  `json:"clientInfo"`
}

type response struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    serverCapabilities `json:"capabilities"`
	ServerInfo      implementation     `json:"serverInfo"`
}

type serverCapabilities struct {
	Tools toolsCapability `json:"tools"`
}

type toolsCapability struct{}

type implementation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
