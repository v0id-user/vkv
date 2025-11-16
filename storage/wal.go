package storage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)


type WAL struct {
	mu   sync.Mutex
	file *os.File
	path string
}

// OpenWAL opens (or creates) the WAL file.
func OpenWAL(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{
		file: f,
		path: path,
	}, nil
}


// AppendSet logs a SET operation to the WAL.
func (w *WAL) AppendSet(key, value string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := fmt.Fprintf(w.file, "SET %s %s\n", key, value)
	return err
}

// AppendDel logs a DEL operation.
func (w *WAL) AppendDel(key string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := fmt.Fprintf(w.file, "DEL %s\n", key)
	return err
}

// Replay reads the WAL and replays all operations into memtable.
func (w *WAL) Replay(mt *Memtable) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Seek to start
	if _, err := w.file.Seek(0, 0); err != nil {
		return err
	}

	scanner := bufio.NewScanner(w.file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")

		switch parts[0] {
		case "SET":
			if len(parts) != 3 {
				continue
			}
			mt.Set(parts[1], parts[2])

		case "DEL":
			if len(parts) != 2 {
				continue
			}
			mt.Del(parts[1])
		}
	}

	return scanner.Err()
}

// Reset clears the WAL file after memtable is flushed.
func (w *WAL) Reset() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.file.Close()

	// Truncate the file
	if err := os.Remove(w.path); err != nil {
		return err
	}

	f, err := os.OpenFile(w.path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w.file = f
	return nil
}

func (w *WAL) Close() error {
	return w.file.Close()
}