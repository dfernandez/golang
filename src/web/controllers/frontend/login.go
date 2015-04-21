package frontend

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/unrolled/render"
	"net/http"
	"web/helpers/mynegroni"
)

func Login(rw http.ResponseWriter, r *http.Request, render *render.Render, s sessions.Session, c *mynegroni.Content) {

	render.HTML(rw, http.StatusOK, "frontend/login", c)
}
