package negroni

import (
	"github.com/codegangsta/negroni"
	"log"
	"net/http"
	"os"
	"time"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "[negroni] ", 0)}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	remoteAddr := r.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr = r.RemoteAddr
	}

	res := rw.(negroni.ResponseWriter)
	l.Printf("%s %s %s %v %s in %v", remoteAddr, r.Method, r.URL.Path, res.Status(), http.StatusText(res.Status()), time.Since(start))
}
