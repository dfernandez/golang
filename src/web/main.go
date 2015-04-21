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

	// Frontend
	mainRouter := mux.NewRouter()
	mainRouter.Handle("/", frontend.Controller(frontend.Index))
	mainRouter.Handle("/login", frontend.Controller(frontend.Login))
	mainRouter.Handle("/profile", negroni.New(
		frontend.LoginRequired,
		negroni.Wrap(frontend.Controller(frontend.Profile)),
	))

	// Backend
	routerBack := mux.NewRouter()
	routerBack.Handle("/backend", backend.Controller(backend.Index))

	// Routers
	mainRouter.Handle("/backend", routerBack)
	mainRouter.NotFoundHandler = http.HandlerFunc(mynegroni.NotFound)
	n.UseHandler(mainRouter)

	// Run negroni run!
	n.Run(":3000")
}
