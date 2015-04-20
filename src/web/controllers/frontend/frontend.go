package frontend

import (
	"github.com/gorilla/context"
	"github.com/unrolled/render"
	"net/http"
	"web/helpers/mynegroni"
)

func Index(rw http.ResponseWriter, r *http.Request) {
	v := mynegroni.NewContainer(r)

	render := context.Get(r, "Render").(*render.Render)
	render.HTML(rw, http.StatusOK, "frontend/index", v)
}

func Login(rw http.ResponseWriter, r *http.Request) {
	v := mynegroni.NewContainer(r)

	render := context.Get(r, "Render").(*render.Render)
	render.HTML(rw, http.StatusOK, "frontend/login", v)
}

func Profile(rw http.ResponseWriter, r *http.Request) {
	v := mynegroni.NewContainer(r)

	render := context.Get(r, "Render").(*render.Render)
	render.HTML(rw, http.StatusOK, "frontend/profile", v)
}
