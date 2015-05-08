package frontend

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/unrolled/render"
	"mynegroni"
	"net/http"
)

func Index(rw http.ResponseWriter, r *http.Request, render *render.Render, s sessions.Session, c *mynegroni.Content) {

	render.HTML(rw, http.StatusOK, "frontend/index", c)
}
