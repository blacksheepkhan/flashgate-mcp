package security

import (
	"errors"
	"os"
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
	root          string
	effectiveRoot string
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

	effectiveRoot, err := filepath.EvalSymlinks(absoluteRoot)
	if err != nil {
		return nil, err
	}

	return &PathGuard{
		root:          absoluteRoot,
		effectiveRoot: effectiveRoot,
	}, nil
}

// Root returns the absolute sandbox root.
func (g *PathGuard) Root() string {
	return g.root
}

// Resolve validates and resolves a user path against the sandbox root.
func (g *PathGuard) Resolve(userPath string) (SafePath, error) {
	return g.ResolveForCreate(userPath)
}

// ResolveExisting validates and resolves an existing user path against the sandbox root.
func (g *PathGuard) ResolveExisting(userPath string) (SafePath, error) {
	resolvedPath, err := g.resolveLexical(userPath)
	if err != nil {
		return SafePath{}, err
	}

	effectivePath, err := filepath.EvalSymlinks(resolvedPath)
	if err != nil {
		return SafePath{}, err
	}

	if !isInsideRoot(g.effectiveRoot, effectivePath) {
		return SafePath{}, ErrOutsideRoot
	}

	return SafePath{
		path: resolvedPath,
	}, nil
}

// ResolveForCreate validates a user path for creation or replacement.
func (g *PathGuard) ResolveForCreate(userPath string) (SafePath, error) {
	resolvedPath, err := g.resolveLexical(userPath)
	if err != nil {
		return SafePath{}, err
	}

	effectivePath, err := filepath.EvalSymlinks(resolvedPath)
	if err == nil {
		if !isInsideRoot(g.effectiveRoot, effectivePath) {
			return SafePath{}, ErrOutsideRoot
		}

		return SafePath{
			path: resolvedPath,
		}, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return SafePath{}, err
	}

	effectiveParent, err := g.resolveExistingParent(resolvedPath)
	if err != nil {
		return SafePath{}, err
	}

	if !isInsideRoot(g.effectiveRoot, effectiveParent) {
		return SafePath{}, ErrOutsideRoot
	}

	return SafePath{
		path: resolvedPath,
	}, nil
}

func (g *PathGuard) resolveLexical(userPath string) (string, error) {
	normalizedUserPath := normalizeUserPath(userPath)

	if filepath.IsAbs(normalizedUserPath) {
		return "", ErrAbsolutePath
	}

	if isTraversalPath(normalizedUserPath) {
		return "", ErrPathTraversal
	}

	resolvedPath, err := filepath.Abs(filepath.Join(g.root, normalizedUserPath))
	if err != nil {
		return "", err
	}

	if !isInsideRoot(g.root, resolvedPath) {
		return "", ErrOutsideRoot
	}

	return resolvedPath, nil
}

func (g *PathGuard) resolveExistingParent(path string) (string, error) {
	current := filepath.Clean(path)

	for {
		parent := filepath.Dir(current)
		if parent == current {
			return "", ErrOutsideRoot
		}

		if !isInsideRoot(g.root, parent) {
			return "", ErrOutsideRoot
		}

		effectiveParent, err := filepath.EvalSymlinks(parent)
		if err == nil {
			return effectiveParent, nil
		}

		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		current = parent
	}
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
