package gee

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

type (
	// RouterGroup implement group routing
	RouterGroup struct {
		prefix      string
		parent      *RouterGroup // support nesting route
		engine      *Engine
		middlewares []HandlerFunc
	}

	// Engine implement the interface of ServeHTTP
	Engine struct {
		*RouterGroup
		*router
		groups []*RouterGroup
	}
)

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create new router group
func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix,
		parent: rg,
		engine: rg.engine,
	}
	rg.engine.groups = append(rg.engine.groups, newGroup)
	return newGroup
}

// Use for adding middleware to the route
func (rg *RouterGroup) Use(middewares ...HandlerFunc) {
	rg.middlewares = append(rg.middlewares, middewares...)
}

func (rg *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := rg.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	rg.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (rg *RouterGroup) GET(pattern string, handler HandlerFunc) {
	rg.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (rg *RouterGroup) POST(pattern string, handler HandlerFunc) {
	rg.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

// ServeHTTP defines the method to serve http request
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	e.router.handle(c)
}
