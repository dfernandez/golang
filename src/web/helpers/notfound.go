package helpers

import (
	sessions "github.com/goincremental/negroni-sessions"
	"mynegroni"
	"net/http"
	"web/models/user"
)

func NotFound(rw http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	content := mynegroni.NewContext(r)

	if session.Get("profile") != nil {
		p := session.Get("profile").(user.Profile)

		content.Set("profile", &p)
	}

	renderer := mynegroni.NewRender()
	renderer.Render.HTML(rw, http.StatusNotFound, "error404", content)
}
