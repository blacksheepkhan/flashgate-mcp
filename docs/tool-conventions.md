# MCP tool conventions

## Naming and exposure

Tool names use stable lower-case snake case and describe the user-visible operation. The Sprint 3.43 baseline is:

```text
list_directory
read_file
get_path_info
write_file
create_directory
delete_path
copy_path
move_path
```

The registry determines deterministic exposure order. `tools/list` uses each implementation's single `Definition()` value for name, title, description, and input schema. MCP names remain in the adapter layer; the filesystem core keeps domain-oriented Go names such as `List`, `Stat`, and `Mkdir`.

## Arguments

Every tool accepts exactly one JSON object. Runtime decoding uses `json.Decoder`, rejects unknown fields and explicit `null` field values, and requires EOF after the first object. Schemas declare `additionalProperties:false`.

Required path fields must be strings and must not be empty or whitespace-only. Validation does not trim or otherwise alter valid path strings. `list_directory.path` is the only optional path; omission defaults to `.`, but an explicit blank value is invalid.

Security and PathGuard policy are enforced server-side and are not delegated to JSON Schema.

## Definitions and schemas

Each tool owns a compact definition containing:

- name;
- human-readable title;
- model-useful description;
- closed input schema with explicit required fields.

`docs/mcp-tool-catalog.json` is the static contract view. A focused contract test compares runtime and catalog name, title, description, and the complete input schema without introducing a general schema engine.

Result schemas remain documentation in Sprint 3.43. Runtime `outputSchema`, `structuredContent`, and a general JSON Schema 2020-12 migration remain later work.

## Path and result conventions

All client paths are relative to the configured root. Results may echo the public relative path supplied by the client; they must never expose the PathGuard-resolved absolute host path.

`get_path_info` reports genuine absence as a successful `{path, exists:false}` result. Existing paths include `name`, `isDir`, and `size`. Policy denials are never converted to absence.

`create_directory.created` is `true` only when the leaf directory was created by that call and `false` for an existing directory.

`copy_path` is file-only. `move_path` is the sole Move/Rename contract and is same-volume only.

## Capability gating

The default profile exposes all eight baseline tools. The read-only profile exposes exactly:

```text
list_directory
read_file
get_path_info
```

`write_file`, `create_directory`, `delete_path`, `copy_path`, and `move_path` are write-gated and absent from the read-only registry.

Client activation must set `MCP_READ_ONLY=true` explicitly; the missing-variable default remains the eight-tool profile. The read-only and negative STDIO smokes require identical generic Invalid params responses for every write-gated and removed legacy name.

## Errors

The adapter retains the existing JSON-RPC architecture:

- parse error `-32700`;
- invalid request `-32600`;
- method not found `-32601`;
- expected argument, path, policy, and filesystem contract failures `-32602`;
- unexpected I/O failures `-32603`.

Internal classification uses `not_found`, `already_exists`, `access_denied`, `invalid_path`, `unsupported_path_type`, `unsupported_operation`, `limit_exceeded`, and `io_error`. Messages are safe and generic. Stable wire-level error objects are deferred.

## Compatibility

FlashGate MCP is pre-1.0 and was not productively deployed when Sprint 3.43 cleaned the contract. Removed names have no alias or deprecation layer. Clients must update their calls and discovery expectations as described in the dated migration.
