package protocol

// Package protocol defines JSON-RPC and MCP protocol data structures.
//
// The package contains request, response, error, and tool definition types used
// by the transport, router, handlers, and tools packages.
//
// It intentionally contains protocol shapes and constants only. It should not
// perform filesystem operations, process IO, routing, or business logic.
