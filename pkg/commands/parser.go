package commands

import (
	"net"
	"strings"
)

func ParseCommands(command string, con net.Conn) {
	parts := strings.Split(command, "=")
	if len(parts) == 1 {
		c, ok := CommandMap[strings.ToUpper(parts[0])]
		if !ok {
			con.Write([]byte("ERR unknown command '" + parts[0] + "'\n"))
			return
		}

		c(con)
	}
}
