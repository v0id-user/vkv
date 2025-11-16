package storage

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

// SSTable represents a single immutable sorted table on disk.
// For v1, we build a full in-memory index (map) when opening.
type SSTable struct {
	path string
	data map[string]string
}

// BuildSSTable writes the given entries to disk as a sorted, immutable table.
//
// The file format is one record per line:
//
//	<key>\t<value>\n
//
// Keys are written in lexicographical order.
func BuildSSTable(path string, entries map[string]string) error {
	// Create/truncate file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Sort keys
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := bufio.NewWriter(f)

	for _, k := range keys {
		v := entries[k]
		// For v1 we assume no tabs in key/value.
		if _, err := w.WriteString(k + "\t" + v + "\n"); err != nil {
			return err
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

// OpenSSTable loads an existing SSTable file and builds an in-memory index.
//
// For v1 this reads the whole file into a map. Later we can optimize this to
// use offsets and on-demand reads instead of storing everything in RAM.
func OpenSSTable(path string) (*SSTable, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := make(map[string]string)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			// malformed line, skip for now
			continue
		}

		key := parts[0]
		value := parts[1]
		data[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &SSTable{
		path: path,
		data: data,
	}, nil
}

// Get returns the value for a key if present in this SSTable.
func (s *SSTable) Get(key string) (string, bool) {
	v, ok := s.data[key]
	return v, ok
}
