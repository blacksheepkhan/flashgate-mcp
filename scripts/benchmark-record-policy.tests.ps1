[CmdletBinding()]
param()

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$Failures = [System.Collections.Generic.List[string]]::new()
$CheckCount = 0
$TemporaryDirectory = $null
$OriginalPath = $env:PATH

function Test-Condition {
    param([bool]$Condition, [string]$Name)
    $script:CheckCount++
    if (-not $Condition) { $script:Failures.Add($Name) }
}

function Invoke-TestWrapper {
    param([string]$WrapperPath, [string[]]$Arguments)

    $ProcessInfo = [Diagnostics.ProcessStartInfo]::new()
    $ProcessInfo.FileName = (Get-Process -Id $PID).Path
    $ProcessInfo.UseShellExecute = $false
    $ProcessInfo.RedirectStandardOutput = $true
    $ProcessInfo.RedirectStandardError = $true
    foreach ($Argument in @('-NoProfile', '-File', $WrapperPath) + $Arguments) {
        [void]$ProcessInfo.ArgumentList.Add($Argument)
    }
    $Process = [Diagnostics.Process]::Start($ProcessInfo)
    $StandardOutput = $Process.StandardOutput.ReadToEnd()
    $StandardError = $Process.StandardError.ReadToEnd()
    $Process.WaitForExit()
    return [pscustomobject]@{ ExitCode = $Process.ExitCode; StandardOutput = $StandardOutput; StandardError = $StandardError }
}

try {
    $TemporaryDirectory = Join-Path ([IO.Path]::GetTempPath()) ('flashgate-record-policy-test-' + [guid]::NewGuid().ToString('N'))
    $FakeRepository = Join-Path $TemporaryDirectory 'repo'
    $FakeScripts = Join-Path $FakeRepository 'scripts'
    $StubDirectory = Join-Path $TemporaryDirectory 'stub'
    New-Item -ItemType Directory -Path $FakeScripts, $StubDirectory | Out-Null

    $Wrapper = Join-Path $FakeScripts 'benchmark.ps1'
    Copy-Item -LiteralPath (Join-Path $PSScriptRoot 'benchmark.ps1') -Destination $Wrapper
    Copy-Item -LiteralPath (Join-Path $PSScriptRoot 'benchmark-window.ps1') -Destination (Join-Path $FakeScripts 'benchmark-window.ps1')
    $Marker = Join-Path $TemporaryDirectory 'go-invoked.marker'
    $Output = Join-Path $TemporaryDirectory 'baseline.windows-amd64.json'
    $Stub = Join-Path $StubDirectory 'go.cmd'
    @(
        '@echo off',
        ('type nul > "{0}"' -f $Marker),
        'exit /b 97'
    ) | Set-Content -LiteralPath $Stub -Encoding ascii
    $env:PATH = $StubDirectory + [IO.Path]::PathSeparator + $OriginalPath

    $RecordResult = Invoke-TestWrapper -WrapperPath $Wrapper -Arguments @('-RecordBaseline', '-OutputPath', $Output)

    Test-Condition ($RecordResult.ExitCode -ne 0) 'record mode exits nonzero'
    Test-Condition ($RecordResult.StandardError -match 'Authoritative baseline recording is not supported by scripts/benchmark.ps1') 'record mode emits policy error on stderr'
    Test-Condition (-not (Test-Path -LiteralPath $Marker)) 'record mode does not invoke go'
    Test-Condition (-not (Test-Path -LiteralPath $Output)) 'record mode does not write baseline output'
    Test-Condition (-not (Test-Path -LiteralPath (Join-Path $FakeRepository 'build'))) 'record mode does not create build directory'
    Test-Condition ([string]::IsNullOrWhiteSpace($RecordResult.StandardOutput)) 'record mode does not emit success output'

    $CanonicalOutput = Join-Path $FakeRepository 'benchmarks\baseline.windows-amd64.json'
    $CanonicalResult = Invoke-TestWrapper -WrapperPath $Wrapper -Arguments @('-OutputPath', $CanonicalOutput)
    Test-Condition ($CanonicalResult.ExitCode -ne 0) 'non-authoritative mode rejects canonical baseline path'
    Test-Condition ($CanonicalResult.StandardOutput -match 'Non-authoritative benchmark runs must not write a canonical versioned baseline path') 'canonical-path rejection is reported'
    Test-Condition (-not (Test-Path -LiteralPath $Marker)) 'canonical-path rejection does not invoke go'
    Test-Condition (-not (Test-Path -LiteralPath $CanonicalOutput)) 'canonical-path rejection does not write output'
    Test-Condition (-not (Test-Path -LiteralPath (Join-Path $FakeRepository 'build'))) 'canonical-path rejection does not create build directory'
}
catch {
    $Failures.Add("unexpected test error: $($_.Exception.Message)")
}
finally {
    $env:PATH = $OriginalPath
    if ($null -ne $TemporaryDirectory -and (Test-Path -LiteralPath $TemporaryDirectory -PathType Container)) {
        Remove-Item -LiteralPath $TemporaryDirectory -Recurse -Force
    }
    [pscustomobject]@{
        Status       = $(if ($Failures.Count -eq 0) { 'PASS' } else { 'FAIL' })
        CheckCount   = $CheckCount
        FailureCount = $Failures.Count
        Failures     = ($Failures -join '; ')
    } | Format-List
}

if ($Failures.Count -gt 0) { exit 1 }
