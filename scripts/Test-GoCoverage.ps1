[CmdletBinding()]
param(
    [Parameter(Mandatory)]
    [ValidateSet('windows', 'linux')]
    [string]$PlatformName,

    [Parameter()]
    [ValidateRange(0.0, 100.0)]
    [double]$MinimumCoverage = 0.0,

    [Parameter()]
    [ValidateNotNullOrEmpty()]
    [string]$OutputRoot = 'build/coverage'
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $false

$status = 'FAIL'
$warningCount = 0
$failureCount = 1
$exitCode = 1
$totalCoverage = $null
$errorMessage = $null

function Resolve-OutputPath {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)]
        [string]$Path
    )

    if ([IO.Path]::IsPathFullyQualified($Path)) {
        return [IO.Path]::GetFullPath($Path)
    }

    return [IO.Path]::GetFullPath(
        (Join-Path -Path (Get-Location) -ChildPath $Path)
    )
}

function Get-OutputTail {
    [CmdletBinding()]
    param(
        [Parameter()]
        [object[]]$Output,

        [Parameter()]
        [ValidateRange(1, 100)]
        [int]$LineCount = 25
    )

    $lines = @($Output | ForEach-Object { $_.ToString() })
    if ($lines.Count -eq 0) {
        return '<keine Ausgabe>'
    }

    return (($lines | Select-Object -Last $LineCount) -join [Environment]::NewLine)
}

try {
    $platformOutputRoot = Resolve-OutputPath -Path (
        Join-Path -Path $OutputRoot -ChildPath $PlatformName
    )

    [IO.Directory]::CreateDirectory($platformOutputRoot) | Out-Null

    $coverageProfile = Join-Path $platformOutputRoot 'coverage.out'
    $textReport = Join-Path $platformOutputRoot 'coverage.txt'
    $htmlReport = Join-Path $platformOutputRoot 'coverage.html'
    $testLog = Join-Path $platformOutputRoot 'test.log'
    $summaryJson = Join-Path $platformOutputRoot 'summary.json'

    foreach ($file in @(
        $coverageProfile,
        $textReport,
        $htmlReport,
        $testLog,
        $summaryJson
    )) {
        if ([IO.File]::Exists($file)) {
            [IO.File]::Delete($file)
        }
    }

    $testOutput = @(
        & go test `
            -covermode=atomic `
            -coverpkg=./... `
            "-coverprofile=$coverageProfile" `
            ./... 2>&1
    )
    $testExitCode = $LASTEXITCODE

    [IO.File]::WriteAllText(
        $testLog,
        (($testOutput | ForEach-Object { $_.ToString() }) -join [Environment]::NewLine) +
            [Environment]::NewLine,
        [Text.UTF8Encoding]::new($false)
    )

    if ($testExitCode -ne 0) {
        $tail = Get-OutputTail -Output $testOutput
        throw "go test ist mit Exitcode $testExitCode fehlgeschlagen.$([Environment]::NewLine)$tail"
    }

    $coverageOutput = @(
        & go tool cover "-func=$coverageProfile" 2>&1
    )
    $coverageExitCode = $LASTEXITCODE

    if ($coverageExitCode -ne 0) {
        $tail = Get-OutputTail -Output $coverageOutput
        throw "go tool cover -func ist mit Exitcode $coverageExitCode fehlgeschlagen.$([Environment]::NewLine)$tail"
    }

    $coverageText = (
        $coverageOutput |
            ForEach-Object { $_.ToString() }
    ) -join [Environment]::NewLine

    [IO.File]::WriteAllText(
        $textReport,
        $coverageText + [Environment]::NewLine,
        [Text.UTF8Encoding]::new($false)
    )

    $totalLine = $coverageOutput |
        ForEach-Object { $_.ToString() } |
        Where-Object {
            $_ -match '^total:\s+\(statements\)\s+\d+(?:\.\d+)?%$'
        } |
        Select-Object -Last 1

    if (-not $totalLine) {
        throw 'Die Gesamt-Coverage konnte nicht aus go tool cover ermittelt werden.'
    }

    $coverageMatch = [regex]::Match(
        $totalLine,
        '^total:\s+\(statements\)\s+(?<coverage>\d+(?:\.\d+)?)%$'
    )

    if (-not $coverageMatch.Success) {
        throw "Die Gesamt-Coverage konnte nicht geparst werden: $totalLine"
    }

    $totalCoverage = [double]::Parse(
        $coverageMatch.Groups['coverage'].Value,
        [Globalization.CultureInfo]::InvariantCulture
    )

    $htmlOutput = @(
        & go tool cover "-html=$coverageProfile" "-o=$htmlReport" 2>&1
    )
    $htmlExitCode = $LASTEXITCODE

    if ($htmlExitCode -ne 0) {
        $tail = Get-OutputTail -Output $htmlOutput
        throw "go tool cover -html ist mit Exitcode $htmlExitCode fehlgeschlagen.$([Environment]::NewLine)$tail"
    }

    $summary = [ordered]@{
        status = 'PASS'
        platform = $PlatformName
        total_coverage_percent = $totalCoverage
        minimum_coverage_percent = $MinimumCoverage
        coverage_profile = $coverageProfile
        text_report = $textReport
        html_report = $htmlReport
        test_log = $testLog
    }

    [IO.File]::WriteAllText(
        $summaryJson,
        ($summary | ConvertTo-Json -Depth 4) + [Environment]::NewLine,
        [Text.UTF8Encoding]::new($false)
    )

    if ($env:GITHUB_STEP_SUMMARY) {
        $summaryMarkdown = @"
## Go Code Coverage — $PlatformName

| Kennzahl | Wert |
|---|---:|
| Gesamt-Coverage | $($totalCoverage.ToString('0.0', [Globalization.CultureInfo]::InvariantCulture)) % |
| Mindestwert | $($MinimumCoverage.ToString('0.0', [Globalization.CultureInfo]::InvariantCulture)) % |
"@

        [IO.File]::AppendAllText(
            $env:GITHUB_STEP_SUMMARY,
            $summaryMarkdown + [Environment]::NewLine,
            [Text.UTF8Encoding]::new($false)
        )
    }

    if ($totalCoverage -lt $MinimumCoverage) {
        throw (
            'Coverage {0:0.0} % unterschreitet den Mindestwert {1:0.0} %.' -f
            $totalCoverage,
            $MinimumCoverage
        )
    }

    $status = 'PASS'
    $failureCount = 0
    $exitCode = 0
}
catch {
    $errorMessage = $_.Exception.Message
}
finally {
    [pscustomobject]@{
        Status          = $status
        Platform        = $PlatformName
        TotalCoverage   = if ($null -eq $totalCoverage) {
            $null
        }
        else {
            '{0:0.0} %' -f $totalCoverage
        }
        MinimumCoverage = '{0:0.0} %' -f $MinimumCoverage
        ReportPath      = Join-Path $OutputRoot $PlatformName
        WarningCount    = $warningCount
        FailureCount    = $failureCount
        NextAction      = if ($status -eq 'PASS') {
            'Baseline dokumentieren und den plattformspezifischen Mindestwert festlegen.'
        }
        else {
            'Fehler beheben und den Coverage-Lauf erneut ausführen.'
        }
        Error           = $errorMessage
    } | Format-List
}

exit $exitCode
