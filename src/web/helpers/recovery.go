package helpers

import (
	"fmt"
	sessions "github.com/goincremental/negroni-sessions"
	"log"
	"mynegroni"
	"net/http"
	"os"
	"web/models/user"
)

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
