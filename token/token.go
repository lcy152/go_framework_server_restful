package token

import (
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
)

func Signed(priKey *rsa.PrivateKey, claims jwt.Claims) (string, error) {
	newToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims) //所有人给xiahualou发送公钥加密的数据，但是只有xiahualou本人可以使用私钥解密
	signedToken, err := newToken.SignedString(priKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func UnSigned(pubKey *rsa.PublicKey, signedToken string, uc jwt.Claims) error {
	token, err := jwt.ParseWithClaims(signedToken, uc, func(token *jwt.Token) (i interface{}, e error) { //使用私钥解密
		return pubKey, nil //这里的返回值必须是公钥，不然解密肯定是失败
	})
	if err != nil {
		return err
	}
	if err := token.Claims.Valid(); err != nil {
		return err
	}
	return nil
}
