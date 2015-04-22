package mynegroni

import (
	"github.com/codegangsta/negroni"
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

func (c *Config) GetDatabase(env string) interface{} {
	return c.Database[env]
}

var Settings = func() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, "config", config)
		next(rw, r)
	}
}()
