package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/blacksheepkhan/flashgate-mcp/internal/fs"
	"github.com/blacksheepkhan/flashgate-mcp/internal/protocol"
)

func TestStatPathToolDefinition(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewStatPathTool(filesystem)

	if tool.Name() != "stat_path" {
		t.Fatalf("expected stat_path, got %q", tool.Name())
	}

	if tool.Description() == "" {
		t.Fatal("expected description")
	}

	if tool.InputSchema() == nil {
		t.Fatal("expected input schema")
	}

	definition := tool.Definition()

	if definition.Name != "stat_path" {
		t.Fatalf("expected definition name stat_path, got %q", definition.Name)
	}

	if definition.Description == "" {
		t.Fatal("expected definition description")
	}

	if definition.InputSchema == nil {
		t.Fatal("expected definition input schema")
	}
}

func TestStatPathToolExecuteFile(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	filesystem.statMetadata = fs.Metadata{
		Name:  "README.md",
		IsDir: false,
		Size:  3595,
	}

	tool := NewStatPathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"path":"README.md"}`),
	)

	if rpcErr != nil {
		t.Fatalf("expected no error, got %v", rpcErr)
	}

	if filesystem.statPath != "README.md" {
		t.Fatalf("expected path README.md, got %q", filesystem.statPath)
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("expected marshalable result, got %v", err)
	}

	var decoded struct {
		Name  string `json:"name"`
		IsDir bool   `json:"isDir"`
		Size  int64  `json:"size"`
	}

	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid result JSON, got %v", err)
	}

	if decoded.Name != "README.md" {
		t.Fatalf("expected README.md, got %q", decoded.Name)
	}

	if decoded.IsDir {
		t.Fatal("expected file, got directory")
	}

	if decoded.Size != 3595 {
		t.Fatalf("expected size 3595, got %d", decoded.Size)
	}
}

func TestStatPathToolExecuteDirectory(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	filesystem.statMetadata = fs.Metadata{
		Name:  "internal",
		IsDir: true,
		Size:  0,
	}

	tool := NewStatPathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"path":"internal"}`),
	)

	if rpcErr != nil {
		t.Fatalf("expected no error, got %v", rpcErr)
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("expected marshalable result, got %v", err)
	}

	var decoded struct {
		Name  string `json:"name"`
		IsDir bool   `json:"isDir"`
		Size  int64  `json:"size"`
	}

	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid result JSON, got %v", err)
	}

	if decoded.Name != "internal" {
		t.Fatalf("expected internal, got %q", decoded.Name)
	}

	if !decoded.IsDir {
		t.Fatal("expected directory")
	}
}

func TestStatPathToolReturnsInvalidParamsForMalformedJSON(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewStatPathTool(filesystem)

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

func TestStatPathToolReturnsInvalidParamsForMissingPath(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewStatPathTool(filesystem)

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

func TestStatPathToolReturnsInvalidParamsForFilesystemError(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	filesystem.statErr = fs.ErrPathIsNotDirectory

	tool := NewStatPathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"path":"missing"}`),
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
