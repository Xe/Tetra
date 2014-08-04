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
	Buffer chan string
	open   bool
}

func (c *Connection) SendLine(line string, stuff ...interface{}) {
	str := fmt.Sprintf(line, stuff...)
	c.Buffer <- str
}

func (c *Connection) sendLinesWait() {
	for {
		str := <-c.Buffer
		c.Log.Printf(">>> " + str)
		fmt.Fprintf(c.Conn, "%s\r\n", str)
	}
}

func (c *Connection) Close() {
	c.open = false
}

func (c *Connection) GetLine() (line string, err error) {
	if c.open {
		line, err = c.Tp.ReadLine()
	} else {
		return "", nil
	}

	return
}
