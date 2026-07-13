package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	mcpserver "github.com/blacksheepkhan/flashgate-mcp/internal/mcp/server"
	"github.com/blacksheepkhan/flashgate-mcp/internal/protocol"
)

func TestToolsListWireOutputSchemasAndPayloadSizes(t *testing.T) {
	tests := []struct {
		name         string
		capabilities toolCapabilities
		toolCount    int
	}{
		{"read-only", capabilitiesFromReadOnly(true), 3},
		{"default", toolCapabilities{filesystemWrite: true}, 8},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			registry := createToolRegistry(noopFileSystem{}, 1024, tc.capabilities)
			request := bytes.NewBufferString("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\",\"params\":{}}\n")
			output := &bytes.Buffer{}
			server := mcpserver.New(request, output, createRouter("test-server", "test-version", registry))
			if err := server.Run(context.Background()); err != nil {
				t.Fatal(err)
			}

			var response struct {
				JSONRPC string          `json:"jsonrpc"`
				ID      json.RawMessage `json:"id"`
				Result  struct {
					Tools []protocol.Tool `json:"tools"`
				} `json:"result"`
			}
			if err := json.Unmarshal(output.Bytes(), &response); err != nil {
				t.Fatalf("invalid tools/list response: %v", err)
			}
			if len(response.Result.Tools) != tc.toolCount {
				t.Fatalf("got %d tools, want %d", len(response.Result.Tools), tc.toolCount)
			}
			schemaCount := 0
			for _, tool := range response.Result.Tools {
				if tool.OutputSchema == nil {
					t.Fatalf("tools/list omitted outputSchema for %s", tool.Name)
				}
				if tool.OutputSchema["type"] != "object" {
					t.Fatalf("%s outputSchema root type=%#v", tool.Name, tool.OutputSchema["type"])
				}
				schemaCount++
			}

			withSchemas := output.Len()
			for index := range response.Result.Tools {
				response.Result.Tools[index].OutputSchema = nil
			}
			historical, err := json.Marshal(response)
			if err != nil {
				t.Fatal(err)
			}
			withoutSchemas := len(historical) + 1 // STDIO JSONL newline, matching the runtime response.
			delta := withSchemas - withoutSchemas
			percent := float64(delta) * 100 / float64(withoutSchemas)
			t.Logf("profile=%s without=%dB with=%dB delta=%dB change=%.2f%% tools=%d schemas=%d",
				tc.name, withoutSchemas, withSchemas, delta, percent, tc.toolCount, schemaCount)
		})
	}
}
