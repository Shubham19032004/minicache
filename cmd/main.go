package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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
			fmt.Println("usage: ./main --port=8000")
			os.Exit(1)
		}
		port, err := strconv.Atoi(val)
		Port = port
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
func getPublicIP() (string, error) {
	services := []string{
		"https://api.ipify.org",
		"https://ifconfig.me/ip",
		"https://icanhazip.com",
	}

	for _, service := range services {
		resp, err := http.Get(service)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err == nil {
			return strings.TrimSpace(string(body)), nil
		}
	}

	return "", fmt.Errorf("unable to determine public IP")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// One reader per connection, so bytes it has buffered aren't lost between reads
	reader := bufio.NewReader(conn)
	address, err := getPublicIP()
	if err != nil {
	}
	for {
		conn.Write([]byte(address + ":" + strconv.Itoa(Port) + ">>"))

		args, err := commands.ReadCommand(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %v", err)
				conn.Write([]byte("-ERR " + err.Error() + "\r\n"))
			}
			return
		}
		fmt.Println("received:", args)
		conn.Write([]byte(address + ":" + strconv.Itoa(Port) + ">>"))
		commands.ParseCommands(args, conn)
	}

}
