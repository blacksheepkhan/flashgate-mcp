package security

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	// ErrEmptyRoot is returned when the configured root path is empty.
	ErrEmptyRoot = errors.New("root path must not be empty")

	// ErrAbsolutePath is returned when a user path is absolute.
	ErrAbsolutePath = errors.New("absolute paths are not allowed")

	// ErrPathTraversal is returned when a user path attempts to escape the root.
	ErrPathTraversal = errors.New("path traversal detected")

	// ErrOutsideRoot is returned when the resolved path is outside the configured root.
	ErrOutsideRoot = errors.New("access outside root denied")
)

// SafePath represents a filesystem path that has passed validation.
type SafePath struct {
	path string
}

// String returns the absolute safe path.
func (p SafePath) String() string {
	return p.path
}

// Base returns the last path element.
func (p SafePath) Base() string {
	return filepath.Base(p.path)
}

// Dir returns the directory component.
func (p SafePath) Dir() string {
	return filepath.Dir(p.path)
}

// PathGuard validates and resolves user-provided paths against a sandbox root.
type PathGuard struct {
	root string
}

// NewPathGuard creates a new PathGuard.
func NewPathGuard(root string) (*PathGuard, error) {
	if strings.TrimSpace(root) == "" {
		return nil, ErrEmptyRoot
	}

	absoluteRoot, err := filepath.Abs(filepath.Clean(root))
	if err != nil {
		return nil, err
	}

	return &PathGuard{
		root: absoluteRoot,
	}, nil
}

// Root returns the absolute sandbox root.
func (g *PathGuard) Root() string {
	return g.root
}

// Resolve validates and resolves a user path against the sandbox root.
func (g *PathGuard) Resolve(userPath string) (SafePath, error) {
	normalizedUserPath := normalizeUserPath(userPath)

	if filepath.IsAbs(normalizedUserPath) {
		return SafePath{}, ErrAbsolutePath
	}

	if isTraversalPath(normalizedUserPath) {
		return SafePath{}, ErrPathTraversal
	}

	resolvedPath, err := filepath.Abs(filepath.Join(g.root, normalizedUserPath))
	if err != nil {
		return SafePath{}, err
	}

	if !isInsideRoot(g.root, resolvedPath) {
		return SafePath{}, ErrOutsideRoot
	}

	return SafePath{
		path: resolvedPath,
	}, nil
}

func normalizeUserPath(userPath string) string {
	if strings.TrimSpace(userPath) == "" {
		return "."
	}

	return filepath.Clean(userPath)
}

func isTraversalPath(path string) bool {
	return path == ".." || strings.HasPrefix(path, ".."+string(filepath.Separator))
}

func isInsideRoot(root string, resolvedPath string) bool {
	relative, err := filepath.Rel(root, resolvedPath)
	if err != nil {
		return false
	}

	return relative == "." || (relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)))
}
