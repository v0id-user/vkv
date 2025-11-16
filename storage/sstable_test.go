package storage

import (
	"os"
	"testing"
)

func TestSSTableBuildAndOpen(t *testing.T) {
	// Temporary test file
	path := "testdata.sst"

	// Cleanup after test
	defer os.Remove(path)

	// Entries in random order to test sorting
	entries := map[string]string{
		"zeta":   "last",
		"alpha":  "first",
		"middle": "center",
	}

	// Build SSTable
	if err := BuildSSTable(path, entries); err != nil {
		t.Fatalf("failed to build sstable: %v", err)
	}

	// Open SSTable
	sst, err := OpenSSTable(path)
	if err != nil {
		t.Fatalf("failed to open sstable: %v", err)
	}

	// Validate Get (order doesn't matter for lookup)
	tests := []struct {
		key   string
		value string
		found bool
	}{
		{"alpha", "first", true},
		{"middle", "center", true},
		{"zeta", "last", true},
		{"unknown", "", false},
	}

	for _, tc := range tests {
		got, ok := sst.Get(tc.key)

		if ok != tc.found {
			t.Fatalf("key %q: expected found=%v, got %v", tc.key, tc.found, ok)
		}

		if ok && got != tc.value {
			t.Fatalf("key %q: expected value=%q, got %q", tc.key, tc.value, got)
		}
	}
}
