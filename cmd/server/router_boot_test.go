package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/handlers"
	"github.com/blacksheepkhan/fileserver-mcp/internal/protocol"
)

func TestCreateRouterRegistersInitialize(t *testing.T) {
	registry := createToolRegistry(noopFileSystem{}, 1024, toolCapabilities{filesystemWrite: true})
	mcpRouter := createRouter("test-server", "test-version", registry)

	params := json.RawMessage(`{
        "protocolVersion": "2025-11-25",
        "capabilities": {},
        "clientInfo": {
            "name": "test-client",
            "version": "test-version"
        }
    }`)

	result, protocolErr := mcpRouter.Dispatch(
		"initialize",
		handlers.Context{Context: context.Background()},
		params,
	)

	if protocolErr != nil {
		t.Fatalf("expected initialize to be registered, got protocol error: %+v", protocolErr)
	}

	if result == nil {
		t.Fatal("expected initialize result")
	}
}

func TestCreateRouterRegistersToolsList(t *testing.T) {
	registry := createToolRegistry(noopFileSystem{}, 1024, toolCapabilities{filesystemWrite: true})
	mcpRouter := createRouter("test-server", "test-version", registry)

	result, protocolErr := mcpRouter.Dispatch(
		"tools/list",
		handlers.Context{Context: context.Background()},
		json.RawMessage(`{}`),
	)

	if protocolErr != nil {
		t.Fatalf("expected tools/list to be registered, got protocol error: %+v", protocolErr)
	}

	if result == nil {
		t.Fatal("expected tools/list result")
	}
}

func TestCreateRouterRegistersToolsCall(t *testing.T) {
	registry := createToolRegistry(noopFileSystem{}, 1024, toolCapabilities{filesystemWrite: true})
	mcpRouter := createRouter("test-server", "test-version", registry)

	_, protocolErr := mcpRouter.Dispatch(
		"tools/call",
		handlers.Context{Context: context.Background()},
		json.RawMessage(`{}`),
	)

	if protocolErr == nil {
		t.Fatal("expected protocol error for invalid tools/call params")
	}

	if protocolErr.Code == protocol.ErrMethodNotFound {
		t.Fatalf("expected tools/call to be registered, got method not found: %+v", protocolErr)
	}

	if protocolErr.Code != protocol.ErrInvalidParams {
		t.Fatalf("expected invalid params from registered tools/call handler, got: %+v", protocolErr)
	}
}

func TestCreateRouterRejectsWriteToolCallWhenReadOnly(t *testing.T) {
	registry := createToolRegistry(noopFileSystem{}, 1024, capabilitiesFromReadOnly(true))
	mcpRouter := createRouter("test-server", "test-version", registry)

	_, protocolErr := mcpRouter.Dispatch(
		"tools/call",
		handlers.Context{Context: context.Background()},
		json.RawMessage(`{"name":"write_file","arguments":{"path":"out.txt","content":"blocked"}}`),
	)

	if protocolErr == nil {
		t.Fatal("expected protocol error for disabled write tool")
	}

	if protocolErr.Code != protocol.ErrMethodNotFound {
		t.Fatalf("expected method not found for disabled write tool, got: %+v", protocolErr)
	}
}

func TestCreateRouterRejectsUnknownMethod(t *testing.T) {
	registry := createToolRegistry(noopFileSystem{}, 1024, toolCapabilities{filesystemWrite: true})
	mcpRouter := createRouter("test-server", "test-version", registry)

	_, protocolErr := mcpRouter.Dispatch(
		"unknown/method",
		handlers.Context{Context: context.Background()},
		json.RawMessage(`{}`),
	)

	if protocolErr == nil {
		t.Fatal("expected protocol error")
	}

	if protocolErr.Code != protocol.ErrMethodNotFound {
		t.Fatalf("expected method not found, got: %+v", protocolErr)
	}
}
