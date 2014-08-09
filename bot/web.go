package tetra

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"gopkg.in/yaml.v1"
)

type client struct {
	Nick    string
	User    string
	Host    string
	Account string
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
			channels = append(channels, convertChannel(in))
		}

		out, err := yaml.Marshal(channels)
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
