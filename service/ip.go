package service

import "encoding/json"

const (
	IPTAG = "ip:"
)

type IPInfo struct {
	IP      string `json:"ip"`
	Time    int64  `json:"time"`
	Count   int64  `json:"count"`
	Disable bool   `json:"disable"`
}

func NewIPKey(key string) string {
	return IPTAG + key
}

func AddRedisRequestCount(ipInfo *IPInfo) error {
	sc := GetContainerInstance()
	message, err := json.Marshal(ipInfo)
	if err != nil {
		return err
	}
	err = sc.RedisService.AddKey(NewIPKey(ipInfo.IP), string(message))
	if err != nil {
		return err
	}
	return nil
}

func GetRedisRequestCount(key string) (*IPInfo, error) {
	sc := GetContainerInstance()
	value, err := sc.RedisService.GetKey(NewIPKey(key))
	if err != nil {
		return nil, err
	}
	sc.RedisService.ExpireKey(NewIPKey(key))
	var info *IPInfo
	err = json.Unmarshal([]byte(value), &info)
	if err != nil {
		return nil, err
	}
	return info, nil
}
