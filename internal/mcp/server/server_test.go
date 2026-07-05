package server

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/handlers"
	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/router"
	"github.com/blacksheepkhan/fileserver-mcp/internal/protocol"
)

func TestServerRunHandlesValidRequest(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"test/ok","params":{"value":"abc"}}` + "\n")
	output := &bytes.Buffer{}

	testRouter := router.New()
	testRouter.Register(&testHandler{
		method: "test/ok",
		result: map[string]any{
			"ok": true,
		},
	})

	server := New(input, output, testRouter)

	if err := server.Run(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	response := decodeSingleResponse(t, output.String())

	if response.JSONRPC != protocol.JSONRPCVersion {
		t.Fatalf("expected JSON-RPC version %q, got %q", protocol.JSONRPCVersion, response.JSONRPC)
	}

	if string(response.ID) != "1" {
		t.Fatalf("expected id 1, got %s", string(response.ID))
	}

	if response.Error != nil {
		t.Fatalf("expected no error, got %#v", response.Error)
	}

	result, ok := response.Result.(map[string]any)
	if !ok {
		t.Fatalf("expected result map, got %#v", response.Result)
	}

	if result["ok"] != true {
		t.Fatalf("expected ok=true, got %#v", result["ok"])
	}
}

func TestServerRunReturnsMethodNotFound(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"missing/method"}` + "\n")
	output := &bytes.Buffer{}

	server := New(input, output, router.New())

	if err := server.Run(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	response := decodeSingleResponse(t, output.String())

	if string(response.ID) != "1" {
		t.Fatalf("expected id 1, got %s", string(response.ID))
	}

	if response.Error == nil {
		t.Fatal("expected error response")
	}

	if response.Error.Code != protocol.ErrMethodNotFound {
		t.Fatalf("expected ErrMethodNotFound, got %d", response.Error.Code)
	}
}

func TestServerRunReturnsParseErrorForInvalidJSON(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":` + "\n")
	output := &bytes.Buffer{}

	server := New(input, output, router.New())

	if err := server.Run(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	response := decodeSingleResponse(t, output.String())

	if response.Error == nil {
		t.Fatal("expected parse error")
	}

	if response.Error.Code != protocol.ErrParseError {
		t.Fatalf("expected ErrParseError, got %d", response.Error.Code)
	}
}

func TestServerRunReturnsHandlerError(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"test/error"}` + "\n")
	output := &bytes.Buffer{}

	testRouter := router.New()
	testRouter.Register(&testHandler{
		method: "test/error",
		err: &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "invalid params",
		},
	})

	server := New(input, output, testRouter)

	if err := server.Run(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	response := decodeSingleResponse(t, output.String())

	if response.Error == nil {
		t.Fatal("expected error response")
	}

	if response.Error.Code != protocol.ErrInvalidParams {
		t.Fatalf("expected ErrInvalidParams, got %d", response.Error.Code)
	}

	if response.Error.Message != "invalid params" {
		t.Fatalf("expected error message %q, got %q", "invalid params", response.Error.Message)
	}
}

func TestServerRunHandlesMultipleRequests(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(
		`{"jsonrpc":"2.0","id":1,"method":"test/one"}` + "\n" +
			`{"jsonrpc":"2.0","id":2,"method":"test/two"}` + "\n",
	)
	output := &bytes.Buffer{}

	testRouter := router.New()
	testRouter.Register(&testHandler{
		method: "test/one",
		result: map[string]any{
			"name": "one",
		},
	})
	testRouter.Register(&testHandler{
		method: "test/two",
		result: map[string]any{
			"name": "two",
		},
	})

	server := New(input, output, testRouter)

	if err := server.Run(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	responses := decodeResponses(t, output.String())

	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}

	if string(responses[0].ID) != "1" {
		t.Fatalf("expected first id 1, got %s", string(responses[0].ID))
	}

	if string(responses[1].ID) != "2" {
		t.Fatalf("expected second id 2, got %s", string(responses[1].ID))
	}
}

func TestServerRunReturnsNilOnEOF(t *testing.T) {
	t.Parallel()

	input := strings.NewReader("")
	output := &bytes.Buffer{}

	server := New(input, output, router.New())

	if err := server.Run(context.Background()); err != nil {
		t.Fatalf("expected nil on EOF, got %v", err)
	}

	if output.Len() != 0 {
		t.Fatalf("expected no output, got %q", output.String())
	}
}

type testHandler struct {
	method string
	result any
	err    *protocol.Error
}

func (h *testHandler) Method() string {
	return h.method
}

func (h *testHandler) Handle(_ handlers.Context, _ json.RawMessage) (any, *protocol.Error) {
	if h.err != nil {
		return nil, h.err
	}

	return h.result, nil
}

type decodedResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *protocol.Error `json:"error,omitempty"`
}

func decodeSingleResponse(t *testing.T, output string) decodedResponse {
	t.Helper()

	responses := decodeResponses(t, output)

	if len(responses) != 1 {
		t.Fatalf("expected 1 response, got %d: %q", len(responses), output)
	}

	return responses[0]
}

func decodeResponses(t *testing.T, output string) []decodedResponse {
	t.Helper()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	responses := make([]decodedResponse, 0, len(lines))

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var response decodedResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			t.Fatalf("expected valid response JSON, got %v: %q", err, line)
		}

		responses = append(responses, response)
	}

	return responses
}
