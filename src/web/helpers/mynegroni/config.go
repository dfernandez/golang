package mynegroni

import (
	"database/sql"
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"net/http"
)

type Configurator interface {
	GetKey(string) interface{}
}

type Config struct {
	OAuth map[string]*struct {
		ClientId string
		SecretId string
	}
	Database map[string]*struct {
		Connector string
		Dns       string
	}
	Domain map[string]*struct {
		Url string
	}
}

func (c *Config) GetDatabase(env string) *sql.DB {
	dbConfig := c.Database[env]

	db, err := sql.Open(dbConfig.Connector, dbConfig.Dns)

	if err != nil {
		panic(err)
	}

	return db
}

var Settings = func() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, "config", config)
		next(rw, r)
	}
}()
