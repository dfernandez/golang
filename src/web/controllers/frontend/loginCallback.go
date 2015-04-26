package frontend

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/gorilla/context"
	"github.com/unrolled/render"
	"net/http"
	"os"
	"web/helpers/mynegroni"
	"web/models/user"
)

func LoginCallback(rw http.ResponseWriter, r *http.Request, render *render.Render, s sessions.Session, c *mynegroni.Content) {

	if context.Get(r, "oauth_profile") == nil {
		return
	}

	config := c.Get("config").(mynegroni.Config)
	db := config.GetDatabase(os.Getenv("ENV"))

	p := context.Get(r, "oauth_profile").(mynegroni.OauthProfile)

	profile := &user.Profile{Name: p.Name, Email: p.Email, Profile: p.Profile, Picture: p.Picture}
	profile.Upsert(db)

	s.Set("profile", profile)

	http.Redirect(rw, r, "/user", http.StatusFound)
}
