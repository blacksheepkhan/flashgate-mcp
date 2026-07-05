package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/blacksheepkhan/fileserver-mcp/internal/protocol"
)

func TestRegistryRegisterAndGet(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	tool := &testTool{
		name: "test_tool",
	}

	registry.Register(tool)

	got, ok := registry.Get("test_tool")
	if !ok {
		t.Fatal("expected tool to exist")
	}

	if got != tool {
		t.Fatalf("expected registered tool, got %#v", got)
	}
}

func TestRegistryGetUnknownTool(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	got, ok := registry.Get("missing_tool")
	if ok {
		t.Fatal("expected tool not to exist")
	}

	if got != nil {
		t.Fatalf("expected nil tool, got %#v", got)
	}
}

func TestRegistryList(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	registry.Register(&testTool{name: "first"})
	registry.Register(&testTool{name: "second"})

	list := registry.List()

	if len(list) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(list))
	}

	names := map[string]bool{}
	for _, tool := range list {
		names[tool.Name()] = true
	}

	if !names["first"] {
		t.Fatal("expected first tool")
	}

	if !names["second"] {
		t.Fatal("expected second tool")
	}
}

func TestRegistryRegisterOverridesExistingTool(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	first := &testTool{
		name:        "test_tool",
		description: "first",
	}

	second := &testTool{
		name:        "test_tool",
		description: "second",
	}

	registry.Register(first)
	registry.Register(second)

	got, ok := registry.Get("test_tool")
	if !ok {
		t.Fatal("expected tool to exist")
	}

	if got != second {
		t.Fatalf("expected second tool, got %#v", got)
	}
}

type testTool struct {
	name        string
	title       string
	description string
	inputSchema any
	result      any
	err         *protocol.Error
	called      int
	arguments   json.RawMessage
}

func (t *testTool) Name() string {
	return t.name
}

func (t *testTool) Title() string {
	return t.title
}

func (t *testTool) Description() string {
	return t.description
}

func (t *testTool) InputSchema() any {
	return t.inputSchema
}

func (t *testTool) Definition() protocol.Tool {
	return protocol.Tool{
		Name:        t.name,
		Title:       t.title,
		Description: t.description,
		InputSchema: t.inputSchema,
	}
}

func (t *testTool) Execute(_ context.Context, arguments json.RawMessage) (any, *protocol.Error) {
	t.called++
	t.arguments = append(json.RawMessage(nil), arguments...)

	if t.err != nil {
		return nil, t.err
	}

	return t.result, nil
}
