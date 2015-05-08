package backend

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/unrolled/render"
	"mynegroni"
	"net/http"
	"web/models/user"
)

func Profiles(rw http.ResponseWriter, r *http.Request, render *render.Render, s sessions.Session, c *mynegroni.Content) {

	profiles := user.GetProfiles()
	c.Set("profiles", profiles)

	render.HTML(rw, http.StatusOK, "backend/profiles", c)
}
