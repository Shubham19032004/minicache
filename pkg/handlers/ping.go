package handlers

import (
	"log"
	"net"
)

func Ping(conn net.Conn) {
	if _, err := conn.Write([]byte("PONG")); err != nil {
		log.Printf("Write error: %v", err)
		return
	}
}
