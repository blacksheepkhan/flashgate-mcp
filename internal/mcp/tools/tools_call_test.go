package tools

import (
	"encoding/json"
	"testing"

	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/handlers"
	"github.com/blacksheepkhan/fileserver-mcp/internal/protocol"
)

func TestCallHandlerCallsRegisteredTool(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	tool := &testTool{
		name: "test_tool",
		result: map[string]any{
			"ok": true,
		},
	}

	registry.Register(tool)

	handler := NewCallHandler(registry)

	result, rpcErr := handler.Handle(
		handlers.Context{},
		json.RawMessage(`{"name":"test_tool","arguments":{"path":"."}}`),
	)

	if rpcErr != nil {
		t.Fatalf("expected no error, got %v", rpcErr)
	}

	if tool.called != 1 {
		t.Fatalf("expected tool to be called once, got %d", tool.called)
	}

	if string(tool.arguments) != `{"path":"."}` {
		t.Fatalf("expected arguments to be forwarded unchanged, got %s", string(tool.arguments))
	}

	resultMap, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected result map, got %#v", result)
	}

	if resultMap["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", resultMap["ok"])
	}
}

func TestCallHandlerReturnsMethodNotFoundForUnknownTool(t *testing.T) {
	t.Parallel()

	handler := NewCallHandler(NewRegistry())

	result, rpcErr := handler.Handle(
		handlers.Context{},
		json.RawMessage(`{"name":"missing_tool","arguments":{}}`),
	)

	if result != nil {
		t.Fatalf("expected nil result, got %#v", result)
	}

	if rpcErr == nil {
		t.Fatal("expected rpc error")
	}

	if rpcErr.Code != protocol.ErrMethodNotFound {
		t.Fatalf("expected ErrMethodNotFound, got %d", rpcErr.Code)
	}
}

func TestCallHandlerReturnsInvalidParamsForMalformedParams(t *testing.T) {
	t.Parallel()

	handler := NewCallHandler(NewRegistry())

	result, rpcErr := handler.Handle(
		handlers.Context{},
		json.RawMessage(`{"name":`),
	)

	if result != nil {
		t.Fatalf("expected nil result, got %#v", result)
	}

	if rpcErr == nil {
		t.Fatal("expected rpc error")
	}

	if rpcErr.Code != protocol.ErrInvalidParams {
		t.Fatalf("expected ErrInvalidParams, got %d", rpcErr.Code)
	}
}

func TestCallHandlerReturnsInvalidParamsForMissingName(t *testing.T) {
	t.Parallel()

	handler := NewCallHandler(NewRegistry())

	result, rpcErr := handler.Handle(
		handlers.Context{},
		json.RawMessage(`{"arguments":{}}`),
	)

	if result != nil {
		t.Fatalf("expected nil result, got %#v", result)
	}

	if rpcErr == nil {
		t.Fatal("expected rpc error")
	}

	if rpcErr.Code != protocol.ErrInvalidParams {
		t.Fatalf("expected ErrInvalidParams, got %d", rpcErr.Code)
	}
}

func TestCallHandlerReturnsToolError(t *testing.T) {
	t.Parallel()

	expectedErr := &protocol.Error{
		Code:    protocol.ErrInvalidParams,
		Message: "invalid path",
	}

	registry := NewRegistry()
	registry.Register(&testTool{
		name: "test_tool",
		err:  expectedErr,
	})

	handler := NewCallHandler(registry)

	result, rpcErr := handler.Handle(
		handlers.Context{},
		json.RawMessage(`{"name":"test_tool","arguments":{"path":"../outside"}}`),
	)

	if result != nil {
		t.Fatalf("expected nil result, got %#v", result)
	}

	if rpcErr != expectedErr {
		t.Fatal("expected tool error to be returned unchanged")
	}
}

func TestCallHandlerMethod(t *testing.T) {
	t.Parallel()

	handler := NewCallHandler(NewRegistry())

	if handler.Method() != "tools/call" {
		t.Fatalf("expected tools/call, got %q", handler.Method())
	}
}
