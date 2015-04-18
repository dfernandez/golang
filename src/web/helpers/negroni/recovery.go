package negroni

import (
	"log"
	"net/http"
	"os"
)

// Recovery is a Negroni middleware that recovers from any panics and writes a 500 if there was one.
type MyRecovery struct {
	Logger     *log.Logger
	PrintStack bool
	StackAll   bool
	StackSize  int
}

// NewRecovery returns a new instance of Recovery
func NewRecovery() *MyRecovery {
	return &MyRecovery{
		Logger:     log.New(os.Stdout, "[negroni] ", 0),
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

func (rec *MyRecovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			// rw.WriteHeader(http.StatusInternalServerError)
			// stack := make([]byte, rec.StackSize)
			// stack = stack[:runtime.Stack(stack, rec.StackAll)]

			rec.Logger.Printf("PANIC: %s\n", err)

			renderer := NewRenderError()
			renderer.Render.HTML(rw, http.StatusInternalServerError, "error500", nil)
		}
	}()

	next(rw, r)
}
