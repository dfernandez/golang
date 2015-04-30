package mynegroni

import (
	"github.com/fatih/color"
	sessions "github.com/goincremental/negroni-sessions"
	"log"
	"net/http"
	"os"
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
	content := NewContext(r)

	GetUserProfile(content, session)

	defer func() {
		if err := recover(); err != nil {
			remoteAddr := r.Header.Get("X-Forwarded-For")
			if remoteAddr == "" {
				remoteAddr = r.RemoteAddr
			}

			color.Set(color.FgRed)
			rec.Logger.Printf("%-25s | %-7s | %s", remoteAddr, "PANIC", err)
			color.Unset()

			renderer := NewRender()
			renderer.Render.HTML(rw, http.StatusInternalServerError, "error500", content)
		}
	}()

	next(rw, r)
}
