package backend

import (
	"encoding/gob"
	"fmt"
	"github.com/codegangsta/negroni"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/gorilla/context"
	"github.com/unrolled/render"
	"mynegroni"
	"net/http"
	"web/models/user"
)

func Controller(action func(http.ResponseWriter, *http.Request, *render.Render, sessions.Session, *mynegroni.Content)) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		render := context.Get(r, "Render").(*render.Render)
		session := sessions.GetSession(r)
		content := mynegroni.NewContext(r)

		userProfile(content, session)

		content.Set("backend", true)

		action(rw, r, render, session, content)
	})
}

var LoginRequired = func() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		token := mynegroni.GetToken(r)

		if token == nil || !token.Valid() {
			http.Redirect(rw, r, "/login", http.StatusFound)
		} else {

			session := sessions.GetSession(r)
			content := mynegroni.NewContext(r)
			profile := session.Get("profile").(user.Profile)

			if !profile.IsAdmin() {
				userProfile(content, session)

				mynegroni.LogMessage(r, mynegroni.LOG_ERROR, fmt.Sprintf("UNAUTHORIZED ACCESS: %s", profile.Email))

				renderer := mynegroni.NewRender()
				renderer.Render.HTML(rw, http.StatusForbidden, "error403", content)
				return
			}

			next(rw, r)
		}
	}
}()

// Set user profile in context
func userProfile(content *mynegroni.Content, session sessions.Session) {
	if session.Get("profile") == nil {
		return
	}

	p := session.Get("profile").(user.Profile)

	content.Set("profile", &p)
}

func init() {
	var userProfile user.Profile
	gob.Register(userProfile)
}
