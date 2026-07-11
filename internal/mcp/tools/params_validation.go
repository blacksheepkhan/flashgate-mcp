package tools

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/blacksheepkhan/flashgate-mcp/internal/protocol"
)

func invalidParamsError() *protocol.Error {
	return &protocol.Error{
		Code:    protocol.ErrInvalidParams,
		Message: "invalid params",
	}
}

func isMissingNullOrObject(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)

	return len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) || isJSONObject(raw)
}

func isJSONObject(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)

	return len(trimmed) > 0 && trimmed[0] == '{'
}

func isJSONNull(raw json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
}

func stringField(fields map[string]json.RawMessage, key string) (string, bool) {
	raw, ok := fields[key]
	if !ok {
		return "", false
	}

	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return "", false
	}

	if strings.TrimSpace(value) == "" {
		return "", false
	}

	return value, true
}
