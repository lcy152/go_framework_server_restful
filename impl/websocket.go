package impl

import (
	"log"
	"net/http"

	framework "tumor_server/framework"
	message "tumor_server/message"
	"tumor_server/model"
	"tumor_server/service"

	"github.com/gorilla/websocket"
)

func WSBaseMessage(c *framework.Context) {
	var (
		wsConn  *websocket.Conn
		err     error
		conn    *framework.Connection
		message []byte
		upgrade websocket.Upgrader
	)
	sendChan := make(chan string)
	sc := service.GetContainerInstance()
	userInfo := GetContextUserInfo(c)
	url := service.GetMessageKey(userInfo.User.ID.String())
	sc.Dispatch.Subscribe(url, sendChan)
	defer sc.Dispatch.Unsubscribe(url, sendChan)
	if wsConn, err = upgrade.Upgrade(c.W, c.Req, nil); err != nil {
		return
	}
	if conn, err = framework.InitConnection(wsConn); err != nil {
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

func WSMessage(c *framework.Context) {
	defer PanicHandler(c)
	var (
		wsConn         *websocket.Conn
		err            error
		conn           *framework.Connection
		receiveMessage []byte
		upgrade        websocket.Upgrader
	)
	upgrade.CheckOrigin = func(r *http.Request) bool {
		return c.IsWebsocket()
	}
	sc := service.GetContainerInstance()
	token := c.GetParam("token")
	userInfo, err := service.TokenValidate(token, c.Req.Host)
	CheckHandler(err, message.ValidateError)
	url := service.GetMessageKey(userInfo.User.ID.String())
	log.Println("websocket open: " + url)

	sendChan := make(chan string)
	sc.Dispatch.Subscribe(url, sendChan)
	defer sc.Dispatch.Unsubscribe(url, sendChan)
	if wsConn, err = upgrade.Upgrade(c.W, c.Req, nil); err != nil {
		return
	}
	if conn, err = framework.InitConnection(wsConn); err != nil {
		goto ERR
	}
	go func() {
		for {
			select {
			case sendMsg := <-sendChan:
				if len(sendMsg) < service.ProtocalLength {
					if sendMsg == "ping" {
						conn.WriteMessage([]byte("pong"))
					} else {
						conn.WriteMessage([]byte(sendMsg))
					}
				} else {
					t, m := service.WsMessageDecode(sendMsg)
					response := model.WSResponse{
						Code:    200,
						Msg:     message.HttpOk,
						Flag:    t,
						Message: m,
					}
					conn.WriteMessage(response.ToJson())
				}
			}
		}
	}()
	for {
		if receiveMessage, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		log.Println(string(receiveMessage))
		sendChan <- string(receiveMessage)
	}

ERR:
	log.Println("websocket close: " + url)
	conn.Close()
}
