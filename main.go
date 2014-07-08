package main

import (
	_ "fmt"
	"github.com/Xe/Tetra/1459"
	"github.com/Xe/Tetra/bot"
)

func main() {
	tetra := tetra.NewTetra()

	tetra.Connect("127.0.0.1", "6667")
	defer tetra.Conn.Conn.Close()

	tetra.Conn.SendLine("PASS shameless TS 6 :420")
	tetra.Conn.SendLine("CAPAB :QS EX IE KLN UNKLN ENCAP SERVICES EUID EOPMO")
	tetra.Conn.SendLine("SERVER tetra.int 1 :Tetra in Go!")

	for _, client := range tetra.Clients.ByUID {
		tetra.Conn.SendLine(client.Euid())
	}

	for {
		line, err := tetra.Conn.GetLine()
		if err != nil {
			panic(err)
		}

		rawline := r1459.NewRawLine(line)

		tetra.Conn.Log.Printf("<<< %s", line)

		if rawline.Verb == "PING" {
			if !tetra.Bursted {
				tetra.Bursted = true
			}
			tetra.Conn.SendLine("PONG :%s", rawline.Args[0])
		}

		if _, present := tetra.Handlers[rawline.Verb]; present {
			for _, handler := range tetra.Handlers[rawline.Verb] {
				if tetra.Bursted {
					go handler.Impl(rawline)
				} else {
					handler.Impl(rawline)
				}
			}
		}
	}
}
