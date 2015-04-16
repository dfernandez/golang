package main

import "github.com/go-martini/martini"
import "github.com/martini-contrib/render"
import "github.com/martini-contrib/sessions"
import "net/http"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "code.google.com/p/gcfg"

import "web/controllers/frontend"
import "web/controllers/backend"
import "web/helpers/authentication"
import "web/helpers/configuration"
import "os"

var db *sql.DB

var config configuration.Application

func init() {
	// Config file
	path := os.Getenv("GOPATH") + "/cfg/app.gcfg"
	err := gcfg.ReadFileInto(&config, path)

	if err != nil {
		panic(err)
	}
}

func main() {

	// Database
	db = getDB()
	defer db.Close()

	// Cookies
	store := sessions.NewCookieStore([]byte("myapp"))

	// Start martini
	m := martini.Classic()

	// Routing frontend
	m.Group("/", func(r martini.Router) {
		r.Get("", frontend.Index)
		r.Get("login", frontend.Login)
		r.Get("profile", authentication.LoginRequired, frontend.Profile)
	})

	// Routing backend
	m.Group("/admin", func(r martini.Router) {
		r.Get("", backend.Index)
	})

	// 404 Handler
	m.NotFound(error404)

	// Sessions & layout
	m.Use(sessions.Sessions("session", store))
	m.Use(render.Renderer(render.Options{Layout: "layout"}))

	// Authentication
	m.Use(authentication.Basic)
	m.Use(authentication.Google)
	m.Use(authentication.Facebook)

	// Dependency injection
	m.Map(db)
	m.Map(config)
	m.Run()
}

func error404(r render.Render) {
	r.HTML(http.StatusNotFound, "error404", nil)
}

func getDB() *sql.DB {
	var err error

	if db == nil {
		db, err = sql.Open("mysql", config.Mysql[martini.Env].Dns)
		if err != nil {
			panic(err)
		}
	}

	return db
}
