package mynegroni

import (
	sessions "github.com/goincremental/negroni-sessions"
	"net/http"
)

func NotFound(rw http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	content := NewContext(r)

	GetUserProfile(content, session)

	renderer := NewRender()
	renderer.Render.HTML(rw, http.StatusNotFound, "error404", content)
}
