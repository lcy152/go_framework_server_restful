package model

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func Test_Trie(t *testing.T) {
	type m struct {
		Description string    `json:"description" bson:"description"`
		InTime      time.Time `json:"in_time" bson:"in_time"`
	}
	s := m{
		Description: "ss",
		InTime:      time.Now(),
	}
	b, err := json.Marshal(s)
	str := string(b)
	log.Println(err)
	log.Println(str)
	n := &m{}
	err = json.Unmarshal([]byte(`{"description":"ss","in_time":"2020-11-05T15:07:21.3360361+08:00"}`), n)
	log.Println(n)
}
