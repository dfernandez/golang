package frontend

import "net/http"
import "web/models/user"
import "github.com/martini-contrib/render"
import "github.com/martini-contrib/sessions"
import "web/helpers/view"

func Index(r render.Render, profile user.Profiler, v view.Container) {
	v.Set("profile", profile)

	r.HTML(http.StatusOK, "frontend/index", v)
}

func Login(r render.Render, profile user.Profiler, s sessions.Session, v view.Container) {
	v.Set("profile", profile)
	v.Set("errors", s.Flashes())

	r.HTML(http.StatusOK, "frontend/login", v)
}

func Profile(r render.Render, profile user.Profiler, v view.Container) {
	v.Set("profile", profile)

	r.HTML(http.StatusOK, "frontend/profile", v)
}
