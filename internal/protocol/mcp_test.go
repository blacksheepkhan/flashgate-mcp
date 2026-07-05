package protocol

import "testing"

func TestProtocolVersionIsSet(t *testing.T) {
	t.Parallel()

	if ProtocolVersion == "" {
		t.Fatal("expected protocol version to be set")
	}
}
