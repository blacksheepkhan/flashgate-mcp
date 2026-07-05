package router

// Package router dispatches JSON-RPC method calls to registered MCP handlers.
//
// The router maps method names such as initialize, tools/list, and tools/call
// to handler implementations. It does not perform protocol IO and does not
// execute filesystem operations directly.
//
// Unknown methods and invalid requests are converted into protocol-level
// errors by the routing and handling layers.
