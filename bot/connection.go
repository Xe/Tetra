package cod

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
	log.Printf(">>> "+line, stuff...)
	fmt.Fprintf(c.Conn, line+"\r\n", stuff...)
}

func (c *Connection) GetLine() (line string, err error) {
	line, err = c.Tp.ReadLine()

	return
}
