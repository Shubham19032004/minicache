package commands

import (
	"net"

	"github.com/Shubham19032004/minicache/pkg/handlers"
)

var CommandMap = map[string]func(net.Conn){
	PING: handlers.Ping,
}
