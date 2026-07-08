$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path "$PSScriptRoot\.."
$binaryPath = Join-Path $repoRoot "build\fileserver-mcp.exe"
$requestPath = Join-Path $repoRoot "build\smoke-jsonrpc-request.jsonl"
$responsePath = Join-Path $repoRoot "build\smoke-jsonrpc-response.jsonl"

if (-not (Test-Path $binaryPath)) {
    throw "Binary not found: $binaryPath. Run: go build -o build/fileserver-mcp.exe ./cmd/server"
}

$env:MCP_ROOT = $repoRoot

$requests = @(
    '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"smoke-test","version":"dev"}}}'
    '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
)

[System.IO.File]::WriteAllText($requestPath, (($requests -join "`n") + "`n"), [System.Text.UTF8Encoding]::new($false))

cmd /c "`"$binaryPath`" < `"$requestPath`" > `"$responsePath`""

if ($LASTEXITCODE -ne 0) {
    throw "Smoke test process failed with exit code $LASTEXITCODE"
}

$responses = Get-Content $responsePath | Where-Object { $_.Trim().Length -gt 0 }

if ($responses.Count -ne 2) {
    throw "Expected 2 JSON-RPC responses, got $($responses.Count). Response file: $responsePath"
}

$initialize = $responses[0] | ConvertFrom-Json
$toolsList = $responses[1] | ConvertFrom-Json

if ($initialize.id -ne 1) {
    throw "Expected initialize response id 1, got $($initialize.id)"
}

if (-not $initialize.result.protocolVersion) {
    throw "Initialize response does not contain protocolVersion"
}

if ($toolsList.id -ne 2) {
    throw "Expected tools/list response id 2, got $($toolsList.id)"
}

if (-not $toolsList.result.tools) {
    throw "tools/list response does not contain tools"
}

$toolNames = @($toolsList.result.tools | ForEach-Object { $_.name })

$expectedTools = @(
    "list_files",
    "read_file",
    "stat_path",
    "exists_path"
)

if ($env:MCP_READ_ONLY -ne "true") {
    $expectedTools += @(
        "write_file",
        "mkdir",
        "delete_path",
        "move_path",
        "copy_path",
        "rename_path"
    )
}

foreach ($expectedTool in $expectedTools) {
    if ($toolNames -notcontains $expectedTool) {
        throw "Expected tool '$expectedTool' was not listed. Actual tools: $($toolNames -join ', ')"
    }
}

foreach ($toolName in $toolNames) {
    if ($expectedTools -notcontains $toolName) {
        throw "Unexpected tool '$toolName' was listed. Expected tools: $($expectedTools -join ', ')"
    }
}

Write-Host "JSON-RPC smoke test passed."
Write-Host "Protocol version: $($initialize.result.protocolVersion)"
Write-Host "Tools: $($toolNames -join ', ')"
