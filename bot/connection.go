package tetra

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
)

type Connection struct {
	Conn   net.Conn
	Log    *log.Logger
	Reader *bufio.Reader
	Tp     *textproto.Reader
}

func (c *Connection) SendLine(line string, stuff ...interface{}) {
	str := fmt.Sprintf(line, stuff...)
	c.Log.Printf(">>> " + str)
	fmt.Fprintf(c.Conn, "%s\r\n", line)
}

func (c *Connection) GetLine() (line string, err error) {
	line, err = c.Tp.ReadLine()

	return
}
