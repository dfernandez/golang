package frontend

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/unrolled/render"
	"net/http"
	"web/helpers/mynegroni"
)

func Index(rw http.ResponseWriter, r *http.Request, render *render.Render, s sessions.Session, c *mynegroni.Content) {

	panic("hola")

	render.HTML(rw, http.StatusOK, "frontend/index", c)
}
