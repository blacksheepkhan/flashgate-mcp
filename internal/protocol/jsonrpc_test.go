package protocol

import (
	"encoding/json"
	"testing"
)

func TestRequestUnmarshal(t *testing.T) {
	t.Parallel()

	const input = `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{"cursor":"abc"}}`

	var request Request
	if err := json.Unmarshal([]byte(input), &request); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if request.JSONRPC != JSONRPCVersion {
		t.Fatalf("expected JSON-RPC version %q, got %q", JSONRPCVersion, request.JSONRPC)
	}

	if string(request.ID) != "1" {
		t.Fatalf("expected id %q, got %q", "1", string(request.ID))
	}

	if request.Method != "tools/list" {
		t.Fatalf("expected method %q, got %q", "tools/list", request.Method)
	}

	if string(request.Params) != `{"cursor":"abc"}` {
		t.Fatalf("unexpected params: %s", string(request.Params))
	}
}

func TestRequestUnmarshalStringID(t *testing.T) {
	t.Parallel()

	const input = `{"jsonrpc":"2.0","id":"request-1","method":"initialize"}`

	var request Request
	if err := json.Unmarshal([]byte(input), &request); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(request.ID) != `"request-1"` {
		t.Fatalf("expected string id, got %s", string(request.ID))
	}
}

func TestNotificationUnmarshal(t *testing.T) {
	t.Parallel()

	const input = `{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}`

	var notification Notification
	if err := json.Unmarshal([]byte(input), &notification); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if notification.JSONRPC != JSONRPCVersion {
		t.Fatalf("expected JSON-RPC version %q, got %q", JSONRPCVersion, notification.JSONRPC)
	}

	if notification.Method != "notifications/initialized" {
		t.Fatalf("expected method %q, got %q", "notifications/initialized", notification.Method)
	}

	if string(notification.Params) != "{}" {
		t.Fatalf("expected empty params object, got %s", string(notification.Params))
	}
}

func TestResponseMarshalResult(t *testing.T) {
	t.Parallel()

	response := Response{
		JSONRPC: JSONRPCVersion,
		ID:      json.RawMessage(`1`),
		Result: map[string]any{
			"ok": true,
		},
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}

	if decoded["jsonrpc"] != JSONRPCVersion {
		t.Fatalf("expected jsonrpc %q, got %v", JSONRPCVersion, decoded["jsonrpc"])
	}

	if decoded["id"] != float64(1) {
		t.Fatalf("expected id 1, got %v", decoded["id"])
	}

	result, ok := decoded["result"].(map[string]any)
	if !ok {
		t.Fatalf("expected result object, got %#v", decoded["result"])
	}

	if result["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", result["ok"])
	}

	if _, exists := decoded["error"]; exists {
		t.Fatal("did not expect error field")
	}
}

func TestResponseMarshalError(t *testing.T) {
	t.Parallel()

	response := Response{
		JSONRPC: JSONRPCVersion,
		ID:      json.RawMessage(`1`),
		Error: &Error{
			Code:    ErrMethodNotFound,
			Message: "method not found",
		},
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var decoded struct {
		JSONRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Error   Error  `json:"error"`
	}

	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}

	if decoded.JSONRPC != JSONRPCVersion {
		t.Fatalf("expected jsonrpc %q, got %q", JSONRPCVersion, decoded.JSONRPC)
	}

	if decoded.ID != 1 {
		t.Fatalf("expected id 1, got %d", decoded.ID)
	}

	if decoded.Error.Code != ErrMethodNotFound {
		t.Fatalf("expected error code %d, got %d", ErrMethodNotFound, decoded.Error.Code)
	}

	if decoded.Error.Message != "method not found" {
		t.Fatalf("expected error message, got %q", decoded.Error.Message)
	}
}

func TestErrorMarshalWithData(t *testing.T) {
	t.Parallel()

	responseError := Error{
		Code:    ErrInvalidParams,
		Message: "invalid params",
		Data:    json.RawMessage(`{"field":"path"}`),
	}

	encoded, err := json.Marshal(responseError)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	const expected = `{"code":-32602,"message":"invalid params","data":{"field":"path"}}`
	if string(encoded) != expected {
		t.Fatalf("expected %s, got %s", expected, string(encoded))
	}
}
