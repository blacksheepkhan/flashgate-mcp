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
// The most important filesystem setting is the configured root directory. All
// filesystem operations exposed through MCP tools are ultimately constrained to
// that root by the filesystem and security layers.
