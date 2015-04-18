package negroni

import (
	"github.com/gorilla/context"
	"github.com/unrolled/render"
	"net/http"
)

type MyRender struct {
	Render *render.Render
}

func NewRender() *MyRender {

	render := render.New(render.Options{Layout: "layout"})

	return &MyRender{render}
}

func NewRenderError() *MyRender {

	render := render.New(render.Options{Layout: ""})

	return &MyRender{render}
}

func (render MyRender) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// Attach the renderer
	context.Set(r, "Render", render.Render)

	// Call the next middleware handler
	next(rw, r)

}
