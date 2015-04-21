package mynegroni

import (
	"code.google.com/p/gcfg"
	"github.com/codegangsta/negroni"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/phyber/negroni-gzip/gzip"
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

	// Static
	static := negroni.NewStatic(http.Dir("public"))

	// Sessions
	session := sessions.Sessions("mysession", cookiestore.New([]byte("myapp")))

	// Render
	render := NewRender()

	// Recovery
	recovery := NewRecovery()

	// Middleware
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.Use(logger)
	n.Use(static)
	n.Use(session)
	n.Use(render)

	// OAuth Authentication
	n.Use(BasicOAuth)
	n.Use(GoogleOAuth)
	n.Use(FacebookOAuth)

	// Recovery
	n.Use(recovery)

	return n
}
