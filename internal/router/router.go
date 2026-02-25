package router

import (
	"net/http"
)

type Router struct {
	routes []Route
}

type Route struct {
	Path    string
	Target  string
	Handler http.HandlerFunc
}

func New() *Router {
	return &Router{}
}

func (r *Router) AddRoute(path, target string, handler http.HandlerFunc) {
	r.routes = append(r.routes, Route{
		Path:    path,
		Target:  target,
		Handler: handler,
	})
}

func (r *Router) Match(path string) (Route, bool) {
	for _, route := range r.routes {
		if route.Path == path {
			return route, true
		}
	}
	return Route{}, false
}
