package net


import (
	"bufio"
	"io"
	"net"

	"github.com/v0id-user/vkv/engine"
	"github.com/v0id-user/vkv/protocol"
)

// Handler is the function signature used by Reactor/WorkerPool.
func HandleConnection(conn net.Conn, eng *engine.Engine) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	dec := protocol.NewDecoder(reader)
	enc := protocol.NewEncoder(writer)

	for {
		// Decode next command from the client
		cmd, err := dec.Decode()
		if err != nil {
			// client closed connection
			if err == io.EOF {
				return
			}
			// protocol error or unreadable input = close
			return
		}

		// Execute command against the engine
		resp := eng.Execute(cmd)

		// Encode response back to client
		if err := enc.Encode(resp); err != nil {
			return
		}

		// Make sure the response is fully sent
		if err := writer.Flush(); err != nil {
			return
		}
	}
}
