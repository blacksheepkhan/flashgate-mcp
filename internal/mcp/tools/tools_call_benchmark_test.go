package tools

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/blacksheepkhan/flashgate-mcp/internal/fs"
	"github.com/blacksheepkhan/flashgate-mcp/internal/protocol"
)

var benchmarkResultSink []byte

type textOnlyCallToolResult struct {
	Content []protocol.TextContent `json:"content"`
}

func BenchmarkCallToolResultSerialization(b *testing.B) {
	fixtures := callToolResultBenchmarkFixtures()
	variants := []struct {
		name      string
		serialize func(any) ([]byte, error)
	}{
		{"historical_direct", json.Marshal},
		{"text_only", serializeTextOnlyCallToolResult},
		{"text_plus_structured", serializeStructuredCallToolResult},
	}

	for _, fixture := range fixtures {
		fixture := fixture
		b.Run(fixture.name, func(b *testing.B) {
			for _, variant := range variants {
				variant := variant
				b.Run(variant.name, func(b *testing.B) {
					sample, err := variant.serialize(fixture.value)
					if err != nil {
						b.Fatal(err)
					}
					b.ReportAllocs()
					b.SetBytes(int64(len(sample)))
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						benchmarkResultSink, err = variant.serialize(fixture.value)
						if err != nil {
							b.Fatal(err)
						}
					}
					b.StopTimer()
					b.ReportMetric(float64(len(sample)), "payload-bytes")
				})
			}
		})
	}
}

func serializeTextOnlyCallToolResult(value any) ([]byte, error) {
	structured, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return json.Marshal(textOnlyCallToolResult{
		Content: []protocol.TextContent{protocol.NewTextContent(string(structured))},
	})
}

func serializeStructuredCallToolResult(value any) ([]byte, error) {
	wrapped, rpcErr := wrapSuccessfulToolResult(value)
	if rpcErr != nil {
		return nil, fmt.Errorf("wrap result: %s", rpcErr.Message)
	}
	return json.Marshal(wrapped)
}

type callToolResultBenchmarkFixture struct {
	name  string
	value any
}

func callToolResultBenchmarkFixtures() []callToolResultBenchmarkFixture {
	largeEntries := make([]fs.Entry, 500)
	for index := range largeEntries {
		largeEntries[index] = fs.Entry{
			Name:  fmt.Sprintf("entry-%04d.txt", index),
			Size:  int64(index * 17),
			IsDir: index%10 == 0,
		}
	}

	return []callToolResultBenchmarkFixture{
		{"path_info_missing", getPathInfoMissingResult{Path: "does-not-exist.txt", Exists: false}},
		{"path_info_existing", getPathInfoExistingResult{Path: "docs/readme.txt", Exists: true, Name: "readme.txt", Size: 128}},
		{"directory_small", listDirectoryResult{Entries: []fs.Entry{{Name: "a.txt", Size: 12}, {Name: "sub", IsDir: true}}}},
		{"directory_500_entries", listDirectoryResult{Entries: largeEntries}},
		{"text_file_small", readFileResult{Content: "FlashGate benchmark text.\n", Size: 26}},
		{"text_file_64kib", readFileResult{Content: strings.Repeat("x", 64*1024), Size: 64 * 1024}},
	}
}
