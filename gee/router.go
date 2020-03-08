package gee

import (
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

// only one * is allowed
func parsePattern(pattern string) []string {
	var parts []string
	for _, part := range strings.Split(pattern, "/") {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	parts := parsePattern(pattern)
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[method+"-"+pattern] = handler
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

func (r *router) getRoute(method, path string) (*node, map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	searchParts := parsePattern(path)
	n := root.search(searchParts, 0)
	if n != nil {
		params := make(map[string]string)
		parts := parsePattern(n.pattern)
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

/* func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}

	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
} */
