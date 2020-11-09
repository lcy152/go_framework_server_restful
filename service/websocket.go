package service

import (
	"strings"
)

const (
	MESSAGE_KEY = "message:"
)

var messageFlag = struct {
	Single      string
	Group       string
	Center      string
	Institution string
	Friend      string
}{
	Single:      "single",
	Group:       "group",
	Center:      "center",
	Institution: "institution",
	Friend:      "friend",
}

const (
	ProtocalLength = 20
)

func WsMessageDecode(sendMessage string) (string, string) {
	tag := ""
	value := ""
	if len(sendMessage) >= ProtocalLength {
		tag = sendMessage[:ProtocalLength]
		tag = strings.ReplaceAll(tag, ".", "")
		value = sendMessage[ProtocalLength:]
	}
	return tag, value
}

func WsMessageEncode(tag, sendMessage string) string {
	for i := len(tag); i < ProtocalLength; i++ {
		tag = tag + "."
	}
	return tag + sendMessage
}

func GetMessageKey(key string) string {
	return MESSAGE_KEY + key
}

func sendMessage(id, tag, value string) {
	sc := GetContainerInstance()
	message := WsMessageEncode(tag, value)
	sc.Dispatch.Post(GetMessageKey(id), message)
}

func SendSingleMessage(id, value string) {
	sendMessage(id, messageFlag.Single, value)
}
func SendGroupMessage(id, value string) {
	sendMessage(id, messageFlag.Group, value)
}
func SendCenterMessage(id, value string) {
	sendMessage(id, messageFlag.Center, value)
}
func SendInstitutionMessage(id, value string) {
	sendMessage(id, messageFlag.Institution, value)
}
func SendFriendMessage(id, value string) {
	sendMessage(id, messageFlag.Friend+value, value)
}
