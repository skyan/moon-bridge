package server

import (
	"strings"
	"testing"
)

func TestOpenAIStreamMonitorTruncatesParsingWithoutBlockingWrites(t *testing.T) {
	monitor := &openAIStreamMonitor{}
	payload := strings.Repeat("x", openAIStreamMonitorMaxParseBytes+1024)

	n, err := monitor.Write([]byte(payload))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != len(payload) {
		t.Fatalf("Write() bytes = %d, want %d", n, len(payload))
	}
	if monitor.Bytes != int64(len(payload)) {
		t.Fatalf("Bytes = %d, want %d", monitor.Bytes, len(payload))
	}
	if !monitor.Truncated {
		t.Fatal("Truncated = false, want true")
	}
	if len(monitor.lineBuf) != 0 || len(monitor.dataLines) != 0 || monitor.eventName != "" {
		t.Fatalf("monitor buffers were not cleared: line=%d data=%d event=%q", len(monitor.lineBuf), len(monitor.dataLines), monitor.eventName)
	}

	monitor.Finish()
	if monitor.Terminal != "" || monitor.ErrorMessage != "" {
		t.Fatalf("Finish() parsed truncated data: terminal=%q error=%q", monitor.Terminal, monitor.ErrorMessage)
	}
}
