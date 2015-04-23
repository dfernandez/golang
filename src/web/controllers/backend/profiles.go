package backend

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/unrolled/render"
	"net/http"
	"os"
	"web/helpers/mynegroni"
	"web/models/user"
)

func Profiles(rw http.ResponseWriter, r *http.Request, render *render.Render, s sessions.Session, c *mynegroni.Content) {

	config := c.Get("config").(mynegroni.Config)
	db := config.GetDatabase(os.Getenv("ENV"))

	profiles := user.GetProfiles(db)
	c.Set("profiles", profiles)

	render.HTML(rw, http.StatusOK, "backend/profiles", c)
}
