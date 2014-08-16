package tetra

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Xe/Tetra/bot/modes"
	"github.com/Xe/Tetra/bot/web"
	"github.com/codegangsta/negroni"
	"github.com/drone/routes"
	"github.com/phyber/negroni-gzip/gzip"
	"gopkg.in/unrolled/render.v1"
)

type client struct {
	Nick     string            `json:"nick"`
	User     string            `json:"user"`
	Host     string            `json:"host"`
	Account  string            `json:"account"`
	Server   string            `json:"server"`
	Joins    []chanuser        `json:"joins"`
	Metadata map[string]string `json:"metadata"`
}

type chanuser struct {
	Client  string `json:"client"`
	Channel string `json:"channel"`
	Prefix  int    `json:"prefix"`
}

type channel struct {
	Name     string            `json:"name"`
	Ts       int64             `json:"ts"`
	Modes    int               `json:"modes"`
	Clients  []chanuser        `json:"clients"`
	Metadata map[string]string `json:"metadata"`
}

func convertChannel(in *Channel) (out channel) {
	out = channel{
		Name:     in.Name,
		Ts:       in.Ts,
		Modes:    in.Modes,
		Metadata: in.Metadata,
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
		Nick:     in.Nick,
		User:     in.User,
		Host:     in.VHost,
		Account:  in.Account,
		Server:   in.Server.Name,
		Metadata: in.Metadata,
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

// WebApp creates the web application and YAML api for Tetra.
func (t *Tetra) WebApp() {
	mux := routes.New()
	r := render.New(render.Options{})

	mux.Get("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "error: No method chosen.")
	})

	mux.Get("/channels.json", func(w http.ResponseWriter, req *http.Request) {
		var channels []channel

		for _, in := range t.Channels {
			if in.Modes&modes.PROP_SECRET != modes.PROP_SECRET {
				channels = append(channels, convertChannel(in))
			}
		}

		r.JSON(w, http.StatusOK, channels)
	})

	mux.Get("/clients.json", func(res http.ResponseWriter, req *http.Request) {
		var clients []client

		for _, in := range t.Clients.ByUID {
			myclient := convertClient(in)

			if len(myclient.Joins) == 0 {
				continue
			}

			clients = append(clients, myclient)
		}

		r.JSON(res, http.StatusOK, clients)
	})

	mux.Get("/yo/:id", func(res http.ResponseWriter, req *http.Request) {
		params := req.URL.Query()
		id := params.Get(":id")
		username := req.URL.Query()["username"][0]

		t.RunHook("YO", username, id)
		t.RunHook("YO_"+id, username)

		fmt.Fprintf(res, "OK")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	t.Log.Printf("listening on %v...\n", port)

	go func() {
		n := negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("public")), web.NewLogger())

		n.Use(gzip.Gzip(gzip.DefaultCompression))
		n.UseHandler(mux)

		err := http.ListenAndServe(":"+port, n)

		if err != nil {
			t.Services["tetra"].ServicesLog("Web app died")
			t.Services["tetra"].ServicesLog(err.Error())
			t.Log.Fatal(err)
		}
	}()
}
