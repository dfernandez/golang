package mynegroni

import (
	"github.com/gorilla/context"
	"github.com/unrolled/render"
	"net/http"
	"os"
)

type MyRender struct {
	Render *render.Render
}

func NewRender() *MyRender {

	var options render.Options

	if os.Getenv("ENV") == "development" {
		options = render.Options{Layout: "layout", IsDevelopment: true}
	} else {
		options = render.Options{Layout: "layout"}
	}

	render := render.New(options)

	return &MyRender{render}
}

func (render MyRender) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// Attach the renderer
	context.Set(r, "Render", render.Render)

	// Call the next middleware handler
	next(rw, r)

}
