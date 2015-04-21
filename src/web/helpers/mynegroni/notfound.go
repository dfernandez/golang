package mynegroni

import (
	"net/http"
)

func NotFound(rw http.ResponseWriter, r *http.Request) {
	v := NewContainer(r)

	renderer := NewRender()
	renderer.Render.HTML(rw, http.StatusNotFound, "error404", v)
}
