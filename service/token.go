package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"tumor_server/message"
	"tumor_server/model"
)

const (
	TokenTAG = "token:"
)

type UserTokenInfo struct {
	User      *model.User `json:"user"`
	Token     string      `json:"token"`
	LoginTime int64       `json:"login_time"`
	IP        string      `json:"ip"`
}

type Token struct {
	ID        string `json:"_id"`
	LoginTime int64  `json:"login_time"`
}

func MarshalToken(t *Token) (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	encodeString := base64.StdEncoding.EncodeToString(data)
	return encodeString, nil
}

func UnmarshalToken(token string) (*Token, error) {
	t := new(Token)
	decodeBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(decodeBytes, t)
	if t.ID == "" {
		return nil, errors.New("empty guid")
	}
	return t, nil
}

func NewUserTokenInfoKey(key string) string {
	return TokenTAG + key
}

func AddUserTokenInfo(info *UserTokenInfo) error {
	sc := GetContainerInstance()
	message, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = sc.RedisService.AddKey(NewUserTokenInfoKey(info.User.ID.String()), string(message))
	if err != nil {
		return err
	}
	return nil
}

func GetUserTokenInfo(guid string) (*UserTokenInfo, error) {
	sc := GetContainerInstance()
	value, err := sc.RedisService.GetKey(NewUserTokenInfoKey(guid))
	if err != nil {
		return nil, err
	}
	sc.RedisService.ExpireKey(NewUserTokenInfoKey(guid))
	var info *UserTokenInfo
	err = json.Unmarshal([]byte(value), &info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func DeleteUserTokenInfo(guid string) error {
	sc := GetContainerInstance()
	err := sc.RedisService.DeleteKey(NewUserTokenInfoKey(guid))
	if err != nil {
		return err
	}
	return nil
}

func TokenValidate(token, ip string) (*UserTokenInfo, error) {
	sc := GetContainerInstance()
	tokenInfo, err := UnmarshalToken(token)
	if err != nil {
		return nil, errors.New(message.ValidateError)
	}
	userInfo, err := GetUserTokenInfo(tokenInfo.ID)
	if err != nil {
		user, err := sc.DB.GetUserByToken(context.TODO(), token)
		if err != nil {
			return nil, errors.New(message.ValidateError)
		}
		userInfo = &UserTokenInfo{
			User:      user,
			Token:     user.Token,
			LoginTime: tokenInfo.LoginTime,
			IP:        ip,
		}
		AddUserTokenInfo(userInfo)
		AddUserRecord(user.ID, model.UserOperationLogin, ip)
	} else if tokenInfo.LoginTime != userInfo.LoginTime {
		return nil, errors.New(message.OtherPlaceLoginError)
	}
	return userInfo, nil
}
