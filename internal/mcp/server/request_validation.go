package server

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/blacksheepkhan/fileserver-mcp/internal/protocol"
)

type validatedRequest struct {
	id           json.RawMessage
	method       string
	params       json.RawMessage
	notification bool
}

type requestValidationError struct {
	id      json.RawMessage
	code    int
	message string
}

func validateRequestMessageWithLimits(message []byte, maxArgumentBytes int64) (validatedRequest, *requestValidationError) {
	if !json.Valid(message) {
		return validatedRequest{}, newRequestValidationError(nullID(), protocol.ErrParseError, "parse error")
	}

	trimmed := bytes.TrimSpace(message)
	if len(trimmed) == 0 {
		return validatedRequest{}, newRequestValidationError(nullID(), protocol.ErrParseError, "parse error")
	}

	if trimmed[0] != '{' {
		return validatedRequest{}, newRequestValidationError(nullID(), protocol.ErrInvalidRequest, "invalid request")
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(trimmed, &fields); err != nil {
		return validatedRequest{}, newRequestValidationError(nullID(), protocol.ErrInvalidRequest, "invalid request")
	}

	version, ok := stringField(fields, "jsonrpc")
	if !ok || version != protocol.JSONRPCVersion {
		return validatedRequest{}, newRequestValidationError(requestIDOrNull(fields), protocol.ErrInvalidRequest, "invalid request")
	}

	id, hasID, ok := requestID(fields)
	if !ok {
		return validatedRequest{}, newRequestValidationError(nullID(), protocol.ErrInvalidRequest, "invalid request")
	}

	method, ok := stringField(fields, "method")
	if !ok || strings.TrimSpace(method) == "" {
		return validatedRequest{}, newRequestValidationError(idForInvalidRequest(id, hasID), protocol.ErrInvalidRequest, "invalid request")
	}

	if !hasID {
		return validatedRequest{
			method:       method,
			params:       fields["params"],
			notification: true,
		}, nil
	}

	params, paramsErr := validateMethodParams(method, fields["params"], maxArgumentBytes)
	if paramsErr != nil {
		paramsErr.id = id
		return validatedRequest{}, paramsErr
	}

	return validatedRequest{
		id:     id,
		method: method,
		params: params,
	}, nil
}

func validateMethodParams(method string, params json.RawMessage, maxArgumentBytes int64) (json.RawMessage, *requestValidationError) {
	switch method {
	case "initialize":
		if !isJSONObject(params) {
			return nil, newRequestValidationError(nil, protocol.ErrInvalidParams, "invalid params")
		}

		var fields map[string]json.RawMessage
		if err := json.Unmarshal(params, &fields); err != nil {
			return nil, newRequestValidationError(nil, protocol.ErrInvalidParams, "invalid params")
		}

		version, ok := stringField(fields, "protocolVersion")
		if !ok || version == "" {
			return nil, newRequestValidationError(nil, protocol.ErrInvalidParams, "invalid params")
		}

		return params, nil
	case "tools/list":
		if len(bytes.TrimSpace(params)) == 0 || isJSONNull(params) || isJSONObject(params) {
			return params, nil
		}

		return nil, newRequestValidationError(nil, protocol.ErrInvalidParams, "invalid params")
	case "tools/call":
		return validateToolsCallParams(params, maxArgumentBytes)
	default:
		return params, nil
	}
}

func validateToolsCallParams(params json.RawMessage, maxArgumentBytes int64) (json.RawMessage, *requestValidationError) {
	if !isJSONObject(params) {
		return nil, newRequestValidationError(nil, protocol.ErrInvalidParams, "invalid params")
	}

	if exceedsLimit(params, maxArgumentBytes) {
		return nil, newRequestValidationError(nil, protocol.ErrInvalidParams, "invalid params")
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(params, &fields); err != nil {
		return nil, newRequestValidationError(nil, protocol.ErrInvalidParams, "invalid params")
	}

	name, ok := stringField(fields, "name")
	if !ok || strings.TrimSpace(name) == "" {
		return nil, newRequestValidationError(nil, protocol.ErrInvalidParams, "invalid params")
	}

	arguments, hasArguments := fields["arguments"]
	if !hasArguments || isJSONNull(arguments) {
		fields["arguments"] = json.RawMessage(`{}`)
		normalized, err := json.Marshal(fields)
		if err != nil {
			return nil, newRequestValidationError(nil, protocol.ErrInternalError, "internal error")
		}

		return normalized, nil
	}

	if !isJSONObject(arguments) {
		return nil, newRequestValidationError(nil, protocol.ErrInvalidParams, "invalid params")
	}

	if exceedsLimit(arguments, maxArgumentBytes) {
		return nil, newRequestValidationError(nil, protocol.ErrInvalidParams, "invalid params")
	}

	return params, nil
}

func exceedsLimit(raw json.RawMessage, limit int64) bool {
	return limit > 0 && int64(len(bytes.TrimSpace(raw))) > limit
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

	return value, true
}

func requestID(fields map[string]json.RawMessage) (json.RawMessage, bool, bool) {
	raw, ok := fields["id"]
	if !ok {
		return nil, false, true
	}

	if !isValidID(raw) {
		return nil, true, false
	}

	return append(json.RawMessage(nil), raw...), true, true
}

func requestIDOrNull(fields map[string]json.RawMessage) json.RawMessage {
	id, hasID, ok := requestID(fields)
	if !hasID || !ok {
		return nullID()
	}

	return id
}

func idForInvalidRequest(id json.RawMessage, hasID bool) json.RawMessage {
	if !hasID {
		return nullID()
	}

	return id
}

func isValidID(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	if bytes.Equal(trimmed, []byte("null")) {
		return true
	}

	var value any
	if err := json.Unmarshal(trimmed, &value); err != nil {
		return false
	}

	switch value.(type) {
	case string, float64:
		return true
	default:
		return false
	}
}

func isJSONObject(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)

	return len(trimmed) > 0 && trimmed[0] == '{'
}

func isJSONNull(raw json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
}

func newRequestValidationError(id json.RawMessage, code int, message string) *requestValidationError {
	if len(id) == 0 {
		id = nullID()
	}

	return &requestValidationError{
		id:      id,
		code:    code,
		message: message,
	}
}

func nullID() json.RawMessage {
	return json.RawMessage(`null`)
}
