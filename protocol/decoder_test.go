package protocol

import (
	"bufio"
	"strings"
	"testing"
)

func TestDecoderDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		kind    string
	}{
		{
			name:  "valid SET",
			input: "SET foo bar\n",
			kind:  "SET",
		},
		{
			name:  "valid GET",
			input: "GET x\n",
			kind:  "GET",
		},
		{
			name:  "valid DEL",
			input: "DEL y\n",
			kind:  "DEL",
		},
		{
			name:    "empty line",
			input:   "\n",
			wantErr: true,
		},
		{
			name:    "invalid command",
			input:   "HELLO WORLD\n",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tc.input))
			dec := NewDecoder(reader)

			cmd, err := dec.Decode()

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cmd.Kind() != tc.kind {
				t.Fatalf("expected %s, got %s", tc.kind, cmd.Kind())
			}
		})
	}
}
