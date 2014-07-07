package main

import (
	_ "fmt"
	"github.com/Xe/Tetra/bot"
	"github.com/Xe/Tetra/1459"
)

func main() {
	cod := cod.NewCod()

	cod.Connect("127.0.0.1", "6667")
	defer cod.Conn.Conn.Close()

	cod.Conn.SendLine("PASS shameless TS 6 :420")
	cod.Conn.SendLine("CAPAB :QS EX IE KLN UNKLN ENCAP SERVICES EUID EOPMO")
	cod.Conn.SendLine("SERVER cod.int 1 :Cod in Go!")

	for _, client := range cod.Clients.ByUID {
		cod.Conn.SendLine(client.Euid())
	}

	for {
		line, err := cod.Conn.GetLine()
		if err != nil {
			panic(err)
		}

		rawline := r1459.NewRawLine(line)

		cod.Conn.Log.Printf("<<< %s", line)

		if rawline.Verb == "PING" {
			if !cod.Bursted {
				cod.Bursted = true
			}
			cod.Conn.SendLine("PONG :%s", rawline.Args[0])
		}

		if _, present := cod.Handlers[rawline.Verb]; present {
			for _, handler := range cod.Handlers[rawline.Verb] {
				if cod.Bursted {
					go handler.Impl(rawline)
				} else {
					handler.Impl(rawline)
				}
			}
		}
	}
}
