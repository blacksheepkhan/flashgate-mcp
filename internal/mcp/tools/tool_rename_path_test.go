package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/blacksheepkhan/fileserver-mcp/internal/fs"
	"github.com/blacksheepkhan/fileserver-mcp/internal/protocol"
)

func TestRenamePathToolDefinition(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewRenamePathTool(filesystem)

	if tool.Name() != "rename_path" {
		t.Fatalf("expected rename_path, got %q", tool.Name())
	}

	if tool.Description() == "" {
		t.Fatal("expected description")
	}

	if tool.InputSchema() == nil {
		t.Fatal("expected input schema")
	}

	definition := tool.Definition()

	if definition.Name != "rename_path" {
		t.Fatalf("expected definition name rename_path, got %q", definition.Name)
	}

	if definition.Description == "" {
		t.Fatal("expected definition description")
	}

	if definition.InputSchema == nil {
		t.Fatal("expected definition input schema")
	}
}

func TestRenamePathToolExecute(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewRenamePathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"source":"old-name.txt","target":"new-name.txt"}`),
	)

	if rpcErr != nil {
		t.Fatalf("expected no error, got %v", rpcErr)
	}

	if filesystem.renameSource != "old-name.txt" {
		t.Fatalf("expected source old-name.txt, got %q", filesystem.renameSource)
	}

	if filesystem.renameTarget != "new-name.txt" {
		t.Fatalf("expected target new-name.txt, got %q", filesystem.renameTarget)
	}

	if filesystem.renameOverwrite {
		t.Fatal("expected overwrite=false")
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("expected marshalable result, got %v", err)
	}

	var decoded struct {
		Source  string `json:"source"`
		Target  string `json:"target"`
		Renamed bool   `json:"renamed"`
	}

	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid result JSON, got %v", err)
	}

	if decoded.Source != "old-name.txt" {
		t.Fatalf("expected source old-name.txt, got %q", decoded.Source)
	}

	if decoded.Target != "new-name.txt" {
		t.Fatalf("expected target new-name.txt, got %q", decoded.Target)
	}

	if !decoded.Renamed {
		t.Fatal("expected renamed=true")
	}
}

func TestRenamePathToolForwardsOverwrite(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewRenamePathTool(filesystem)

	_, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"source":"old.txt","target":"new.txt","overwrite":true}`),
	)

	if rpcErr != nil {
		t.Fatalf("expected no error, got %v", rpcErr)
	}

	if !filesystem.renameOverwrite {
		t.Fatal("expected overwrite=true")
	}
}

func TestRenamePathToolReturnsInvalidParamsForMalformedJSON(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewRenamePathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"source":`),
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

func TestRenamePathToolReturnsInvalidParamsForMissingSource(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewRenamePathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"target":"new.txt"}`),
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

func TestRenamePathToolReturnsInvalidParamsForMissingTarget(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewRenamePathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"source":"old.txt"}`),
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

func TestRenamePathToolReturnsInvalidParamsForTargetExists(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	filesystem.renameErr = fs.ErrFileExists

	tool := NewRenamePathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"source":"old.txt","target":"new.txt","overwrite":false}`),
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

func TestRenamePathToolReturnsInvalidParamsForFilesystemError(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	filesystem.renameErr = fs.ErrPathIsDirectory

	tool := NewRenamePathTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"source":"invalid","target":"target.txt"}`),
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
