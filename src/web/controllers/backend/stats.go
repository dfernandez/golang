package backend

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/unrolled/render"
	"mynegroni"
	"net/http"
)

func Stats(rw http.ResponseWriter, r *http.Request, render *render.Render, s sessions.Session, c *mynegroni.Content) {

	render.HTML(rw, http.StatusOK, "backend/stats", c)
}
