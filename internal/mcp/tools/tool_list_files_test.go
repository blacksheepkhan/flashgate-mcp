package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/blacksheepkhan/fileserver-mcp/internal/fs"
	"github.com/blacksheepkhan/fileserver-mcp/internal/protocol"
)

func TestListFilesToolDefinition(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewListFilesTool(filesystem)

	if tool.Name() != "list_files" {
		t.Fatalf("expected list_files, got %q", tool.Name())
	}

	if tool.Description() == "" {
		t.Fatal("expected description")
	}

	schema := tool.InputSchema()
	if schema == nil {
		t.Fatal("expected input schema")
	}

	definition := tool.Definition()

	if definition.Name != "list_files" {
		t.Fatalf("expected definition name list_files, got %q", definition.Name)
	}

	if definition.Description == "" {
		t.Fatal("expected definition description")
	}

	if definition.InputSchema == nil {
		t.Fatal("expected definition input schema")
	}
}

func TestListFilesToolExecute(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	filesystem.entries = []fs.Entry{
		{
			Name:  "docs",
			IsDir: true,
			Size:  0,
		},
		{
			Name:  "README.md",
			IsDir: false,
			Size:  128,
		},
	}

	tool := NewListFilesTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"path":"."}`),
	)

	if rpcErr != nil {
		t.Fatalf("expected no error, got %v", rpcErr)
	}

	if filesystem.listPath != "." {
		t.Fatalf("expected path '.', got %q", filesystem.listPath)
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("expected marshalable result, got %v", err)
	}

	var decoded struct {
		Entries []struct {
			Name  string `json:"name"`
			IsDir bool   `json:"isDir"`
			Size  int64  `json:"size"`
		} `json:"entries"`
	}

	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid result JSON, got %v", err)
	}

	if len(decoded.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(decoded.Entries))
	}

	if decoded.Entries[0].Name != "docs" {
		t.Fatalf("expected docs, got %q", decoded.Entries[0].Name)
	}

	if !decoded.Entries[0].IsDir {
		t.Fatal("expected docs to be directory")
	}

	if decoded.Entries[1].Name != "README.md" {
		t.Fatalf("expected README.md, got %q", decoded.Entries[1].Name)
	}

	if decoded.Entries[1].Size != 128 {
		t.Fatalf("expected size 128, got %d", decoded.Entries[1].Size)
	}
}

func TestListFilesToolExecuteDefaultsPathToDot(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewListFilesTool(filesystem)

	_, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)

	if rpcErr != nil {
		t.Fatalf("expected no error, got %v", rpcErr)
	}

	if filesystem.listPath != "." {
		t.Fatalf("expected default path '.', got %q", filesystem.listPath)
	}
}

func TestListFilesToolReturnsInvalidParamsForMalformedJSON(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	tool := NewListFilesTool(filesystem)

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

func TestListFilesToolReturnsInvalidParamsForFilesystemError(t *testing.T) {
	t.Parallel()

	filesystem := newFakeFileSystem()
	filesystem.err = fs.ErrPathIsNotDirectory

	tool := NewListFilesTool(filesystem)

	result, rpcErr := tool.Execute(
		context.Background(),
		json.RawMessage(`{"path":"file.txt"}`),
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

type fakeFileSystem struct {
	entries      []fs.Entry
	err          error
	listPath     string
	readPath     string
	readMaxBytes int64
	readContent  []byte
	readErr      error
}

func newFakeFileSystem() *fakeFileSystem {
	return &fakeFileSystem{}
}

func (f *fakeFileSystem) List(path string) ([]fs.Entry, error) {
	f.listPath = path

	if f.err != nil {
		return nil, f.err
	}

	return f.entries, nil
}

func (f *fakeFileSystem) Read(path string, maxBytes int64) ([]byte, error) {
	f.readPath = path
	f.readMaxBytes = maxBytes

	if f.readErr != nil {
		return nil, f.readErr
	}

	return f.readContent, nil
}

func (f *fakeFileSystem) Stat(_ string) (fs.Metadata, error) {
	panic("not implemented")
}

func (f *fakeFileSystem) Exists(_ string) (bool, error) {
	panic("not implemented")
}

func (f *fakeFileSystem) Write(_ string, _ []byte, _ bool) error {
	panic("not implemented")
}

func (f *fakeFileSystem) Mkdir(_ string) error {
	panic("not implemented")
}

func (f *fakeFileSystem) Delete(_ string, _ bool) error {
	panic("not implemented")
}

func (f *fakeFileSystem) Move(_ string, _ string, _ bool) error {
	panic("not implemented")
}

func (f *fakeFileSystem) Copy(_ string, _ string, _ bool) error {
	panic("not implemented")
}

func (f *fakeFileSystem) Rename(_ string, _ string, _ bool) error {
	panic("not implemented")
}
