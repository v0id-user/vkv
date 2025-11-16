package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("connect error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to vkv on 127.0.0.1:9999")
	fmt.Println("Type commands like: SET foo bar | GET foo | DEL foo")
	fmt.Println("Ctrl+C to exit")

	reader := bufio.NewReader(os.Stdin)
	server := bufio.NewReader(conn)

	for {
		fmt.Print("> ")

		// Read user input
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("input error:", err)
			return
		}

		// Trim spaces
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Send to server
		_, err = conn.Write([]byte(line + "\n"))
		if err != nil {
			fmt.Println("send error:", err)
			return
		}

		// Read server response
		resp, err := server.ReadString('\n')
		if err != nil {
			fmt.Println("server closed")
			return
		}

		fmt.Println(strings.TrimSpace(resp))
	}
}
