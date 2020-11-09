package model

import (
	"encoding/json"
	"time"
)

type HttpResponse struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Flag  string      `json:"flag"`
	Count int64       `json:"count"`
	Data  interface{} `json:"data"`
	Time  time.Time   `json:"time"`
}

func (r *HttpResponse) ToJson() []byte {
	r.Time = time.Now()
	res, err := json.Marshal(r)
	if err != nil {
		r.Msg = "json marshal error"
		r.Data = nil
		res, _ = json.Marshal(r)
	}
	return res
}

type WSResponse struct {
	Code    int       `json:"code"`
	Msg     string    `json:"msg"`
	Flag    string    `json:"flag"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func (r *WSResponse) ToJson() []byte {
	r.Time = time.Now()
	res, _ := json.Marshal(r)
	return res
}
