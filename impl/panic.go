package impl

import (
	"log"
	"runtime"
	framework "tumor_server/framework"
	"tumor_server/message"
	"tumor_server/model"
)

type NetError struct {
	Code int
	Msg  string
	Type string
}

func CheckHandler(err interface{}, msg string) {
	isError := false
	switch e := err.(type) {
	case error:
		isError = true
		log.Printf("error: %s\n", e.Error())
	case bool:
		isError = e
	}
	if isError {
		funcName, file, line, ok := runtime.Caller(1)
		funcNameStr := runtime.FuncForPC(funcName).Name()
		if ok {
			log.Printf("error: %s  %d  %s  %s\n", funcNameStr, line, msg, file)
			// funcName2, file2, line2, ok2 := runtime.Caller(2)
			// funcNameStr2 := runtime.FuncForPC(funcName2).Name()
			// if ok2 {
			// 	log.Printf("error: %s  %d  %s\n", funcNameStr2, line2, file2)
			// }
		} else {
			log.Printf("error: %s  getFunc: %s\n", msg, "fail")
		}
		myError := NetError{
			Msg:  msg,
			Code: 203,
			Type: "http",
		}
		panic(myError)
	}
}

func PanicHandler(c *framework.Context) {
	if r := recover(); r != nil {
		log.Println(r)
		response := model.HttpResponse{Code: 500}
		switch e := r.(type) {
		case NetError:
			if e.Type == "http" {
				response.Code = e.Code
				response.Msg = e.Msg
			}
		case string:
			response.Code = 500
			response.Msg = e
		}
		c.Error(response.ToJson())
	}
}

func HttpReponseHandler(c *framework.Context, data interface{}) {
	response := model.HttpResponse{Code: 200, Msg: message.HttpOk, Data: data}
	c.Error(response.ToJson())
}

func HttpReponseErrorHandler(c *framework.Context, msg string) {
	response := model.HttpResponse{Code: 200, Msg: msg}
	c.Error(response.ToJson())
}

func HttpReponseListHandler(c *framework.Context, count int64, data interface{}) {
	response := model.HttpResponse{Code: 200, Msg: message.HttpOk, Count: count, Data: data}
	c.Error(response.ToJson())
}

func HttpReponseExtraListHandler(c *framework.Context, count int64, extra string, data interface{}) {
	response := model.HttpResponse{Code: 200, Msg: message.HttpOk, Count: count, Flag: extra, Data: data}
	c.Error(response.ToJson())
}
