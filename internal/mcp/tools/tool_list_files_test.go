package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/blacksheepkhan/flashgate-mcp/internal/fs"
	"github.com/blacksheepkhan/flashgate-mcp/internal/protocol"
	"github.com/blacksheepkhan/flashgate-mcp/internal/security"
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

func TestMapFilesystemErrorReturnsGenericInvalidParamsForSecurityDenial(t *testing.T) {
	t.Parallel()

	hostPath := t.TempDir()
	rpcErr := mapFilesystemError(fmt.Errorf("%w: %s", security.ErrOutsideRoot, hostPath))

	if rpcErr.Code != protocol.ErrInvalidParams {
		t.Fatalf("expected ErrInvalidParams, got %d", rpcErr.Code)
	}

	if strings.Contains(rpcErr.Message, hostPath) {
		t.Fatalf("expected generic error message without host path, got %q", rpcErr.Message)
	}

	if rpcErr.Message != "filesystem error: invalid path" {
		t.Fatalf("expected generic invalid path message, got %q", rpcErr.Message)
	}
}

func TestMapFilesystemErrorReturnsGenericInvalidParamsForPolicyDenials(t *testing.T) {
	t.Parallel()

	hostPath := t.TempDir()
	testCases := []error{
		security.ErrHiddenPathDenied,
		security.ErrUNCPathDenied,
		security.ErrSymlinkDenied,
		security.ErrReparsePointDenied,
	}

	for _, testErr := range testCases {
		rpcErr := mapFilesystemError(fmt.Errorf("%w: %s", testErr, hostPath))

		if rpcErr.Code != protocol.ErrInvalidParams {
			t.Fatalf("expected ErrInvalidParams for %v, got %d", testErr, rpcErr.Code)
		}

		if rpcErr.Message != "filesystem error: invalid path" {
			t.Fatalf("expected generic invalid path message for %v, got %q", testErr, rpcErr.Message)
		}

		if strings.Contains(rpcErr.Message, hostPath) {
			t.Fatalf("expected message without host path for %v, got %q", testErr, rpcErr.Message)
		}
	}
}

type fakeFileSystem struct {
	entries         []fs.Entry
	err             error
	listPath        string
	readPath        string
	readMaxBytes    int64
	readContent     []byte
	readErr         error
	statPath        string
	statMetadata    fs.Metadata
	statErr         error
	existsPath      string
	exists          bool
	existsErr       error
	writePath       string
	writeContent    []byte
	writeOverwrite  bool
	writeErr        error
	mkdirPath       string
	mkdirErr        error
	deletePath      string
	deleteRecursive bool
	deleteErr       error
	moveSource      string
	moveTarget      string
	moveOverwrite   bool
	moveErr         error
	copySource      string
	copyTarget      string
	copyOverwrite   bool
	copyErr         error
	renameSource    string
	renameTarget    string
	renameOverwrite bool
	renameErr       error
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

func (f *fakeFileSystem) Stat(path string) (fs.Metadata, error) {
	f.statPath = path

	if f.statErr != nil {
		return fs.Metadata{}, f.statErr
	}

	return f.statMetadata, nil
}

func (f *fakeFileSystem) Exists(path string) (bool, error) {
	f.existsPath = path

	if f.existsErr != nil {
		return false, f.existsErr
	}

	return f.exists, nil
}

func (f *fakeFileSystem) Write(path string, content []byte, overwrite bool) error {
	f.writePath = path
	f.writeContent = append([]byte(nil), content...)
	f.writeOverwrite = overwrite

	if f.writeErr != nil {
		return f.writeErr
	}

	return nil
}

func (f *fakeFileSystem) Mkdir(path string) error {
	f.mkdirPath = path

	if f.mkdirErr != nil {
		return f.mkdirErr
	}

	return nil
}

func (f *fakeFileSystem) Delete(path string, recursive bool) error {
	f.deletePath = path
	f.deleteRecursive = recursive

	if f.deleteErr != nil {
		return f.deleteErr
	}

	return nil
}

func (f *fakeFileSystem) Move(source string, target string, overwrite bool) error {
	f.moveSource = source
	f.moveTarget = target
	f.moveOverwrite = overwrite

	if f.moveErr != nil {
		return f.moveErr
	}

	return nil
}

func (f *fakeFileSystem) Copy(source string, target string, overwrite bool) error {
	f.copySource = source
	f.copyTarget = target
	f.copyOverwrite = overwrite

	if f.copyErr != nil {
		return f.copyErr
	}

	return nil
}

func (f *fakeFileSystem) Rename(source string, target string, overwrite bool) error {
	f.renameSource = source
	f.renameTarget = target
	f.renameOverwrite = overwrite

	if f.renameErr != nil {
		return f.renameErr
	}

	return nil
}
