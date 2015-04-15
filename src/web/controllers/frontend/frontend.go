package frontend

import "net/http"
import "web/models/user"
import "github.com/martini-contrib/render"
import "github.com/martini-contrib/sessions"

var view map[string]interface{}

func init() {
	view = make(map[string]interface{})
}

func Index(r render.Render, profile user.Profiler) {
	view["profile"] = profile

	r.HTML(http.StatusOK, "frontend/index", view)
}

func Login(r render.Render, profile user.Profiler, s sessions.Session) {
	view["profile"] = profile
	view["errors"] = s.Flashes()

	r.HTML(http.StatusOK, "frontend/login", view)
}

func Profile(r render.Render, profile user.Profiler) {
	view["profile"] = profile

	r.HTML(http.StatusOK, "frontend/profile", view)
}
