package protocol

import (
	"bufio"
	"bytes"
	"testing"
)

func TestEncoderEncode(t *testing.T) {
	tests := []struct {
		name     string
		resp     Response
		expected string
	}{
		{
			name:     "OK",
			resp:     ResponseOK(),
			expected: "OK\n",
		},
		{
			name:     "VALUE",
			resp:     ResponseValue("bar"),
			expected: "VALUE bar\n",
		},
		{
			name:     "NIL",
			resp:     ResponseNil(),
			expected: "NIL\n",
		},
		{
			name:     "ERR",
			resp:     ResponseErr("bad"),
			expected: "ERR bad\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)
			enc := NewEncoder(writer)

			if err := enc.Encode(tc.resp); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			writer.Flush()

			got := buf.String()
			if got != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}
