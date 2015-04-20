package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"web/controllers/frontend"
	mynegroni "web/helpers/negroni"
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

	// Routing
	router := mux.NewRouter()
	router.HandleFunc("/", frontend.Index)
	router.HandleFunc("/login", frontend.Login)

	// Routing for /profile
	profileRouter := mux.NewRouter()
	profileRouter.HandleFunc("/profile", frontend.Profile)

	// Login required middleware for /profile
	router.Handle("/profile", negroni.New(
		mynegroni.LoginRequired,
		negroni.Wrap(profileRouter),
	))

	router.NotFoundHandler = http.HandlerFunc(NotFound)

	// Routing
	n.UseHandler(router)

	// Run negroni run!
	n.Run(":3000")
}

func NotFound(rw http.ResponseWriter, r *http.Request) {
	v := mynegroni.NewContainer(r)

	renderer := mynegroni.NewRender()
	renderer.Render.HTML(rw, http.StatusNotFound, "error404", v)
}
