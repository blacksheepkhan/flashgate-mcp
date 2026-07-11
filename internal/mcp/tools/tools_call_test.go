package tools

import (
	"encoding/json"
	"testing"

	"github.com/blacksheepkhan/flashgate-mcp/internal/mcp/handlers"
	"github.com/blacksheepkhan/flashgate-mcp/internal/protocol"
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

func TestCallHandlerReturnsInvalidParamsForUnknownTool(t *testing.T) {
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

	if rpcErr.Code != protocol.ErrInvalidParams {
		t.Fatalf("expected ErrInvalidParams, got %d", rpcErr.Code)
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

func TestCallHandlerReturnsInvalidParamsForInvalidParamsShape(t *testing.T) {
	t.Parallel()

	testCases := map[string]json.RawMessage{
		"missing params": nil,
		"null params":    json.RawMessage(`null`),
		"string params":  json.RawMessage(`"bad"`),
		"array params":   json.RawMessage(`[]`),
		"number params":  json.RawMessage(`1`),
		"bool params":    json.RawMessage(`true`),
		"non-string name": json.RawMessage(
			`{"name":123,"arguments":{}}`,
		),
		"empty name":           json.RawMessage(`{"name":"","arguments":{}}`),
		"non-object arguments": json.RawMessage(`{"name":"test_tool","arguments":"bad"}`),
	}

	for name, rawParams := range testCases {
		name := name
		rawParams := rawParams

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler := NewCallHandler(NewRegistry())

			result, rpcErr := handler.Handle(handlers.Context{}, rawParams)
			if result != nil {
				t.Fatalf("expected nil result, got %#v", result)
			}

			if rpcErr == nil {
				t.Fatal("expected rpc error")
			}

			if rpcErr.Code != protocol.ErrInvalidParams {
				t.Fatalf("expected ErrInvalidParams, got %d", rpcErr.Code)
			}

			if rpcErr.Message != "invalid params" {
				t.Fatalf("expected generic invalid params message, got %q", rpcErr.Message)
			}
		})
	}
}

func TestCallHandlerTreatsMissingAndNullArgumentsAsEmptyObject(t *testing.T) {
	t.Parallel()

	testCases := map[string]json.RawMessage{
		"missing arguments": json.RawMessage(`{"name":"test_tool"}`),
		"null arguments":    json.RawMessage(`{"name":"test_tool","arguments":null}`),
	}

	for name, rawParams := range testCases {
		name := name
		rawParams := rawParams

		t.Run(name, func(t *testing.T) {
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

			result, rpcErr := handler.Handle(handlers.Context{}, rawParams)
			if rpcErr != nil {
				t.Fatalf("expected no error, got %v", rpcErr)
			}

			if result == nil {
				t.Fatal("expected result")
			}

			if string(tool.arguments) != `{}` {
				t.Fatalf("expected empty object arguments, got %s", string(tool.arguments))
			}
		})
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
