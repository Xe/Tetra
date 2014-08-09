package tetra

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Xe/Tetra/bot/modes"
	"github.com/codegangsta/negroni"
	"gopkg.in/yaml.v1"
)

type client struct {
	Nick    string
	User    string
	Host    string
	Account string
	Server  string
	Joins   []chanuser
}

type chanuser struct {
	Client  string
	Channel string
	Prefix  int
}

type channel struct {
	Name    string
	Ts      int64
	Modes   int
	Clients []chanuser
}

func convertChannel(in *Channel) (out channel) {
	out = channel{
		Name:  in.Name,
		Ts:    in.Ts,
		Modes: in.Modes,
	}

	for _, chanclient := range in.Clients {
		out.Clients = append(out.Clients, chanuser{
			Client:  chanclient.Client.Nick,
			Channel: chanclient.Channel.Name,
			Prefix:  chanclient.Prefix,
		})
	}

	return
}

func convertClient(in *Client) (out client) {
	out = client{
		Nick:    in.Nick,
		User:    in.User,
		Host:    in.VHost,
		Account: in.Account,
		Server:  in.Server.Name,
	}

	for _, mychan := range in.Channels {
		if mychan.Modes&modes.PROP_SECRET != modes.PROP_SECRET {
			chanclient := mychan.Clients[in.Uid]
			out.Joins = append(out.Joins, chanuser{
				Channel: mychan.Name,
				Client:  in.Nick,
				Prefix:  chanclient.Prefix,
			})
		}
	}

	return
}

func (t *Tetra) WebApp() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "error: No method chosen.")
	})
	mux.HandleFunc("/config.yaml", func(res http.ResponseWriter, req *http.Request) {
		out, err := yaml.Marshal(t.Config)
		if err != nil {
			res.WriteHeader(500)
			fmt.Fprintf(res, `error: Bad yaml`)
			return
		}

		fmt.Fprintf(res, "%s", out)
	})
	mux.HandleFunc("/channels.yaml", func(res http.ResponseWriter, req *http.Request) {
		var channels []channel

		for _, in := range t.Channels {
			if in.Modes&modes.PROP_SECRET != modes.PROP_SECRET {
				channels = append(channels, convertChannel(in))
			}
		}

		out, err := yaml.Marshal(channels)
		if err != nil {
			res.WriteHeader(500)
			fmt.Fprintf(res, `error: Bad yaml`)
			return
		}

		fmt.Fprintf(res, "%s", out)
	})
	mux.HandleFunc("/clients.yaml", func(res http.ResponseWriter, req *http.Request) {
		var clients []client

		for _, in := range t.Clients.ByUID {
			myclient := convertClient(in)

			if len(myclient.Joins) == 0 {
				continue
			}

			clients = append(clients, myclient)
		}

		out, err := yaml.Marshal(clients)
		if err != nil {
			res.WriteHeader(500)
			fmt.Fprintf(res, `error: Bad yaml`)
			return
		}

		fmt.Fprintf(res, "%s", out)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("listening on %v...\n", port)

	go func() {
		n := negroni.Classic()
		n.UseHandler(mux)

		err := http.ListenAndServe(":"+port, n)

		if err != nil {
			t.Services["tetra"].ServicesLog("Web app died")
			t.Services["tetra"].ServicesLog(err.Error())
			t.Log.Fatal(err)
		}
	}()
}
