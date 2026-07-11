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
$readOnlyWriteRelative = "build/smoke-readonly-$stamp-write.txt"
$readOnlyCreateRelative = "build/smoke-readonly-$stamp-directory"
$readOnlyDeleteRelative = "build/smoke-readonly-$stamp-delete.txt"
$readOnlyCopySourceRelative = "build/smoke-readonly-$stamp-copy-source.txt"
$readOnlyCopyTargetRelative = "build/smoke-readonly-$stamp-copy-target.txt"
$readOnlyMoveSourceRelative = "build/smoke-readonly-$stamp-move-source.txt"
$readOnlyMoveTargetRelative = "build/smoke-readonly-$stamp-move-target.txt"
$readOnlyWritePath = Join-Path $repoRoot $readOnlyWriteRelative
$readOnlyCreatePath = Join-Path $repoRoot $readOnlyCreateRelative
$readOnlyDeletePath = Join-Path $repoRoot $readOnlyDeleteRelative
$readOnlyCopySourcePath = Join-Path $repoRoot $readOnlyCopySourceRelative
$readOnlyCopyTargetPath = Join-Path $repoRoot $readOnlyCopyTargetRelative
$readOnlyMoveSourcePath = Join-Path $repoRoot $readOnlyMoveSourceRelative
$readOnlyMoveTargetPath = Join-Path $repoRoot $readOnlyMoveTargetRelative

if (-not (Test-Path $binaryPath)) {
    throw "Binary not found: $binaryPath. Run: go build -o build/flashgate-mcp.exe ./cmd/server"
}

$env:MCP_ROOT = $repoRoot

New-Item -ItemType Directory -Path $buildDir -Force | Out-Null

try {
    if ($env:MCP_READ_ONLY -eq "true") {
        foreach ($fixturePath in @($readOnlyDeletePath, $readOnlyCopySourcePath, $readOnlyMoveSourcePath)) {
            [System.IO.File]::WriteAllText($fixturePath, "read-only-smoke-fixture", [System.Text.UTF8Encoding]::new($false))
        }
    }

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
    } else {
        $requests += @(
            ('{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"write_file","arguments":{"path":"' + $readOnlyWriteRelative + '","content":"blocked"}}}')
            ('{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"create_directory","arguments":{"path":"' + $readOnlyCreateRelative + '"}}}')
            ('{"jsonrpc":"2.0","id":8,"method":"tools/call","params":{"name":"delete_path","arguments":{"path":"' + $readOnlyDeleteRelative + '"}}}')
            ('{"jsonrpc":"2.0","id":9,"method":"tools/call","params":{"name":"copy_path","arguments":{"source":"' + $readOnlyCopySourceRelative + '","target":"' + $readOnlyCopyTargetRelative + '"}}}')
            ('{"jsonrpc":"2.0","id":10,"method":"tools/call","params":{"name":"move_path","arguments":{"source":"' + $readOnlyMoveSourceRelative + '","target":"' + $readOnlyMoveTargetRelative + '"}}}')
        )
    }

    [System.IO.File]::WriteAllText($requestPath, (($requests -join "`n") + "`n"), [System.Text.UTF8Encoding]::new($false))

    cmd /c "`"$binaryPath`" < `"$requestPath`" > `"$responsePath`""

    if ($LASTEXITCODE -ne 0) {
        throw "Smoke test process failed with exit code $LASTEXITCODE"
    }

    $responses = Get-Content $responsePath | Where-Object { $_.Trim().Length -gt 0 }

    $expectedResponseCount = if ($env:MCP_READ_ONLY -eq "true") { 10 } else { 6 }
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

    if (($toolNames -join "`0") -cne ($expectedTools -join "`0")) {
        throw "Unexpected tool order. Expected: $($expectedTools -join ', '). Actual: $($toolNames -join ', ')"
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
    } else {
        $writeToolResponses = @($responses[5..9] | ForEach-Object { $_ | ConvertFrom-Json })
        for ($index = 0; $index -lt $writeToolResponses.Count; $index++) {
            $response = $writeToolResponses[$index]
            $expectedId = $index + 6
            if ($response.id -ne $expectedId -or $response.error.code -ne -32602 -or $response.error.message -ne "invalid params") {
                throw "Expected read-only-gated write tool id $expectedId to return generic Invalid params"
            }
        }
        foreach ($expectedFixture in @($readOnlyDeletePath, $readOnlyCopySourcePath, $readOnlyMoveSourcePath)) {
            if (-not (Test-Path -LiteralPath $expectedFixture -PathType Leaf)) {
                throw "Read-only smoke fixture was modified or removed: $expectedFixture"
            }
        }
        foreach ($unexpectedPath in @($readOnlyWritePath, $readOnlyCreatePath, $readOnlyCopyTargetPath, $readOnlyMoveTargetPath)) {
            if (Test-Path -LiteralPath $unexpectedPath) {
                throw "Read-only smoke created an unexpected path: $unexpectedPath"
            }
        }
    }

    Write-Host "JSON-RPC smoke test passed."
    Write-Host "Protocol version: $($initialize.result.protocolVersion)"
    Write-Host "Tools: $($toolNames -join ', ')"
} finally {
    Remove-Item -LiteralPath $requestPath, $responsePath, $moveSourcePath, $moveTargetPath, $readOnlyWritePath, $readOnlyDeletePath, $readOnlyCopySourcePath, $readOnlyCopyTargetPath, $readOnlyMoveSourcePath, $readOnlyMoveTargetPath -Force -ErrorAction SilentlyContinue
    Remove-Item -LiteralPath $readOnlyCreatePath -Recurse -Force -ErrorAction SilentlyContinue
}
