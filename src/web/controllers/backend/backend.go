package backend

import (
	_ "github.com/codegangsta/negroni"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/gorilla/context"
	"github.com/unrolled/render"
	"net/http"
	"web/helpers/mynegroni"
)

func Controller(action func(http.ResponseWriter, *http.Request, *render.Render, sessions.Session, *mynegroni.Content)) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		render := context.Get(r, "Render").(*render.Render)
		session := sessions.GetSession(r)
		context := mynegroni.NewContext(r)

		action(rw, r, render, session, context)
	})
}
