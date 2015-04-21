package frontend

import (
	"github.com/codegangsta/negroni"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/gorilla/context"
	"github.com/unrolled/render"
	"net/http"
	"web/helpers/mynegroni"
)

var Controller = func() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		render := context.Get(r, "Render").(*render.Render)
		session := sessions.GetSession(r)
		container := mynegroni.NewContainer(r)

		switch r.URL.Path {

		case "/":
			Index(rw, r, render, session, container)

		case "/login":
			Login(rw, r, render, session, container)

		case "/profile":
			Profile(rw, r, render, session, container)

		default:
			next(rw, r)
		}
	}
}()

var Restricted = func() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		switch r.URL.Path {
		case "/profile":
			token := mynegroni.GetToken(r)

			if token == nil || !token.Valid() {
				http.Redirect(rw, r, "/login", http.StatusFound)
			} else {
				next(rw, r)
			}
		default:
			next(rw, r)
		}
	}
}()

var NotFound = func() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		render := context.Get(r, "Render").(*render.Render)
		container := mynegroni.NewContainer(r)

		render.HTML(rw, http.StatusNotFound, "error404", container)
	}
}()
