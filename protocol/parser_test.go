package protocol

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantKind  string
		wantKey   string
		wantValue string
		wantErr   bool
	}{
		{
			name:      "SET valid",
			input:     "SET foo bar",
			wantKind:  "SET",
			wantKey:   "foo",
			wantValue: "bar",
		},
		{
			name:     "GET valid",
			input:    "GET key123",
			wantKind: "GET",
			wantKey:  "key123",
		},
		{
			name:     "DEL valid",
			input:    "DEL x",
			wantKind: "DEL",
			wantKey:  "x",
		},
		{
			name:    "Unknown command",
			input:   "HELLO world",
			wantErr: true,
		},
		{
			name:    "SET missing value",
			input:   "SET foo",
			wantErr: true,
		},
		{
			name:    "GET missing key",
			input:   "GET",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := ParseLine(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cmd.Kind() != tc.wantKind {
				t.Fatalf("expected kind %q, got %q", tc.wantKind, cmd.Kind())
			}

			// For commands that have a Key field, check it.
			switch c := cmd.(type) {
			case Set:
				if c.Key != tc.wantKey {
					t.Fatalf("expected key %q, got %q", tc.wantKey, c.Key)
				}
				if c.Value != tc.wantValue {
					t.Fatalf("expected value %q, got %q", tc.wantValue, c.Value)
				}
			case Get:
				if c.Key != tc.wantKey {
					t.Fatalf("expected key %q, got %q", tc.wantKey, c.Key)
				}
			case Del:
				if c.Key != tc.wantKey {
					t.Fatalf("expected key %q, got %q", tc.wantKey, c.Key)
				}
			}
		})
	}
}
