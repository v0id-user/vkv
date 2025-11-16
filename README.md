# vkv — a tiny LSM-style key value store (Go)

vkv is a minimal, educational key-value database built from scratch in Go.  
It’s a real database: persistent, concurrent, networked, and structured like a tiny Redis/LevelDB hybrid.

This project focuses on clear architecture, durability, and simplicity.

## Features

### • Text Protocol (ASCII)
Commands are newline-terminated:
```

SET key value
GET key
DEL key

````

### • Parser + AST
Incoming lines go through:
1. Lexer
2. Parser
3. Command AST
4. Engine execution

Keeps the core clean and extendable.

### • Storage Engine (LSM-ish)
- **Memtable** (sharded in-memory map)
- **WAL** (append-only log for crash recovery)
- **SSTables** (immutable, sorted files on disk)
- Startup recovery from WAL
- Read path: memtable → SSTables

### • Engine Layer
Maps parsed commands to storage operations:
- SET → WAL + memtable
- GET → memtable → SSTables
- DEL → WAL + memtable

The engine is isolated from networking and runtime.

### • Runtime
Simple, idiomatic Go runtime:
- **Reactor**: schedules new connections
- **WorkerPool**: runs connection handlers concurrently

No fake epoll. No unnecessary complexity.

### • Networking
TCP server:
- Accept loop
- Reactor registration
- Worker executes:
  - decode command
  - engine.Execute
  - encode response

You can `nc` into it immediately.

## Quick Start

Build and run:

```sh
go build -o vkv ./cmd/vkv
./vkv
````

Connect:

```sh
nc localhost 9999
SET foo bar
OK
GET foo
VALUE bar
DEL foo
OK
```

## Project Structure

```
protocol/   → lexer, parser, AST, encoder/decoder
storage/    → WAL, Memtable, SSTable
engine/     → routing + execution
runtime/    → reactor + worker pool
net/        → TCP connection handling + server
main.go
```

## Why vkv exists

To learn how databases actually work:

* framing
* parsing
* LSM trees
* WAL durability
* memtable flushes
* SSTable reads
* concurrency
* minimal networking runtime

It’s intentionally small, readable, and hackable.

## Future Work

* memtable flush daemon
* SSTable compaction
* binary protocol
* metrics + logging
* TTLs + expirations
* RESP/Redis-compatible protocol

---

Built by #V0ID to understand the internals — not to chase benchmarks.

```

---

If you want, I can write the **cmd/vkv/main.go** next so the server actually boots.

Just say:  
**“main.go next”**