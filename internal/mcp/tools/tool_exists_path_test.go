package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/blacksheepkhan/fileserver-mcp/internal/fs"
	"github.com/blacksheepkhan/fileserver-mcp/internal/protocol"
)

func TestExistsPathToolDefinition(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewExistsPathTool(filesystem)

	if tool.Name() != "exists_path" {
		t.Fatalf("expected exists_path, got %q", tool.Name())
	}

	if tool.Description() == "" {
		t.Fatal("expected description")
	}

	if tool.InputSchema() == nil {
		t.Fatal("expected input schema")
	}

	definition := tool.Definition()

	if definition.Name != "exists_path" {
		t.Fatalf("expected definition name exists_path, got %q", definition.Name)
	}

	if definition.Description == "" {
		t.Fatal("expected definition description")
	}

	if definition.InputSchema == nil {
		t.Fatal("expected definition input schema")
	}
}

func TestExistsPathToolExecuteExistingPath(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	filesystem.exists = true

	tool := NewExistsPathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"path":"README.md"}`),
	)

	if rpcErr != nil {
		t.Fatalf("expected no error, got %v", rpcErr)
	}

	if filesystem.existsPath != "README.md" {
		t.Fatalf("expected path README.md, got %q", filesystem.existsPath)
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("expected marshalable result, got %v", err)
	}

	var decoded struct {
		Exists bool `json:"exists"`
	}

	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid result JSON, got %v", err)
	}

	if !decoded.Exists {
		t.Fatal("expected path to exist")
	}
}

func TestExistsPathToolExecuteMissingPath(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	filesystem.exists = false

	tool := NewExistsPathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"path":"missing.txt"}`),
	)

	if rpcErr != nil {
		t.Fatalf("expected no error, got %v", rpcErr)
	}

	if filesystem.existsPath != "missing.txt" {
		t.Fatalf("expected path missing.txt, got %q", filesystem.existsPath)
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("expected marshalable result, got %v", err)
	}

	var decoded struct {
		Exists bool `json:"exists"`
	}

	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid result JSON, got %v", err)
	}

	if decoded.Exists {
		t.Fatal("expected path not to exist")
	}
}

func TestExistsPathToolReturnsInvalidParamsForMalformedJSON(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewExistsPathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"path":`),
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

func TestExistsPathToolReturnsInvalidParamsForMissingPathArgument(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewExistsPathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
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

func TestExistsPathToolReturnsInvalidParamsForFilesystemError(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	filesystem.existsErr = fs.ErrPathIsNotDirectory

	tool := NewExistsPathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"path":"invalid"}`),
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
