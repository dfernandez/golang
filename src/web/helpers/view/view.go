package view

type Container interface {
	Initialize()
	Get(string) interface{}
	Set(string, interface{})
}

type Content struct {
	vars map[string]interface{}
}

func (c *Content) Initialize() {
	c.vars = make(map[string]interface{})
}

func (c *Content) Get(key string) interface{} {
	return c.vars[key]
}

func (c *Content) Set(key string, val interface{}) {
	c.vars[key] = val
}
