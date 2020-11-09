package impl

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type IhuyiResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Smid string `json:"smsid"`
}

func SendMessageToIhuyi(phone, text string) *IhuyiResponse {
	GetMd5String := func(s string) string {
		h := md5.New()
		h.Write([]byte(s))
		return hex.EncodeToString(h.Sum(nil))
	}

	v := url.Values{}
	_now := strconv.FormatInt(time.Now().Unix(), 10)
	_account := "cf_dutsf"
	_password := "ef9db60da48633bbf421e86b9a879c18"
	_mobile := phone
	_content := text
	v.Set("account", _account)
	v.Set("password", GetMd5String(_account+_password+_mobile+_content+_now))
	v.Set("mobile", _mobile)
	v.Set("content", _content)
	v.Set("time", _now)
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://106.ihuyi.com/webservice/sms.php?method=Submit&format=json", body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("error: SendMessageToIhuyi: ", err)
		return nil
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	ihuyi := new(IhuyiResponse)
	err = json.Unmarshal(data, ihuyi)
	if err != nil {
		log.Println("error: SendMessageToIhuyi: ", err)
		return nil
	}
	return ihuyi
}
