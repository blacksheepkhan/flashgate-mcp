package fs

// Package fs provides the filesystem abstraction used by MCP tools.
//
// The package exposes a FileSystem interface for listing, reading, writing,
// creating, deleting, moving, copying, renaming, and inspecting paths below a
// configured root directory.
//
// Implementations are responsible for enforcing filesystem safety rules before
// touching the host filesystem. Callers pass relative paths; implementations
// resolve and validate those paths through the security layer.
//
// MCP tools depend on the FileSystem interface rather than direct operating
// system calls. This keeps path validation centralized and makes tool behavior
// testable through fake filesystem implementations.
