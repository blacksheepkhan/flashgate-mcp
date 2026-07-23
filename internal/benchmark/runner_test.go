package benchmark

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestBenchmarkEnvironmentIsMinimalAndDeterministic(t *testing.T) {
	t.Setenv("MCP_MAX_LIST_ENTRIES", "1")
	t.Setenv("FLASHGATE_PRIVATE_TEST_VALUE", "must-not-be-copied")

	environment := benchmarkEnvironment("benchmark-root", true)
	joined := strings.Join(environment, "\n")
	for _, required := range []string{
		"MCP_ROOT=benchmark-root",
		"MCP_READ_ONLY=true",
		"MCP_ALLOW_CWD_ROOT=false",
		"MCP_DEBUG=false",
	} {
		if !strings.Contains(joined, required) {
			t.Fatalf("benchmark environment omitted %q: %q", required, environment)
		}
	}
	for _, forbidden := range []string{"MCP_MAX_LIST_ENTRIES=", "FLASHGATE_PRIVATE_TEST_VALUE="} {
		if strings.Contains(joined, forbidden) {
			t.Fatalf("benchmark environment copied %q", forbidden)
		}
	}
}

func TestApplyBudgetEvaluationKeepsWarningChannelsSeparate(t *testing.T) {
	result := completeBudgetTestResult()
	result.Warnings = collectWarnings([]processSample{{warnings: []string{"runtime warning"}}})
	result.Workflows[0].DurationNS = MetricSummary{
		Samples: 1,
		Min:     100_000_001,
		P50:     100_000_001,
		P95:     100_000_001,
		Max:     100_000_001,
	}

	if err := applyBudgetEvaluation(filepath.Join("..", "..", "benchmarks", "budgets.json"), &result); err != nil {
		t.Fatal(err)
	}
	if result.BudgetEvaluation.HardFailures != 0 || result.BudgetEvaluation.SoftWarnings != 1 {
		t.Fatalf("budget evaluation=%+v, want zero hard and one soft warning", result.BudgetEvaluation)
	}
	if len(result.BudgetEvaluation.Messages) != 1 || !strings.HasPrefix(result.BudgetEvaluation.Messages[0], "soft: ") {
		t.Fatalf("budget messages=%q, want one soft message", result.BudgetEvaluation.Messages)
	}
	if len(result.Warnings) != 1 || result.Warnings[0] != "runtime warning" {
		t.Fatalf("general warnings=%q, want only the runtime warning", result.Warnings)
	}
}

func TestInitializedNotificationAndInitializeValidation(t *testing.T) {
	var output bytes.Buffer
	written, err := writeNotification(&output, initializedNotificationSpec())
	if err != nil {
		t.Fatal(err)
	}
	notification := output.Bytes()
	want := "{\"jsonrpc\":\"2.0\",\"method\":\"notifications/initialized\"}\n"
	if string(notification) != want {
		t.Fatalf("initialized notification=%q, want %q", notification, want)
	}
	if written != uint64(len(notification)) || written != 55 {
		t.Fatalf("initialized notification bytes=%d, want 55", written)
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(notification, &envelope); err != nil {
		t.Fatal(err)
	}
	if _, ok := envelope["id"]; ok {
		t.Fatal("initialized notification must not contain id")
	}

	valid := json.RawMessage(`{"protocolVersion":"2025-11-25","capabilities":{"tools":{}},"serverInfo":{"name":"flashgate-mcp","version":"dev"}}`)
	if err := validateInitializeResult(valid); err != nil {
		t.Fatalf("valid initialize result rejected: %v", err)
	}
	invalid := map[string]json.RawMessage{
		"wrong protocol":   json.RawMessage(`{"protocolVersion":"2024-01-01","capabilities":{},"serverInfo":{"name":"server","version":"v1"}}`),
		"missing protocol": json.RawMessage(`{"capabilities":{},"serverInfo":{"name":"server","version":"v1"}}`),
		"missing info":     json.RawMessage(`{"protocolVersion":"2025-11-25","capabilities":{}}`),
		"missing caps":     json.RawMessage(`{"protocolVersion":"2025-11-25","serverInfo":{"name":"server","version":"v1"}}`),
		"null caps":        json.RawMessage(`{"protocolVersion":"2025-11-25","capabilities":null,"serverInfo":{"name":"server","version":"v1"}}`),
	}
	for name, result := range invalid {
		t.Run(name, func(t *testing.T) {
			if err := validateInitializeResult(result); err == nil {
				t.Fatal("invalid initialize result accepted")
			}
		})
	}
}

func TestReadFileCountersDoNotCountScannedBytes(t *testing.T) {
	result := json.RawMessage(`{"structuredContent":{"content":"FlashGate benchmark text.\n","size":26}}`)
	var counters Counters
	if err := addResultCounters("tools/call", result, &counters); err != nil {
		t.Fatal(err)
	}
	if counters.ReadBytes != 26 || counters.ScannedBytes != 0 {
		t.Fatalf("read=%d scanned=%d, want 26/0", counters.ReadBytes, counters.ScannedBytes)
	}
}

func TestSanitizeStderrRedactsOnlyNontrivialAbsolutePaths(t *testing.T) {
	tests := []struct {
		name       string
		stderr     string
		root       string
		binaryPath string
		want       string
	}{
		{"relative binary", "file.txt. remains intact.", ".", "flashgate-mcp.exe", "file.txt. remains intact."},
		{"dot binary", "extension .txt and sentence.", ".", ".", "extension .txt and sentence."},
		{"windows case", `failed at c:\users\alice\corpus\file.txt`, `C:\Users\Alice\Corpus`, `C:\Build\flashgate.exe`, `failed at [REDACTED_PATH]\file.txt`},
		{"windows separators", `failed at c:/users/alice/corpus/file.txt`, `C:\Users\Alice\Corpus`, `C:\Build\flashgate.exe`, `failed at [REDACTED_PATH]/file.txt`},
		{"linux", "failed at /home/alice/corpus/file.txt", "/home/alice/corpus", "/opt/flashgate/server", "failed at [REDACTED_PATH]/file.txt"},
		{"macOS", "failed at /Users/alice/corpus/file.txt", "/Users/alice/corpus", "/Applications/FlashGate/server", "failed at [REDACTED_PATH]/file.txt"},
		{"normal filename", "cannot open file.txt.", ".", "server", "cannot open file.txt."},
		{"multiple paths", `C:\Corpus\a.txt failed via C:/Build/server.exe`, `C:\Corpus`, `C:\Build\server.exe`, `[REDACTED_PATH]\a.txt failed via [REDACTED_PATH]/server.exe`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sanitizeStderr(tc.stderr, tc.root, tc.binaryPath)
			if len(got) != 1 || got[0] != tc.want {
				t.Fatalf("sanitize=%q, want %q", got, tc.want)
			}
		})
	}
	if paths := validRedactionPaths("", ".", "/", `C:\`); len(paths) != 0 {
		t.Fatalf("trivial redaction paths accepted: %d", len(paths))
	}
}

func TestUnsupportedSnapshotNamesEveryOmittedMetric(t *testing.T) {
	snapshot := unsupportedSnapshot()
	if snapshot.status != "not_supported" {
		t.Fatalf("status=%q, want not_supported", snapshot.status)
	}
	if snapshot.workingSetBytes != nil || snapshot.peakWorkingSetBytes != nil || snapshot.userCPUNS != nil || snapshot.systemCPUNS != nil {
		t.Fatal("unsupported snapshot must omit numeric metrics")
	}
	want := []string{
		"idle_working_set_bytes",
		"peak_working_set_bytes",
		"user_cpu_ns",
		"system_cpu_ns",
	}
	if strings.Join(snapshot.unsupported, "\n") != strings.Join(want, "\n") {
		t.Fatalf("unsupported=%v, want %v", snapshot.unsupported, want)
	}
}
