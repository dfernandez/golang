package frontend

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/gorilla/context"
	"github.com/unrolled/render"
	"net/http"
	"web/helpers/mynegroni"
	"web/models/user"
)

func Login(rw http.ResponseWriter, r *http.Request, render *render.Render, s sessions.Session, c *mynegroni.Content) {

	c.Set("errors", s.Flashes("oauth"))

	render.HTML(rw, http.StatusOK, "frontend/login", c)
}

func LoginCallback(rw http.ResponseWriter, r *http.Request, render *render.Render, s sessions.Session, c *mynegroni.Content) {

	if context.Get(r, "oauth_profile") == nil {
		return
	}

	p := context.Get(r, "oauth_profile").(mynegroni.OauthProfile)

	profile := &user.Profile{Name: p.Name, Email: p.Email, Gender: p.Gender, Profile: p.Profile, Picture: p.Picture}
	profile.Upsert()

	s.Set("profile", profile)

	http.Redirect(rw, r, "/user", http.StatusFound)
}
