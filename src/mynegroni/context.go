package mynegroni

import (
	"github.com/gorilla/context"
	"net/http"
)

type Container interface {
	Initialize()
	Get(interface{}) interface{}
	Set(interface{}, interface{})
}

type Content struct {
	vars map[interface{}]interface{}
}

func NewContext(r *http.Request) *Content {
	c := new(Content)
	c.vars = context.GetAll(r)

	return c
}

func (c *Content) Get(key interface{}) interface{} {
	return c.vars[key]
}

func (c *Content) Set(key interface{}, val interface{}) {
	c.vars[key] = val
}
