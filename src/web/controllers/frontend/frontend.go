package frontend

import (
	"net/http"
	"github.com/unrolled/render"
	"github.com/gorilla/context"
	helpers "web/helpers/negroni"
)	

func Index(rw http.ResponseWriter, r *http.Request) {
	v := helpers.NewContainer(r)

	render := context.Get(r, "Render").(*render.Render)
	render.HTML(rw, http.StatusOK, "frontend/index", v)
}

func Login(rw http.ResponseWriter, r *http.Request) {
	v := helpers.NewContainer(r)

	render := context.Get(r, "Render").(*render.Render)
	render.HTML(rw, http.StatusOK, "frontend/login", v)
}

func Profile(rw http.ResponseWriter, r *http.Request) {
	v := helpers.NewContainer(r)

	render := context.Get(r, "Render").(*render.Render)
	render.HTML(rw, http.StatusOK, "frontend/profile", v)
}
