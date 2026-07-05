package server

// Package server coordinates the MCP server loop.
//
// The server connects transport input and output with the JSON-RPC router. It
// owns the request processing loop but delegates method handling to registered
// MCP handlers.
//
// The package should remain small: protocol parsing belongs to the transport
// layer, method dispatch belongs to the router layer, and filesystem behavior
// belongs to tools and the filesystem abstraction.
