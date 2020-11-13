package token

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserClaim struct { //这个结构体主要是用来宣示当前公钥的使用者是谁，只有使用者和公钥的签名者是同一个人才可以用来正确的解密，还可以设置其他的属性，可以去百度一下
	UserName           string `json:"username"`
	jwt.StandardClaims        //嵌套了这个结构体就实现了Claim接口
}

func Test_GenerateToken(t *testing.T) {
	err := GenRSAPubAndPri(1024, "../pem") //1024是长度，长度越长安全性越高，但是性能也就越差
	if err != nil {
		log.Fatal(err)
	}
}

func Test_Token1(t *testing.T) {
	priBytes, err := ioutil.ReadFile("../pem/private.pem")
	if err != nil {
		log.Fatal("私钥文件读取失败")
	}
	pubBytes, err := ioutil.ReadFile("../pem/public.pem")
	if err != nil {
		log.Fatal("公钥文件读取失败")
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		log.Fatal("公钥文件不正确")
	}
	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(priBytes)
	if err != nil {
		log.Fatal("私钥文件不正确")
	}

	uc2 := &UserClaim{
		UserName: "xiahualou",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 2,
		},
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodRS256, uc2) //所有人给xiahualou发送公钥加密的数据，但是只有xiahualou本人可以使用私钥解密
	signedToken, err := newToken.SignedString(priKey)
	if err != nil {
		println(err)
		return
	}

	uc := &UserClaim{}
	getToken, err := jwt.ParseWithClaims(signedToken, uc, func(token *jwt.Token) (i interface{}, e error) { //使用私钥解密
		return pubKey, nil //这里的返回值必须是公钥，不然解密肯定是失败
	})
	if err != nil {
		println(err)
		return
	}

	if getToken.Valid { //服务端验证token是否有效
		fmt.Println(getToken.Claims.(*UserClaim).UserName)
	}
	if err := getToken.Claims.Valid(); err != nil {
		print(err)
		return
	}
	fmt.Println("ok")
}

func Test_Token2(t *testing.T) {
	priBytes, err := ioutil.ReadFile("../pem/private.pem")
	if err != nil {
		log.Fatal("私钥文件读取失败")
	}
	pubBytes, err := ioutil.ReadFile("../pem/public.pem")
	if err != nil {
		log.Fatal("公钥文件读取失败")
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		log.Fatal("公钥文件不正确")
	}
	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(priBytes)
	if err != nil {
		log.Fatal("私钥文件不正确")
	}

	uc := &UserClaim{
		UserName: "xiahualou",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 5,
		},
	}
	signedToken, err := Signed(priKey, uc)
	if err != nil {
		println(err)
		return
	}

	uc2 := &UserClaim{}
	err = UnSigned(pubKey, signedToken, uc2)
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println(uc2)
}
