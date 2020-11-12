package service

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"time"
	"tumor_server/message"
)

const (
	ShortMessageTAG = "short_message:"
)

type ShortMessage struct {
	Phone string `json:"phone"`
	Time  int64  `json:"time"`
	Count int64  `json:"count"`
	Code  string `json:"code"`
}

func GenerateCode(number int) string {
	base := 9
	for i := 1; i < number; i++ {
		base = base*10 + 9
	}
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(base)
	return strconv.Itoa(code)
}

func NewShortMessageKey(key string) string {
	return ShortMessageTAG + key
}

func ShortMessageValidate(phone string, code string) bool {
	sc := GetContainerInstance()
	if code == "" {
		return false
	}
	codeInfo, err := GetRedisShortMessage(phone)
	if err != nil {
		return false
	}
	tn := time.Now().Unix()
	if (tn - codeInfo.Time) > int64(sc.Config.ShortMessageInvalidTime) {
		return false
	}
	if codeInfo.Code == code {
		return true
	}
	return false
}

func AddRedisShortMessage(info *ShortMessage) error {
	sc := GetContainerInstance()
	message, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = sc.RedisService.AddKey(NewShortMessageKey(info.Phone), string(message))
	if err != nil {
		return err
	}
	return nil
}

func GetRedisShortMessage(key string) (*ShortMessage, error) {
	sc := GetContainerInstance()
	value, err := sc.RedisService.GetKey(NewShortMessageKey(key))
	if err != nil {
		return nil, err
	}
	sc.RedisService.ExpireKey(key)
	var info *ShortMessage
	err = json.Unmarshal([]byte(value), &info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func StoreShortMessage(phone string, code string) error {
	sc := GetContainerInstance()
	tn := time.Now().Unix()
	codeInfo, err := GetRedisShortMessage(phone)
	if err != nil {
		codeInfo = &ShortMessage{
			Phone: phone,
			Time:  tn,
			Code:  code,
			Count: 1,
		}
	} else {
		codeInfo.Count++
		if codeInfo.Count > int64(sc.Config.ShortMessageLimitedCount) {
			return errors.New(message.ShortMessageLimitedError)
		} else if (tn - codeInfo.Time) < int64(sc.Config.ShortMessageSpaceTime) {
			return errors.New(message.ShortMessageFrequentlyError)
		}
		codeInfo.Code = code
		codeInfo.Time = tn
	}
	err = AddRedisShortMessage(codeInfo)
	if err != nil {
		return errors.New(message.HttpError)
	}
	return nil
}
