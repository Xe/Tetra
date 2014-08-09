package tetra

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
)

func (t *Tetra) WebApp() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

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

func index(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "{\"error\": \"No method chosen.\"}\n")
}
