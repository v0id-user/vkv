package net


import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/v0id-user/vkv/engine"
	runtime2 "github.com/v0id-user/vkv/runtime"
)

// StartServer starts the TCP server with a reactor + worker pool.
func StartServer(addr string, eng *engine.Engine, workers int) error {
	// Worker pool
	pool := runtime2.NewWorkerPool(workers)
	pool.Start()

	// Reactor
	reactor := runtime2.NewReactor(pool)

	// TCP listen
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer ln.Close()

	fmt.Printf("vkv server listening on %s\n", addr)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		fmt.Println("shutting down server...")
		reactor.Stop()
		pool.Stop()
		ln.Close()
		os.Exit(0)
	}()

	// Accept loop
	for {
		conn, err := ln.Accept()
		if err != nil {
			// listener closed or fatal error
			return err
		}

		// Register connection with reactor
		reactor.Register(conn, func(c net.Conn) {
			HandleConnection(c, eng)
		})
	}
}