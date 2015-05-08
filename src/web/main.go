package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"mynegroni"
	"net/http"
	"os"
	"web/controllers/backend"
	"web/controllers/frontend"
	"web/helpers"
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

	// Panic recovery
	recovery := helpers.NewRecovery()
	n.Use(recovery)

	mainRouter := mux.NewRouter()
	mainRouter.NotFoundHandler = http.HandlerFunc(mynegroni.NotFound)

	// Frontend
	mainRouter.Handle("/", frontend.Controller(frontend.Index))
	mainRouter.Handle("/login", frontend.Controller(frontend.Login))
	mainRouter.Handle("/login/google/callback", frontend.Controller(frontend.LoginCallback))
	mainRouter.Handle("/login/facebook/callback", frontend.Controller(frontend.LoginCallback))

	// User profile
	mainRouter.Handle("/user", negroni.New(frontend.LoginRequired, negroni.Wrap(frontend.Controller(frontend.Dashboard))))
	mainRouter.Handle("/user/profile", negroni.New(frontend.LoginRequired, negroni.Wrap(frontend.Controller(frontend.Profile))))

	// Backend
	mainRouter.Handle("/backend", negroni.New(backend.LoginRequired, negroni.Wrap(backend.Controller(backend.Index))))
	mainRouter.Handle("/backend/profiles", negroni.New(backend.LoginRequired, negroni.Wrap(backend.Controller(backend.Profiles))))
	mainRouter.Handle("/backend/stats", negroni.New(backend.LoginRequired, negroni.Wrap(backend.Controller(backend.Stats))))

	// Routers
	n.UseHandler(mainRouter)

	// Run negroni run!
	n.Run(":3000")
}
