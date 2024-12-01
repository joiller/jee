package jee

import (
	"encoding/json"
	"net/http"
)

type H map[string]interface{}

type HandlerFunc func(c *Context)

type Context struct {
	Writer   http.ResponseWriter
	Req      *http.Request
	index    int
	handlers []HandlerFunc
	engine   *Engine
	params   map[string]interface{}
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	context := &Context{
		Writer: w,
		Req:    r,
		index:  -1,
	}
	return context
}

func (c *Context) Next() {
	c.index++
	l := len(c.handlers)
	for ; c.index < l; c.index++ {
		c.handlers[c.index](c)
	}
}

// if code=200, it's not needed to be called
func (c *Context) err(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) Fail(code int, err error) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) JSON(code int, data interface{}) {
	c.Writer.WriteHeader(code)
	c.Writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(c.Writer).Encode(data)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) Param(name string) interface{} {
	return c.params[name]
}
