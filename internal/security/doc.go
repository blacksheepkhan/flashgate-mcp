package security

// Package security contains path validation primitives used to constrain
// filesystem access.
//
// The central responsibility of this package is to ensure that user-supplied
// paths cannot escape the configured filesystem root. This includes rejecting
// absolute paths, parent traversal outside the root, and other unsafe path
// forms.
//
// Higher-level MCP tools must not duplicate or bypass this validation. They
// should pass relative paths to the filesystem layer and rely on this package,
// through the filesystem implementation, to enforce root confinement.
