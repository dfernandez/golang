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
import "web/models/user"

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
	m := martini.New()

	// Logger.
	m.Use(martini.Logger())

	// Default layout.
	m.Use(render.Renderer(render.Options{Layout: "layout"}))

	// Statics.
	m.Use(martini.Static("public"))

	// Sessions.
	m.Use(sessions.Sessions("session", store))

	// OAuth Authentication
	m.Use(authentication.Basic)
	m.Use(authentication.Google)
	m.Use(authentication.Facebook)

	// Recovery. Needs to be after render, logger and oauth.
	m.Use(error500())

	// Routing
	r := martini.NewRouter()

	// Routing frontend
	r.Group("/", func(r martini.Router) {
		r.Get("", frontend.Index)
		r.Get("login", frontend.Login)
		r.Get("profile", authentication.LoginRequired, frontend.Profile)
	})

	// Routing backend
	r.Group("/admin", func(r martini.Router) {
		r.Get("", backend.Index)
	})

	// 404 Handler
	r.NotFound(error404)

	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)

	// Dependency injection
	m.Map(db)
	m.Map(config)
	m.Run()
}

// Recover from application panics.
func error500() martini.Handler {
	return func(c martini.Context, r render.Render, profile user.Profiler) {
		view := make(map[string]interface{})
		view["profile"] = profile

		defer func() {
			if err := recover(); err != nil {
				view["panic"] = err
				r.HTML(http.StatusInternalServerError, "error500", view)
			}
		}()

		c.Next()
	}
}

// Handler for 404 pages.
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
