package impl

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_T1(t *testing.T) {
	id := primitive.NewObjectID()
	t1 := time.Now().Unix()
	fmt.Println(byte(uint32(t1)))
	fmt.Println(id)
	type A struct {
		ID primitive.ObjectID `json:"id"`
	}
	a := A{
		ID: id,
	}
	msg, _ := json.Marshal(a)
	fmt.Println(string(msg))
	b := &A{}
	json.Unmarshal(msg, b)
	fmt.Println(b.ID)
}
