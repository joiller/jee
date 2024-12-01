package jee

import (
	"net/http"
	"testing"
)

func newRouterTest() *router {
	return newRouter()
}

func TestRoute(t *testing.T) {
	r := newRouterTest()
	r.addRoute("GET", "/hello/:name", func(c *Context) {
		c.JSON(http.StatusOK, H{
			"name": c.Param("name"),
		})
	})
	if f, _ := r.getRoute("GET", "/hello/joiller"); f == nil {
		t.Fail()
	}
}
