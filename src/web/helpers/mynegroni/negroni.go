package mynegroni

import (
	"code.google.com/p/gcfg"
	"github.com/codegangsta/negroni"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"net/http"
	"os"
)

var config Config

func init() {
	// Config file
	path := os.Getenv("GOPATH") + "/cfg/app.gcfg"
	err := gcfg.ReadFileInto(&config, path)

	if err != nil {
		panic(err)
	}
}

func New() *negroni.Negroni {

	n := negroni.New()

	// Logger
	logger := NewLogger()

	// Recovery
	recovery := NewRecovery()

	// Render
	render := NewRender()

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
	n.Use(BasicOAuth)
	n.Use(GoogleOAuth)
	n.Use(FacebookOAuth)

	return n
}

func NotFound(rw http.ResponseWriter, r *http.Request) {
	v := NewContainer(r)

	renderer := NewRender()
	renderer.Render.HTML(rw, http.StatusNotFound, "error404", v)
}
