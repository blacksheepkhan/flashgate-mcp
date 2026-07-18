[CmdletBinding()]
param(
    [switch]$Quick,
    [switch]$RecordBaseline,
    [string]$OutputPath
)

if ($RecordBaseline) {
    throw 'Authoritative baseline recording is not supported by scripts/benchmark.ps1. Use the documented two-phase prebuilt workflow in the local benchmark workspace C:\Voxtronic\Codex\Temp\Benchmarks.'
}

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $true
. (Join-Path $PSScriptRoot 'benchmark-window.ps1')

$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$BuildDirectory = Join-Path $RepoRoot 'build'
$ServerBinary = Join-Path $BuildDirectory 'flashgate-mcp.exe'
$BenchmarkBinary = Join-Path $BuildDirectory 'flashgate-benchmark.exe'
$BudgetPath = Join-Path $RepoRoot 'benchmarks\budgets.json'
$Status = 'FAIL'
$WarningCount = 0
$FailureCount = 0
$NextAction = 'Inspect the reported failure.'
$ReportPath = $null
$FailureMessage = $null
$RunOutputPath = $null
$PerformanceContaminated = $false
$MeasurementWarning = $null

try {
    if ([string]::IsNullOrWhiteSpace($OutputPath)) {
        $OutputPath = Join-Path $BuildDirectory 'benchmark-current.windows-amd64.json'
    }
    elseif (-not [IO.Path]::IsPathRooted($OutputPath)) {
        $OutputPath = Join-Path $RepoRoot $OutputPath
    }

    $OutputFullPath = [IO.Path]::GetFullPath($OutputPath)
    $BenchmarkDirectory = [IO.Path]::GetFullPath((Join-Path $RepoRoot 'benchmarks'))
    if ([IO.Path]::GetDirectoryName($OutputFullPath).Equals($BenchmarkDirectory, [StringComparison]::OrdinalIgnoreCase) -and [IO.Path]::GetFileName($OutputFullPath) -like 'baseline.*-*.json') {
        throw 'Non-authoritative benchmark runs must not write a canonical versioned baseline path.'
    }
    $OutputPath = $OutputFullPath

    $PerformanceContaminated = (Get-EuropeViennaMeasurementWindowStatus).IsBlocked

    $InitialStatus = @(& git status --porcelain --untracked-files=all)
    $WorkingTreeDirty = $InitialStatus.Count -gt 0

    New-Item -ItemType Directory -Path $BuildDirectory -Force | Out-Null
    & go build -o $ServerBinary ./cmd/server 2>&1 | Out-Null
    & go build -o $BenchmarkBinary ./cmd/benchmark 2>&1 | Out-Null
    $Commit = (& git rev-parse HEAD).Trim()
    $RunOutputPath = $OutputPath

    $Arguments = @(
        '-binary', $ServerBinary,
        '-output', $RunOutputPath,
        '-commit', $Commit,
        '-budgets', $BudgetPath
    )
    if ($WorkingTreeDirty) {
        $Arguments += '-working-tree-dirty'
    }
    if ($Quick) {
        $Arguments += '-quick'
    }

    $BenchmarkOutput = @(& $BenchmarkBinary @Arguments 2>&1)
    $Result = Get-Content -LiteralPath $RunOutputPath -Raw | ConvertFrom-Json
    $WarningCount = @($Result.warnings).Count
    $FailureCount = [int]$Result.budget_evaluation.hard_failures
    if ($FailureCount -gt 0) {
        throw "Benchmark reported $FailureCount hard budget failure(s)."
    }

    if ((Get-EuropeViennaMeasurementWindowStatus).IsBlocked) {
        $PerformanceContaminated = $true
    }

    if ($PerformanceContaminated) {
        $MeasurementWarning = Get-ContaminatedPerformanceWarning
        $WarningCount++
    }

    $Status = 'PASS'
    $NextAction = 'Compare the current result with the versioned baseline.'
}
catch {
    $FailureCount = [Math]::Max(1, $FailureCount)
    $FailureMessage = $_.Exception.Message
}
finally {
    $ReportPath = if ([string]::IsNullOrWhiteSpace($OutputPath)) { '[not-created]' } else { $OutputPath }
    [pscustomobject]@{
        Status         = $Status
        ReportPath     = $ReportPath
        WarningCount   = $WarningCount
        FailureCount   = $FailureCount
        NextAction     = $NextAction
        MeasurementWarning = $MeasurementWarning
        FailureMessage = $FailureMessage
    } | Format-List
}

if ($FailureCount -gt 0) {
    exit 1
}
