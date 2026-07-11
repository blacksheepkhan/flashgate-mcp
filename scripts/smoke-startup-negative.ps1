$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path "$PSScriptRoot\..").Path
$binaryPath = Join-Path $repoRoot "build\flashgate-mcp.exe"
$stamp = "{0}-{1}" -f ([DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()), $PID
$rootFile = Join-Path $repoRoot "build\startup-root-file-$stamp.txt"
$missingRoot = Join-Path $repoRoot "build\startup-missing-root-$stamp"

if (-not (Test-Path -LiteralPath $binaryPath)) {
    throw "Binary not found: $binaryPath. Run: go build -o build/flashgate-mcp.exe ./cmd/server"
}

function Assert-NoHostPathLeak {
    param([string]$Value)

    if ($Value -match '[A-Za-z]:[\\/]' -or $Value -match '\\\\[^\\]+' -or $Value -match '/(home|etc|private|tmp)/') {
        throw "Startup diagnostics leaked a host path: $Value"
    }
}

function Invoke-StartupCase {
    param(
        [string]$Name,
        [hashtable]$Environment,
        [int]$ExpectedExitCode,
        [string]$ExpectedStderr
    )

    $startInfo = [System.Diagnostics.ProcessStartInfo]::new()
    $startInfo.FileName = $binaryPath
    $startInfo.WorkingDirectory = $repoRoot
    $startInfo.UseShellExecute = $false
    $startInfo.RedirectStandardInput = $true
    $startInfo.RedirectStandardOutput = $true
    $startInfo.RedirectStandardError = $true
    $startInfo.CreateNoWindow = $true

    foreach ($nameToRemove in @("MCP_ROOT", "MCP_ALLOW_CWD_ROOT", "MCP_READ_ONLY")) {
        [void]$startInfo.Environment.Remove($nameToRemove)
    }
    foreach ($entry in $Environment.GetEnumerator()) {
        $startInfo.Environment[$entry.Key] = [string]$entry.Value
    }

    $process = [System.Diagnostics.Process]::new()
    $process.StartInfo = $startInfo
    try {
        if (-not $process.Start()) {
            throw "Failed to start case '$Name'"
        }
        $process.StandardInput.Close()
        $stdout = $process.StandardOutput.ReadToEnd()
        $stderr = $process.StandardError.ReadToEnd()
        $process.WaitForExit()

        if ($process.ExitCode -ne $ExpectedExitCode) {
            throw "Case '$Name': expected exit code $ExpectedExitCode, got $($process.ExitCode). stderr=$stderr"
        }
        if ($stdout.Length -ne 0) {
            throw "Case '$Name': expected empty stdout, got $stdout"
        }

        $normalizedStderr = $stderr -replace "`r`n", "`n"
        if ($normalizedStderr -ne $ExpectedStderr) {
            throw "Case '$Name': expected stderr '$ExpectedStderr', got '$normalizedStderr'"
        }
        Assert-NoHostPathLeak -Value $normalizedStderr
    } finally {
        if (-not $process.HasExited) {
            $process.Kill($true)
            $process.WaitForExit()
        }
        $process.Dispose()
    }
}

New-Item -ItemType Directory -Path (Split-Path $rootFile) -Force | Out-Null
[System.IO.File]::WriteAllText($rootFile, "not a directory", [System.Text.UTF8Encoding]::new($false))

try {
    $cases = @(
        @{ Name = "missing root"; Environment = @{}; Exit = 3; Stderr = "flashgate-mcp: startup failed (missing_root)`n" }
        @{ Name = "empty root"; Environment = @{ MCP_ROOT = "" }; Exit = 3; Stderr = "flashgate-mcp: startup failed (invalid_root)`n" }
        @{ Name = "whitespace root"; Environment = @{ MCP_ROOT = "   " }; Exit = 3; Stderr = "flashgate-mcp: startup failed (invalid_root)`n" }
        @{ Name = "relative root"; Environment = @{ MCP_ROOT = "subdir" }; Exit = 3; Stderr = "flashgate-mcp: startup failed (invalid_root)`n" }
        @{ Name = "parent root"; Environment = @{ MCP_ROOT = ".." }; Exit = 3; Stderr = "flashgate-mcp: startup failed (invalid_root)`n" }
        @{ Name = "dot without opt-in"; Environment = @{ MCP_ROOT = "." }; Exit = 3; Stderr = "flashgate-mcp: startup failed (invalid_root)`n" }
        @{ Name = "opt-in without root"; Environment = @{ MCP_ALLOW_CWD_ROOT = "true" }; Exit = 3; Stderr = "flashgate-mcp: startup failed (missing_root)`n" }
        @{ Name = "invalid development option"; Environment = @{ MCP_ROOT = $repoRoot; MCP_ALLOW_CWD_ROOT = "TRUE" }; Exit = 3; Stderr = "flashgate-mcp: startup failed (invalid_development_option)`n" }
        @{ Name = "missing filesystem root"; Environment = @{ MCP_ROOT = $missingRoot }; Exit = 3; Stderr = "flashgate-mcp: startup failed (root_not_found)`n" }
        @{ Name = "file root"; Environment = @{ MCP_ROOT = $rootFile }; Exit = 3; Stderr = "flashgate-mcp: startup failed (root_not_directory)`n" }
        @{ Name = "invalid read-only profile"; Environment = @{ MCP_ROOT = $repoRoot; MCP_READ_ONLY = "invalid" }; Exit = 3; Stderr = "flashgate-mcp: startup failed (invalid_profile)`n" }
        @{ Name = "valid absolute root"; Environment = @{ MCP_ROOT = $repoRoot; MCP_READ_ONLY = "true" }; Exit = 0; Stderr = "" }
        @{ Name = "development CWD"; Environment = @{ MCP_ROOT = "."; MCP_ALLOW_CWD_ROOT = "true"; MCP_READ_ONLY = "true" }; Exit = 0; Stderr = "flashgate-mcp: warning: development CWD root enabled`n" }
    )

    foreach ($case in $cases) {
        Invoke-StartupCase -Name $case.Name -Environment $case.Environment -ExpectedExitCode $case.Exit -ExpectedStderr $case.Stderr
    }

    Write-Host "Startup negative smoke test passed."
} finally {
    Remove-Item -LiteralPath $rootFile, $missingRoot -Force -Recurse -ErrorAction SilentlyContinue
}
