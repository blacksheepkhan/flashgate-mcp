package config

// Package config loads and exposes runtime configuration for FlashGate MCP.
//
// The package is responsible for translating environment variables and other
// configuration sources into typed configuration values used by the server.
//
// Configuration values are exposed through narrow accessor methods instead of
// mutable public fields. This keeps initialization explicit and avoids spreading
// environment parsing throughout the codebase.
//
// MCP_ROOT is required and production roots must be absolute. The process
// working directory is accepted only when MCP_ROOT=. and the explicit
// MCP_ALLOW_CWD_ROOT=true development option are both set. All filesystem
// operations exposed through MCP tools are ultimately constrained to the
// validated root by the filesystem and security layers.
