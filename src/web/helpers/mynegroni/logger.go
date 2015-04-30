package mynegroni

import (
	"github.com/codegangsta/negroni"
	"github.com/fatih/color"
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

const LOG_PANIC = "PANIC"
const LOG_ERROR = "ERROR"

func LogMessage(r *http.Request, errorType string, errorMessage string) {

	logTime := time.Now().Format("2006-01-02 15:04:01")

	f, err := os.OpenFile("log/messages", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("%s", err)
	}
	defer f.Close()

	remoteAddr := r.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr = r.RemoteAddr
	}

	l := log.New(os.Stdout, "[negroni] ", 0)

	if errorType == LOG_PANIC {
		color.Set(color.FgRed)
	} else {
		color.Set(color.FgYellow)
	}

	l.Printf("%-25s | %-7s | %s", remoteAddr, errorType, errorMessage)
	color.Unset()

	log.SetOutput(f)
	log.Printf("%s - %s - %s - %s - %s: %s", logTime, remoteAddr, r.Method, r.URL.Path, errorType, errorMessage)
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
	l.Printf("%-25s | %-7s | %-60s | %v %-25s | %12v", remoteAddr, r.Method, r.URL.Path, res.Status(), http.StatusText(res.Status()), time.Since(start))
}
