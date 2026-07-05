package protocol

import "testing"

func TestJSONRPCErrorCodes(t *testing.T) {
	t.Parallel()

	tests := map[string]int{
		"parse error":      ErrParseError,
		"invalid request":  ErrInvalidRequest,
		"method not found": ErrMethodNotFound,
		"invalid params":   ErrInvalidParams,
		"internal error":   ErrInternalError,
	}

	expected := map[string]int{
		"parse error":      -32700,
		"invalid request":  -32600,
		"method not found": -32601,
		"invalid params":   -32602,
		"internal error":   -32603,
	}

	for name, code := range tests {
		if code != expected[name] {
			t.Fatalf("expected %s to be %d, got %d", name, expected[name], code)
		}
	}
}
