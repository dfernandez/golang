package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
	"os"
	helpers "web/helpers/negroni"
)

func main() {

	// Negroni
	n := negroni.New()

	// Routing
	router := mux.NewRouter()
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		v := helpers.NewContainer(r)

		render := context.Get(r, "Render").(*render.Render)
		render.HTML(rw, http.StatusOK, "frontend/index", v)
	})
	router.HandleFunc("/login", func(rw http.ResponseWriter, r *http.Request) {
		v := helpers.NewContainer(r)

		render := context.Get(r, "Render").(*render.Render)
		render.HTML(rw, http.StatusOK, "frontend/login", v)
	})

	// Routing for /profile
	profileRouter := mux.NewRouter()
	profileRouter.HandleFunc("/profile", func(rw http.ResponseWriter, r *http.Request) {
		v := helpers.NewContainer(r)

		render := context.Get(r, "Render").(*render.Render)
		render.HTML(rw, http.StatusOK, "frontend/profile", v)
	})

	// Login required middleware for /profile
	router.Handle("/profile", negroni.New(
		helpers.LoginRequired,
		negroni.Wrap(profileRouter),
	))

	router.NotFoundHandler = http.HandlerFunc(NotFound)

	// Logger
	logger := negroni.NewLogger()

	// Recovery
	recovery := helpers.NewRecovery()

	// Render
	render := helpers.NewRender()

	// Static
	static := negroni.NewStatic(http.Dir("public"))

	// Sessions
	store := cookiestore.New([]byte("myapp"))
	session := sessions.Sessions("mysession", store)

	// Middleware
	n.Use(logger)
	n.Use(static)
	n.Use(session)
	n.Use(render)
	n.Use(recovery)

	// OAuth Authentication
	n.Use(helpers.BasicOAuth)
	n.Use(helpers.GoogleOAuth)
	n.Use(helpers.FacebookOAuth)

	// Routing
	n.UseHandler(router)

	// Run negroni run!
	fmt.Printf("[negroni] %s\n", os.Getenv("ENV"))
	n.Run(":3000")
}

func NotFound(rw http.ResponseWriter, r *http.Request) {
	v := helpers.NewContainer(r)

	renderer := helpers.NewRender()
	renderer.Render.HTML(rw, http.StatusNotFound, "error404", v)
}