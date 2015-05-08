package frontend

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/unrolled/render"
	"mynegroni"
	"net/http"
)

func Dashboard(rw http.ResponseWriter, r *http.Request, render *render.Render, s sessions.Session, c *mynegroni.Content) {
	c.Set("template", "dashboard")

	render.HTML(rw, http.StatusOK, "frontend/dashboard", c)
}

func Profile(rw http.ResponseWriter, r *http.Request, render *render.Render, s sessions.Session, c *mynegroni.Content) {
	c.Set("template", "profile")

	render.HTML(rw, http.StatusOK, "frontend/profile", c)
}
