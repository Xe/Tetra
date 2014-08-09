package tetra

import (
	"fmt"
	"net/http"
	"os"
)

func (t *Tetra) WebApp() {
	http.HandleFunc("/", index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("listening on %v...\n", port)

	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			panic(err)
		}
	}()
}

func index(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "{\"error\": \"No method chosen.\"}\n")
}
