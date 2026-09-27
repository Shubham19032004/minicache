package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Shubham19032004/minicache/pkg/commands"
)

var Port = 6379

func main() {
	// Basic port getting
	if len(os.Args) < 2 {
		fmt.Println("Port is node define using default port 6379")
		makeServer(Port)
	} else {
		arg := os.Args[1]
		key, val, ok := strings.Cut(arg, "=")
		if !ok || key != "--port" {
			fmt.Println("usage: ./main port=8000")
			os.Exit(1)
		}
		port, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("port must be a number")
			os.Exit(1)
		}
		makeServer(port)
	}

}
func makeServer(port int) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// One reader per connection, so bytes it has buffered aren't lost between reads
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %v", err)
			}
			return
		}
		command := strings.TrimSpace(message)
		fmt.Println("received:", command)
		commands.ParseCommands(command,conn)
	}

}
