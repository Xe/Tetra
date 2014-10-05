package tetra

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"net/textproto"
)

// Struct Connection contains everything needed for the socket connection Tetra
// uses.
type Connection struct {
	Conn   net.Conn
	Log    *log.Logger
	Reader *bufio.Reader
	Tp     *textproto.Reader
	Buffer chan string
	open   bool
	Debug  bool
}

// SendLine buffers a line to be sent to the server.
func (c *Connection) SendLine(line string, stuff ...interface{}) {
	str := fmt.Sprintf(line, stuff...)
	c.Buffer <- str
}

func (c *Connection) sendLinesWait() {
	for {
		str := <-c.Buffer

		debugf(">>> " + str)

		fmt.Fprintf(c.Conn, "%s\r\n", str)
	}
}

// Close kills the connection
func (c *Connection) Close() {
	c.open = false
}

// GetLine returns a new line from the server.
func (c *Connection) GetLine() (line string, err error) {
	if c.open {
		line, err = c.Tp.ReadLine()
	} else {
		return "", errors.New("Conection is closed")
	}

	return
}
