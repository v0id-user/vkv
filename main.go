package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/v0id-user/vkv/engine"
	"github.com/v0id-user/vkv/net"
	"github.com/v0id-user/vkv/storage"
)

func main() {
	addr := "0.0.0.0:9999"
	workers := 8
	walPath := "vkv.wal"

	// WAL
	wal, err := storage.OpenWAL(walPath)
	if err != nil {
		fmt.Println("Failed to open WAL:", err)
		os.Exit(1)
	}

	// Memtable
	mem := storage.NewMemtable()

	// Recovery: replay WAL → memtable
	if err := wal.Replay(mem); err != nil {
		fmt.Println("Failed to replay WAL:", err)
		os.Exit(1)
	}

	// Engine
	eng := engine.New(mem, wal)

	// 5) Load existing SSTables
	sstPaths, err := filepath.Glob("data/*.sst")
	if err != nil {
		fmt.Println("failed to scan SSTable directory:", err)
		os.Exit(1)
	}

	// Load oldest → newest
	sort.Strings(sstPaths)

	for _, p := range sstPaths {
		sst, err := storage.OpenSSTable(p)
		if err != nil {
			fmt.Printf("failed to open sstable %s: %v\n", p, err)
			continue
		}
		eng.AddSSTable(sst)
		fmt.Println("loaded sstable:", p)
	}

	// Start server
	if err := net.StartServer(addr, eng, workers); err != nil {
		fmt.Println("server error:", err)
		os.Exit(1)
	}
}
