package main

import (
	_ "fmt"
	"github.com/Xe/Tetra/1459"
	"github.com/Xe/Tetra/bot"
)

func main() {
	tetra := tetra.NewTetra("etc/config.json")

	tetra.Connect("127.0.0.1", "6667")
	defer tetra.Conn.Conn.Close()

	tetra.Auth()

	for _, script := range tetra.Config.Autoload {
		tetra.LoadScript(script)
	}

	for _, sclient := range tetra.Config.Services {
		tetra.AddService(sclient.Name, sclient.Nick, sclient.User, sclient.Host, sclient.Gecos)
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
				if svc, ok := tetra.Services["tetra"]; !ok {
					panic("No service tetra!")
				} else {
					for _, client := range tetra.Services {
						tetra.Conn.SendLine(client.Euid())
					}
					svc.Join("#services")
				}
			}
			tetra.Conn.SendLine("PONG :" + rawline.Args[0])
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
