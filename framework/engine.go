package framework

import (
	"net/http"
	"strings"
)

const (
	// WSROOT .
	WSROOT = "WS"
)

// Engine server.
type Engine struct {
	node          *node
	fsHandler     map[string]http.Handler
	staticHandler http.Handler
}

// New server.
func NewEngine() *Engine {
	e := &Engine{node: &node{}, fsHandler: make(map[string]http.Handler)}
	return e
}

// AddMiddleware registers route middleware.
func (e *Engine) AddMiddleware(pattern string, handler HandlerFunc) {
	for _, method := range []string{"GET", "POST", "PUT", "DELETE"} {
		parts := parsePattern(pattern)
		parts = append([]string{method}, parts...)
		e.node.insertMiddleware(parts, handler, nil)
	}
}

// AddAfterMiddleware registers route middleware.
func (e *Engine) AddAfterMiddleware(pattern string, handler HandlerFunc) {
	for _, method := range []string{"GET", "POST", "PUT", "DELETE"} {
		parts := parsePattern(pattern)
		parts = append([]string{method}, parts...)
		e.node.insertMiddleware(parts, nil, handler)
	}
}

// AddWSMiddleware registers route middleware.
func (e *Engine) AddWSMiddleware(pattern string, handler HandlerFunc) {
	for _, method := range []string{WSROOT} {
		parts := parsePattern(pattern)
		parts = append([]string{method}, parts...)
		e.node.insertMiddleware(parts, handler, nil)
	}
}

// AddWSAfterMiddleware registers route middleware.
func (e *Engine) AddWSAfterMiddleware(pattern string, handler HandlerFunc) {
	for _, method := range []string{WSROOT} {
		parts := parsePattern(pattern)
		parts = append([]string{method}, parts...)
		e.node.insertMiddleware(parts, nil, handler)
	}
}

// GET registers get handler.
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	parts = append([]string{"GET"}, parts...)
	e.node.insert(parts, pattern, handler)
}

// POST registers post handler.
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	parts = append([]string{"POST"}, parts...)
	e.node.insert(parts, pattern, handler)
}

// PUT registers put handler.
func (e *Engine) PUT(pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	parts = append([]string{"PUT"}, parts...)
	e.node.insert(parts, pattern, handler)
}

// DELETE registers delete handler.
func (e *Engine) DELETE(pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	parts = append([]string{"DELETE"}, parts...)
	e.node.insert(parts, pattern, handler)
}

// WS registers get handler.
func (e *Engine) WS(pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	parts = append([]string{WSROOT}, parts...)
	e.node.insert(parts, pattern, handler)
}

// Static registers a single route in order to serve a single file of the local filesystem.
func (e *Engine) Static(folder string) {
	e.staticHandler = http.FileServer(http.Dir(folder))
}

// FsFile registers a single route in order to serve a single file of the local filesystem.
func (e *Engine) FsFile(prefix string, folder string) {
	e.fsHandler[prefix] = http.StripPrefix(prefix, http.FileServer(http.Dir(folder)))
}

// Run start server.
func (e *Engine) Run(addr string) (err error) {
	server := http.Server{
		Addr:    addr,
		Handler: e,
	}
	return server.ListenAndServe()
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	parts := parsePattern(c.Path)
	method := c.Method
	if c.IsWebsocket() {
		method = WSROOT
	}
	methodParts := append([]string{method}, parts...)
	beforeMiddleware := []HandlerFunc{}
	afterMiddleware := []HandlerFunc{}
	searchNode, beforeMiddleware, afterMiddleware := e.node.searchMiddleware(methodParts, beforeMiddleware, afterMiddleware)
	if searchNode != nil && searchNode.isLeaf {
		c.Params = make(map[string]string)
		for i, key := range searchNode.paramKey {
			if i < len(searchNode.paramValue) {
				c.Params[key] = searchNode.paramValue[i]
			}
		}
		c.handlers = append(beforeMiddleware, searchNode.handler)
		c.handlers = append(c.handlers, afterMiddleware...)
		c.Next()
	} else {
		ok := false
		for k, h := range e.fsHandler {
			if strings.Contains(k, parts[0]) {
				h.ServeHTTP(w, req)
				ok = true
				break
			}
		}
		if !ok && e.staticHandler != nil {
			e.staticHandler.ServeHTTP(w, req)
		}
	}
}
