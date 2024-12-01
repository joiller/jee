package jee

import (
	"errors"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	n := r.roots[method]
	n.addNode(pattern, splitPattern(pattern), 0)
	r.handlers[method+"-"+pattern] = handler
}

func (r *router) getRoute(method string, path string) (HandlerFunc, map[string]interface{}) {
	if _, ok := r.roots[method]; !ok {
		return nil, nil
	}
	paths := splitPattern(path)
	searchNode := r.roots[method].searchNode(paths, 0)
	if searchNode != nil && searchNode.pattern != "" {
		pattern := splitPattern(searchNode.pattern)
		params := map[string]interface{}{}
		for i, p := range pattern {
			if p[0] == ':' {
				params[p[1:]] = paths[i]
			} else if p[0] == '*' && len(pattern) > i+1 {
				params[p[1:]] = paths[i][1:] + strings.Join(pattern[i+1:], "/")
				break
			}
		}
		return r.handlers[method+"-"+searchNode.pattern], params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	route, params := r.getRoute(c.Req.Method, c.Req.RequestURI)
	if route != nil {
		c.params = params
		c.handlers = append(c.handlers, route)
	} else {
		c.Fail(http.StatusNotFound, errors.New("route not found"))
	}
	c.Next()
}
