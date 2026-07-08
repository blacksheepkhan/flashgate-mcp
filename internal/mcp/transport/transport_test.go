package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestReadMessageReadsSingleLine(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(`{"jsonrpc":"2.0","method":"tools/list"}` + "\n")
	output := &bytes.Buffer{}
	transport := New(input, output)

	message, err := transport.ReadMessage()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := `{"jsonrpc":"2.0","method":"tools/list"}` + "\n"
	if string(message) != expected {
		t.Fatalf("expected %q, got %q", expected, string(message))
	}
}

func TestReadMessageReadsMultipleLines(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(
		`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n" +
			`{"jsonrpc":"2.0","id":2,"method":"tools/list"}` + "\n",
	)
	output := &bytes.Buffer{}
	transport := New(input, output)

	first, err := transport.ReadMessage()
	if err != nil {
		t.Fatalf("expected no error for first message, got %v", err)
	}

	second, err := transport.ReadMessage()
	if err != nil {
		t.Fatalf("expected no error for second message, got %v", err)
	}

	expectedFirst := `{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"
	expectedSecond := `{"jsonrpc":"2.0","id":2,"method":"tools/list"}` + "\n"

	if string(first) != expectedFirst {
		t.Fatalf("expected first message %q, got %q", expectedFirst, string(first))
	}

	if string(second) != expectedSecond {
		t.Fatalf("expected second message %q, got %q", expectedSecond, string(second))
	}
}

func TestReadMessageReturnsEOF(t *testing.T) {
	t.Parallel()

	input := strings.NewReader("")
	output := &bytes.Buffer{}
	transport := New(input, output)

	_, err := transport.ReadMessage()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestReadMessageReturnsMessageTooLarge(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(`{"jsonrpc":"2.0","method":"tools/list"}` + "\n")
	output := &bytes.Buffer{}
	transport := NewWithLimits(input, output, 10, 0)

	_, err := transport.ReadMessage()
	if !errors.Is(err, ErrMessageTooLarge) {
		t.Fatalf("expected ErrMessageTooLarge, got %v", err)
	}
}

func TestWriteMessageWritesJSONLine(t *testing.T) {
	t.Parallel()

	input := strings.NewReader("")
	output := &bytes.Buffer{}
	transport := New(input, output)

	value := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"result": map[string]any{
			"ok": true,
		},
	}

	if err := transport.WriteMessage(value); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	written := output.String()

	if !strings.HasSuffix(written, "\n") {
		t.Fatalf("expected output to end with newline, got %q", written)
	}

	var decoded map[string]any
	if err := json.Unmarshal([]byte(written), &decoded); err != nil {
		t.Fatalf("expected valid JSON, got %v", err)
	}

	if decoded["jsonrpc"] != "2.0" {
		t.Fatalf("expected jsonrpc 2.0, got %v", decoded["jsonrpc"])
	}

	if decoded["id"] != float64(1) {
		t.Fatalf("expected id 1, got %v", decoded["id"])
	}
}

func TestWriteMessageFlushesOutput(t *testing.T) {
	t.Parallel()

	input := strings.NewReader("")
	output := &bytes.Buffer{}
	transport := New(input, output)

	if err := transport.WriteMessage(map[string]any{"ok": true}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.Len() == 0 {
		t.Fatal("expected output to be flushed")
	}
}

func TestWriteMessageReturnsResponseTooLarge(t *testing.T) {
	t.Parallel()

	input := strings.NewReader("")
	output := &bytes.Buffer{}
	transport := NewWithLimits(input, output, 0, 10)

	err := transport.WriteMessage(map[string]any{"content": strings.Repeat("x", 50)})
	if !errors.Is(err, ErrResponseTooLarge) {
		t.Fatalf("expected ErrResponseTooLarge, got %v", err)
	}

	if output.Len() != 0 {
		t.Fatalf("expected no oversized output, got %q", output.String())
	}
}
