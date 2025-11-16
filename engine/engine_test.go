package engine

import (
	"os"
	"testing"

	"github.com/v0id-user/vkv/protocol"
	"github.com/v0id-user/vkv/storage"
)

func TestEngineBasic(t *testing.T) {
	// TEMP WAL
	walPath := "test-engine.wal"
	defer os.Remove(walPath)

	wal, err := storage.OpenWAL(walPath)
	if err != nil {
		t.Fatalf("failed to open WAL: %v", err)
	}
	defer wal.Close()

	// Empty memtable
	mem := storage.NewMemtable()

	// Engine
	e := New(mem, wal)

	// --- SET ---
	setCmd := protocol.Set{Key: "foo", Value: "bar"}
	resp := e.Execute(setCmd)

	if resp.Kind() != "OK" {
		t.Fatalf("SET: expected OK, got: %s", resp.Kind())
	}

	// --- GET (in memtable) ---
	getCmd := protocol.Get{Key: "foo"}
	resp = e.Execute(getCmd)

	if resp.Kind() != "VALUE" {
		t.Fatalf("GET: expected VALUE, got: %s", resp.Kind())
	}

	val := resp.(protocol.RespValue).Value
	if val != "bar" {
		t.Fatalf("GET: expected bar, got: %s", val)
	}

	// --- DEL ---
	delCmd := protocol.Del{Key: "foo"}
	resp = e.Execute(delCmd)

	if resp.Kind() != "OK" {
		t.Fatalf("DEL: expected OK, got: %s", resp.Kind())
	}

	// --- GET after DEL (NIL) ---
	resp = e.Execute(getCmd)

	if resp.Kind() != "NIL" {
		t.Fatalf("GET after DEL: expected NIL, got: %s", resp.Kind())
	}
}

func TestEngineWithSSTable(t *testing.T) {
	// build SSTable
	path := "test-sst.sst"
	defer os.Remove(path)

	entries := map[string]string{
		"alpha": "first",
	}

	if err := storage.BuildSSTable(path, entries); err != nil {
		t.Fatalf("failed to build sstable: %v", err)
	}

	sst, err := storage.OpenSSTable(path)
	if err != nil {
		t.Fatalf("failed to open sstable: %v", err)
	}

	// engine with empty memtable
	wal, _ := storage.OpenWAL("test-sst.wal")
	defer os.Remove("test-sst.wal")
	defer wal.Close()

	mem := storage.NewMemtable()
	e := New(mem, wal)

	// register SSTable
	e.AddSSTable(sst)

	// GET key from sstable
	cmd := protocol.Get{Key: "alpha"}
	resp := e.Execute(cmd)

	if resp.Kind() != "VALUE" {
		t.Fatalf("expected VALUE, got %s", resp.Kind())
	}

	v := resp.(protocol.RespValue).Value
	if v != "first" {
		t.Fatalf("expected value 'first', got %q", v)
	}
}
