package tools

// Package tools implements MCP tool definitions and tool execution.
//
// Each filesystem operation is exposed as a dedicated Tool implementation. A
// tool is responsible for defining its name, title, description, input schema,
// argument decoding, and result shape.
//
// Tools must keep path validation centralized by delegating filesystem access
// to the fs.FileSystem interface. They should not call os package filesystem
// functions directly.
//
// Tool results should be structured JSON-compatible objects. Errors from the
// filesystem layer are normalized into protocol errors before being returned to
// the MCP client.
//
// The Registry stores tools by name and returns them in deterministic
// registration order so that tools/list is stable for clients, debugging, and
// tests.
