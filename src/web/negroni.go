package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"web/controllers/frontend"
	helpers "web/helpers/negroni"
)

func main() {

	// Negroni
	n := negroni.New()

	// Routing
	router := mux.NewRouter()
	router.HandleFunc("/", frontend.Index)
	router.HandleFunc("/login", frontend.Login)

	// Routing for /profile
	profileRouter := mux.NewRouter()
	profileRouter.HandleFunc("/profile", frontend.Profile)

	// Login required middleware for /profile
	router.Handle("/profile", negroni.New(
		helpers.LoginRequired,
		negroni.Wrap(profileRouter),
	))

	router.NotFoundHandler = http.HandlerFunc(NotFound)

	// Logger
	logger := helpers.NewLogger()

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
