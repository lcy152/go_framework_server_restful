package framework

import (
	"log"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

func Test_ServerStart(t *testing.T) {
	s := NewEngine()
	GetContainerInstance("")
	// http
	s.AddMiddleware("", BeforeMiddleware1)
	s.AddMiddleware("", BeforeMiddleware2)
	s.AddAfterMiddleware("", AfterMiddleware1)
	s.AddAfterMiddleware("", AfterMiddleware2)
	s.GET("/login/:id/:password", Login)

	// websocket
	s.AddWSMiddleware("", WSBeforeMiddleware1)
	s.AddWSAfterMiddleware("", WSAfterMiddleware1)
	s.WS("/ws", WSTest)

	// file
	s.Static("../dist")
	s.FsFile("/fs", "../fs")
	s.Run(":2325")
}

var instance *Container
var once sync.Once

func GetExistContainerInstance() *Container {
	return instance
}

func GetContainerInstance(configPath string) *Container {
	once.Do(
		func() {
			instance = NewContainer(configPath)
		},
	)
	return instance
}

func NewContainer(configPath string) *Container {
	container := new(Container)
	return container
}

type Container struct {
	DB string
}

func Login(c *Context) {
	sc := GetExistContainerInstance()
	log.Print("Login", sc.DB)
	c.Success(nil)
}

func BeforeMiddleware1(c *Context) {
	log.Print("BeforeMiddleware1")
	c.Next()
}

func BeforeMiddleware2(c *Context) {
	log.Print("BeforeMiddleware2")
	c.Next()
}

func AfterMiddleware1(c *Context) {
	log.Print("AfterMiddleware1")
	c.Next()
}

func AfterMiddleware2(c *Context) {
	log.Print("AfterMiddleware2")
	c.Next()
}

func WSBeforeMiddleware1(c *Context) {
	log.Print("WSBeforeMiddleware1")
	c.Next()
}

func WSAfterMiddleware1(c *Context) {
	log.Print("WSAfterMiddleware1")
	c.Next()
}

func WSTest(c *Context) {
	var (
		wsConn  *websocket.Conn
		err     error
		conn    *Connection
		message []byte
		upgrade websocket.Upgrader
	)
	sendChan := make(chan string)
	if wsConn, err = upgrade.Upgrade(c.W, c.Req, nil); err != nil {
		return
	}
	if conn, err = InitConnection(wsConn); err != nil {
		goto ERR
	}

	go func() {
		for {
			select {
			case message := <-sendChan:
				conn.WriteMessage([]byte(message))
			}
		}
	}()

	for {
		if message, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		log.Println(string(message))
		sendChan <- "hello client"
	}

ERR:
	log.Println("close")
	conn.Close()
}
