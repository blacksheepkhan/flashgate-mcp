package transport

// Package transport implements JSON-RPC message transport concerns.
//
// The transport layer is responsible for reading protocol messages from input
// streams and writing protocol responses to output streams. It does not route
// methods, execute tools, or access the filesystem.
//
// For the server process, standard output is reserved for JSON-RPC protocol
// messages. Diagnostic output must be written outside the protocol stream.
