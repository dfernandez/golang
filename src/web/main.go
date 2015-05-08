package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"web/controllers/backend"
	"web/controllers/frontend"
	"web/helpers/mynegroni"
	"web/models/user"
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

	recovery := NewRecovery()
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

// Recovery is a Negroni middleware that recovers from any panics and writes a 500 if there was one.
type MyRecovery struct {
	Logger *log.Logger
}

// NewRecovery returns a new instance of Recovery
func NewRecovery() *MyRecovery {
	return &MyRecovery{
		Logger: log.New(os.Stdout, "[negroni] ", 0),
	}
}

func (rec *MyRecovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session := sessions.GetSession(r)
	content := mynegroni.NewContext(r)

	mynegroni.GetUserProfile(content, session)

	if session.Get("profile") != nil {
		p := session.Get("profile").(user.Profile)

		content.Set("profile", &p)
	}

	defer func() {
		if err := recover(); err != nil {

			mynegroni.LogMessage(r, mynegroni.LOG_PANIC, fmt.Sprintf("%s", err))

			renderer := mynegroni.NewRender()
			renderer.Render.HTML(rw, http.StatusInternalServerError, "error500", content)
		}
	}()

	next(rw, r)
}
