$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path "$PSScriptRoot\.."
$binaryPath = Join-Path $repoRoot "build\flashgate-mcp.exe"
$buildDir = Join-Path $repoRoot "build"
$stamp = "{0}-{1}" -f ([DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()), $PID
$requestPath = Join-Path $buildDir "smoke-jsonrpc-$stamp-request.jsonl"
$responsePath = Join-Path $buildDir "smoke-jsonrpc-$stamp-response.jsonl"
$moveSourceRelative = "build/smoke-move-$stamp-source.txt"
$moveTargetRelative = "build/smoke-move-$stamp-target.txt"
$moveSourcePath = Join-Path $repoRoot $moveSourceRelative
$moveTargetPath = Join-Path $repoRoot $moveTargetRelative

if (-not (Test-Path $binaryPath)) {
    throw "Binary not found: $binaryPath. Run: go build -o build/flashgate-mcp.exe ./cmd/server"
}

$env:MCP_ROOT = $repoRoot

New-Item -ItemType Directory -Path $buildDir -Force | Out-Null

try {
    $requests = @(
        '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"smoke-test","version":"dev"}}}'
        '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
        '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"list_directory","arguments":{}}}'
        '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"get_path_info","arguments":{"path":"README.md"}}}'
        '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"get_path_info","arguments":{"path":"smoke-missing-path"}}}'
    )

    if ($env:MCP_READ_ONLY -ne "true") {
        [System.IO.File]::WriteAllText($moveSourcePath, "move-smoke", [System.Text.UTF8Encoding]::new($false))
        $requests += '{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"move_path","arguments":{"source":"' + $moveSourceRelative + '","target":"' + $moveTargetRelative + '"}}}'
    }

    [System.IO.File]::WriteAllText($requestPath, (($requests -join "`n") + "`n"), [System.Text.UTF8Encoding]::new($false))

    cmd /c "`"$binaryPath`" < `"$requestPath`" > `"$responsePath`""

    if ($LASTEXITCODE -ne 0) {
        throw "Smoke test process failed with exit code $LASTEXITCODE"
    }

    $responses = Get-Content $responsePath | Where-Object { $_.Trim().Length -gt 0 }

    $expectedResponseCount = if ($env:MCP_READ_ONLY -eq "true") { 5 } else { 6 }
    if ($responses.Count -ne $expectedResponseCount) {
        throw "Expected $expectedResponseCount JSON-RPC responses, got $($responses.Count). Response file: $responsePath"
    }

    $initialize = $responses[0] | ConvertFrom-Json
    $toolsList = $responses[1] | ConvertFrom-Json
    $listDirectory = $responses[2] | ConvertFrom-Json
    $existingPathInfo = $responses[3] | ConvertFrom-Json
    $missingPathInfo = $responses[4] | ConvertFrom-Json

    if ($initialize.id -ne 1) {
        throw "Expected initialize response id 1, got $($initialize.id)"
    }

    if (-not $initialize.result.protocolVersion) {
        throw "Initialize response does not contain protocolVersion"
    }

    if ($initialize.result.serverInfo.name -ne "flashgate") {
        throw "Expected serverInfo.name flashgate, got $($initialize.result.serverInfo.name)"
    }

    if ($toolsList.id -ne 2) {
        throw "Expected tools/list response id 2, got $($toolsList.id)"
    }

    if (-not $toolsList.result.tools) {
        throw "tools/list response does not contain tools"
    }

    $toolNames = @($toolsList.result.tools | ForEach-Object { $_.name })

    $expectedTools = @(
        "list_directory",
        "read_file",
        "get_path_info"
    )

    if ($env:MCP_READ_ONLY -ne "true") {
        $expectedTools += @(
            "write_file",
            "create_directory",
            "delete_path",
            "copy_path",
            "move_path"
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

    if ($listDirectory.id -ne 3 -or $null -eq $listDirectory.result.entries) {
        throw "list_directory did not return entries"
    }
    if ($existingPathInfo.id -ne 4 -or -not $existingPathInfo.result.exists -or $existingPathInfo.result.path -ne "README.md") {
        throw "get_path_info did not report README.md as existing"
    }
    if ($missingPathInfo.id -ne 5 -or $missingPathInfo.result.exists -or $missingPathInfo.result.path -ne "smoke-missing-path") {
        throw "get_path_info did not report the missing path correctly"
    }
    if ($env:MCP_READ_ONLY -ne "true") {
        $moveResult = $responses[5] | ConvertFrom-Json
        if ($moveResult.id -ne 6 -or -not $moveResult.result.moved -or -not (Test-Path -LiteralPath $moveTargetPath)) {
            throw "move_path did not perform rename semantics"
        }
    }

    Write-Host "JSON-RPC smoke test passed."
    Write-Host "Protocol version: $($initialize.result.protocolVersion)"
    Write-Host "Tools: $($toolNames -join ', ')"
} finally {
    Remove-Item -LiteralPath $requestPath, $responsePath, $moveSourcePath, $moveTargetPath -Force -ErrorAction SilentlyContinue
}
