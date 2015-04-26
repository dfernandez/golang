package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"web/controllers/backend"
	"web/controllers/frontend"
	"web/helpers/mynegroni"
)

func init() {
	env := os.Getenv("ENV")

	if env == "" {
		panic("ENV not set")
	}

	fmt.Printf("[negroni] %s\n", env)
}

func main() {
	// Negroni
	n := mynegroni.New()

	mainRouter := mux.NewRouter()
	mainRouter.NotFoundHandler = http.HandlerFunc(mynegroni.NotFound)

	// Frontend
	mainRouter.Handle("/", frontend.Controller(frontend.Index))
	mainRouter.Handle("/login", frontend.Controller(frontend.Login))
	mainRouter.Handle("/login/google/callback", frontend.Controller(frontend.LoginCallback))
	mainRouter.Handle("/login/facebook/callback", frontend.Controller(frontend.LoginCallback))
	mainRouter.Handle("/user", negroni.New(
		frontend.LoginRequired,
		negroni.Wrap(frontend.Controller(frontend.Dashboard)),
	))
	mainRouter.Handle("/user/profile", negroni.New(
		frontend.LoginRequired,
		negroni.Wrap(frontend.Controller(frontend.Profile)),
	))

	// Backend
	mainRouter.Handle("/backend", negroni.New(
		backend.LoginRequired,
		negroni.Wrap(backend.Controller(backend.Index)),
	))
	mainRouter.Handle("/backend/profiles", negroni.New(
		backend.LoginRequired,
		negroni.Wrap(backend.Controller(backend.Profiles)),
	))

	// Routers
	n.UseHandler(mainRouter)

	// Run negroni run!
	n.Run(":3000")
}
