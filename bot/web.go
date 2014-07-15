package tetra

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func (t *Tetra) WebApp() {
	m := martini.Classic()

	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.JSON(501, map[string]interface{} {"error": "No method selected"})
	})

	m.Get("/clients", func(r render.Render) {

	})

	go m.Run()
}
