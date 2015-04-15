package backend

import "net/http"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "github.com/martini-contrib/render"

var view map[string]interface{}

func init() {
	view = make(map[string]interface{})
}

func Index(req *http.Request, r render.Render, db *sql.DB) {
	r.HTML(http.StatusOK, "backend/index", view)
}
