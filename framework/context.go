package framework

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// HandlerFunc .
type HandlerFunc func(*Context)

// Context .
type Context struct {
	W        http.ResponseWriter
	Req      *http.Request
	Path     string
	Method   string
	Params   map[string]string
	body     []byte
	extra    []byte
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	c := &Context{
		W:      w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
	return c
}

// Next middleware handle success func
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return len(c.handlers) == c.index
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
func (c *Context) Abort() {
	c.index = len(c.handlers)
}

// AbortWithStatus calls `Abort()` and `Status()` internally.
func (c *Context) AbortWithStatus(code int) {
	c.Abort()
	c.SetStatusCode(code)
}

// AbortWithJSON calls `Abort()` and `JSON()` internally.
func (c *Context) AbortWithJSON(code int, obj interface{}) {
	c.Abort()
	c.JSON(code, obj)
}

// AbortWithString calls `Abort()` and `String()` internally.
func (c *Context) AbortWithString(code int, format string, values ...interface{}) {
	c.Abort()
	c.String(code, format, values...)
}

// GetURLParam get url param
func (c *Context) GetURLParam(key string) string {
	return c.Req.URL.Query().Get(key)
}

// GetFormValue get form param
func (c *Context) GetFormValue(key string) string {
	return c.Req.FormValue(key)
}

// GetBody parse body to []byte
func (c *Context) GetBody() []byte {
	if c.body == nil {
		b, _ := ioutil.ReadAll(c.Req.Body)
		c.body = b
		if c.body == nil {
			c.body = []byte{}
		}
		c.Req.Body.Close()
	}
	return c.body
}

// SetStatusCode set response
func (c *Context) SetStatusCode(code int) *Context {
	c.W.WriteHeader(code)
	return c
}

// SetHeader set response
func (c *Context) SetHeader(key string, value string) *Context {
	c.W.Header().Set(key, value)
	return c
}

// String formats according to a format specifier and put the resulting string in the response body.
// It also sets the Content-Type as "text/plain".
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatusCode(code)
	c.W.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json; charset=UTF-8")
	c.SetStatusCode(code)
	encoder := json.NewEncoder(c.W)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.W, err.Error(), 500)
	}
}

func (c *Context) requestHeader(key string) string {
	return c.Req.Header.Get(key)
}

// IsWebsocket returns true if the request headers indicate that a websocket
// handshake is being initiated by the client.
func (c *Context) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(c.requestHeader("Connection")), "upgrade") &&
		strings.EqualFold(c.requestHeader("Upgrade"), "websocket") {
		return true
	}
	return false
}

// Success response
func (c *Context) Success(res []byte) {
	c.W.Header().Set("content-Type", "application/json; charset=UTF-8")
	c.W.WriteHeader(http.StatusOK)
	c.W.Write(res)
}

// Error response
func (c *Context) Error(res []byte) {
	c.W.Header().Set("content-Type", "application/json; charset=UTF-8")
	c.W.WriteHeader(http.StatusNonAuthoritativeInfo)
	c.W.Write(res)
}

// Warn response
func (c *Context) Warn(res []byte) {
	c.W.Header().Set("content-Type", "application/json; charset=UTF-8")
	c.W.WriteHeader(http.StatusCreated)
	c.W.Write(res)
}

func (c *Context) ParseBody(data interface{}) bool {
	b := c.GetBody()
	err := json.Unmarshal(b, data)
	if err != nil {
		return false
	}
	return true
}

func (c *Context) ParseExtra(data interface{}) {
	json.Unmarshal(c.extra, data)
	return
}

func (c *Context) SetExtra(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	c.extra = b
	return nil
}

// GetAuthorization calls `Abort()` and `String()` internally.
func (c *Context) GetAuthorization() string {
	return c.Req.Header.Get("Authorization")
}

// AbortWithString calls `Abort()` and `String()` internally.
func (c *Context) GetParam(key string) string {
	if val, ok := c.Params[key]; ok {
		return val
	}
	return ""
}
