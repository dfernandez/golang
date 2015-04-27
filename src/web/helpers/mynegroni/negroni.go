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

	// Stores application config
	n.Use(Settings)

	// Logger
	logger := NewLogger()

	// Static
	static := negroni.NewStatic(http.Dir("public"))

	// Sessions
	var cookieOptions sessions.Options

	if os.Getenv("ENV") == "production" {
        cookieOptions.Path = "/"
        cookieOptions.MaxAge = 604800
		cookieOptions.Secure = true
		cookieOptions.HTTPOnly = true
	}

	cookie := cookiestore.New([]byte("colernio"))
	cookie.Options(cookieOptions)
	session := sessions.Sessions("colernio", cookie)

	// Render
	render := NewRender()

	// Recovery
	recovery := NewRecovery()

	// Middleware
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
