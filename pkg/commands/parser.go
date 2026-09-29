package commands

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// ReadCommand reads one full command from the client.
// It accepts RESP arrays (*1\r\n$4\r\nPING\r\n), which is what redis-cli sends,
// and inline commands (PING\r\n), which is what you type into nc.
func ReadCommand(reader *bufio.Reader) ([]string, error) {
	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(line, "*") {
		return strings.Fields(line), nil
	}

	count, err := strconv.Atoi(line[1:])
	if err != nil || count < 0 {
		return nil, fmt.Errorf("Protocol error: invalid multibulk length")
	}
	args := make([]string, 0, count)
	for i := 0; i < count; i++ {
		header, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(header, "$") {
			return nil, fmt.Errorf("Protocol error: expected '$', got '%s'", header)
		}
		size, err := strconv.Atoi(header[1:])
		if err != nil || size < 0 {
			return nil, fmt.Errorf("Protocol error: invalid bulk length")
		}
		// The length prefix is the truth: read exactly size bytes, then the trailing \r\n
		buf := make([]byte, size+2)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return nil, err
		}
		args = append(args, string(buf[:size]))
	}
	return args, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func ParseCommands(args []string, con net.Conn) {
	if len(args) == 0 {
		return
	}
	c, ok := CommandMap[strings.ToUpper(args[0])]
	if !ok {
		con.Write([]byte("-ERR unknown command '" + args[0] + "'\r\n"))
		return
	}

	c(con)
}
