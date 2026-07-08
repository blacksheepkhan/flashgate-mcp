package diagnostics

import (
	"bytes"
	"strings"
	"testing"
)

func TestRedactRemovesSecretsAndHostPaths(t *testing.T) {
	t.Parallel()

	input := strings.Join([]string{
		`Authorization: Bearer abc123`,
		`Authorization: Basic abc123`,
		`password=secret`,
		`"api_key":"secret"`,
		`-----BEGIN RSA PRIVATE KEY-----`,
		`postgres://user:pass@example.test/db`,
		`C:\Users\Example\secret.txt`,
		`\\server\share\secret.txt`,
		`/home/example/secret.txt`,
	}, "\n")

	redacted := Redact(input)

	for _, forbidden := range []string{
		"abc123",
		"password=secret",
		`"api_key":"secret"`,
		"BEGIN RSA PRIVATE KEY",
		"user:pass@",
		`C:\Users\Example`,
		`\\server\share`,
		"/home/example",
	} {
		if strings.Contains(redacted, forbidden) {
			t.Fatalf("expected %q to be redacted from %q", forbidden, redacted)
		}
	}

	if !strings.Contains(redacted, "[REDACTED]") && !strings.Contains(redacted, "[REDACTED_PATH]") {
		t.Fatalf("expected redaction markers, got %q", redacted)
	}
}

func TestLoggerWritesOnlyWhenDebugEnabled(t *testing.T) {
	t.Parallel()

	disabledOutput := &bytes.Buffer{}
	NewLogger(false, disabledOutput).Debugf("password=secret")
	if disabledOutput.Len() != 0 {
		t.Fatalf("expected no output when debug is disabled, got %q", disabledOutput.String())
	}

	enabledOutput := &bytes.Buffer{}
	NewLogger(true, enabledOutput).Debugf("password=secret")
	output := enabledOutput.String()
	if output == "" {
		t.Fatal("expected debug output")
	}

	if strings.Contains(output, "secret") {
		t.Fatalf("expected secret to be redacted, got %q", output)
	}
}
