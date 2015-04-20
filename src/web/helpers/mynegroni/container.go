package mynegroni

import (
	"encoding/gob"
	sessions "github.com/goincremental/negroni-sessions"
	"net/http"
	"web/models/user"
)

type Container interface {
	Initialize()
	Get(string) interface{}
	Set(string, interface{})
}

type Content struct {
	vars map[string]interface{}
}

func NewContainer(r *http.Request) *Content {
	c := new(Content)
	c.vars = make(map[string]interface{})

	s := sessions.GetSession(r)

	var profile user.Profile
	gob.Register(profile)

	p := s.Get("profile")
	c.Set("profile", p)

	return c
}

func (c *Content) Get(key string) interface{} {
	return c.vars[key]
}

func (c *Content) Set(key string, val interface{}) {
	c.vars[key] = val
}
