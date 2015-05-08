package mynegroni

import (
	"github.com/codegangsta/negroni"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"net/http"
	"os"
)

func New() *negroni.Negroni {

	n := negroni.New()

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
	} else {
		cookieOptions.Path = "/"
	}

	cookie := cookiestore.New([]byte("colernio"))
	cookie.Options(cookieOptions)
	session := sessions.Sessions("colernio", cookie)

	// Render
	render := NewRender()

	// Middleware
	n.Use(logger)
	n.Use(static)
	n.Use(session)
	n.Use(render)

	// OAuth Authentication
	n.Use(BasicOAuth)
	n.Use(GoogleOAuth)
	n.Use(FacebookOAuth)

	return n
}
