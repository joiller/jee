package jee

import (
	"net/http"
	"path"
	"strings"
)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix      string
	engine      *Engine
	middlewares []HandlerFunc
}

func New() *Engine {
	e := &Engine{
		router: newRouter(),
	}
	e.RouterGroup = &RouterGroup{
		engine: e,
	}
	e.groups = append(e.groups, e.RouterGroup)
	e.engine = e
	return e
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	r2 := &RouterGroup{
		prefix: r.prefix + prefix,
		engine: r.engine,
	}
	r2.engine.groups = append(r2.engine.groups, r2)
	return r2
}

func (r *RouterGroup) Use(handlers ...HandlerFunc) {
	r.middlewares = append(r.middlewares, handlers...)
}

func (r *RouterGroup) Static(route string, dir string) {
	handler := func(c *Context) {
		h := http.StripPrefix(route, http.FileServer(http.Dir(dir)))
		h.ServeHTTP(c.Writer, c.Req)
	}
	r.GET(path.Join(route, "/*filepath"), handler)
}

func (r *RouterGroup) GET(pattern string, handler HandlerFunc) {
	r.engine.router.addRoute("GET", pattern, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := newContext(w, r)
	context.engine = e
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(r.RequestURI, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	context.handlers = middlewares
	e.router.handle(context)
}
