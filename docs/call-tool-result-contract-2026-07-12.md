# CallToolResult contract migration

Date: 2026-07-12
Sprint: 3.45a – MCP CallToolResult foundation

## Reason

A direct STDIO preflight accepted syntactically valid JSON-RPC and correct domain values, but the first Codex read-only Canary failed with `Unexpected response type`. Successful tools returned the domain object directly as JSON-RPC `result`, which is not an MCP 2025-11-25 `CallToolResult`.

Previous form:

```json
{"jsonrpc":"2.0","id":3,"result":{"entries":[]}}
```

## New successful wire form

All eight filesystem tools now cross one central adapter wrapper:

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "content": [
      { "type": "text", "text": "{\"entries\":[]}" }
    ],
    "structuredContent": {
      "entries": []
    }
  }
}
```

The text is compact deterministic JSON produced with Go's standard `encoding/json`. The same serialized bytes are used as `structuredContent`, so both representations decode to deeply equal objects. No domain field is placed directly on the outer result, and successful calls omit `isError`.

## `read_file` content fields

`read_file` retains its domain object:

```json
{"content":"file text","size":9}
```

The outer MCP `content` is an array; the inner `structuredContent.content` remains a string:

```json
{
  "content": [{"type":"text","text":"{\"content\":\"file text\",\"size\":9}"}],
  "structuredContent": {"content":"file text","size":9}
}
```

## Missing semantics

A genuinely missing `get_path_info` path remains a successful call with exactly the domain fields `path` and `exists:false` in both representations. It is not a JSON-RPC error and does not set `isError=true`.

## Client and smoke impact

Clients must read domain values from `structuredContent` or decode the JSON text block. Windows and Bash smokes now validate the outer envelope, text-block shape, compact JSON, semantic parity, Missing behavior, and existing domain assertions. A strict local decoder preserves the old unwrapped forms as negative regression fixtures.

No tool name, discovery order, input schema, argument, filesystem behavior, root policy, or capability gate changes. Existing safe JSON-RPC tool errors remain unchanged; the general `isError=true` migration is deferred to BL-203.

The catalog's `resultSchema` values continue to describe domain objects. Runtime `outputSchema` is not advertised and remains the next separate all-eight-tool Sprint 3.45 gate.

## Canary gate

Sprint 3.45a does not reactivate the Canary and does not build a new Canary binary. Reuse requires review, commit, PR, green Windows/Ubuntu CI, merge, post-merge gates, a new versioned binary, strict direct preflight, separate user approval, full Codex restart, and a successful model-driven end-to-end test.
