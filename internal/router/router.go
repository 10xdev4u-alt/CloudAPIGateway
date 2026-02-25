package router

import (
	"net/http"
)

type Router struct {
	routes []Route
}

type Route struct {
	Method  string
	Path    string
	Target  string
	Handler http.HandlerFunc
}

func New() *Router {
	return &Router{}
}

func (r *Router) AddRoute(method, path, target string, handler http.HandlerFunc) {
	r.routes = append(r.routes, Route{
		Method:  method,
		Path:    path,
		Target:  target,
		Handler: handler,
	})
}

func (r *Router) Match(method, path string) (Route, bool) {
	for _, route := range r.routes {
		if route.Method == method && route.Path == path {
			return route, true
		}
	}
	return Route{}, false
}
