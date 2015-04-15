package main

import "github.com/go-martini/martini"
import "github.com/martini-contrib/render"
import "github.com/martini-contrib/sessions"
import "net/http"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import "web/controllers/frontend"
import "web/controllers/backend"
import "web/helpers/authentication"

var db *sql.DB

func main() {

	// Database
	db = getDB()
	defer db.Close()

	// Cookies
	store := sessions.NewCookieStore([]byte("myapp"))

	m := martini.Classic()

	m.Group("/", func(r martini.Router) {
		r.Get("", frontend.Index)
		r.Get("login", frontend.Login)
		r.Get("profile", authentication.LoginRequired, frontend.Profile)
	})

	m.Group("/admin", func(r martini.Router) {
		r.Get("", backend.Index)
	})

	m.NotFound(error404)

	m.Use(sessions.Sessions("session", store))
	m.Use(render.Renderer(render.Options{Layout: "layout"}))

	m.Use(authentication.Basic)
	m.Use(authentication.Google)
	m.Use(authentication.Facebook)

	m.Map(db)
	m.Run()
}

func error404(r render.Render) {
	r.HTML(http.StatusNotFound, "error404", nil)
}

func getDB() *sql.DB {
	var err error

	if db == nil {
		db, err = sql.Open("mysql", "username:password@tcp(hostname:3306)/database")
		if err != nil {
			panic(err)
		}
	}

	return db
}
