/**
 * @Author: Resynz
 * @Date: 2021/9/24 16:47
 */
package common

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	goaes "github.com/rosbit/go-aes"
	"resync/lib/crypto"
	"time"
)

const (
	Secret = "immortals-resync-resynz"
	AesKey = "keep your pride."
	TokenData = "immortals-resync"
)

type Auth struct {
	Id     int64            `json:"id"`
	Name   string           `json:"name"`
	Ip     string           `json:"ip"`
	Expire int64            `json:"expire"`
}

func (s *Auth) string() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func (s *Auth) GenToken() (string, error) {
	s.Expire = time.Now().Unix() + 7 * 86400
	c, err := goaes.AesEncrypt([]byte(s.string()), []byte(AesKey))
	if err != nil {
		return "", err
	}
	m := jwt.MapClaims{TokenData: crypto.Base64EncodeToString(c)}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, m)
	return t.SignedString([]byte(Secret))
}

func (s *Auth) ParseToken(token string) error {
	t, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("[ParseToken] Unexpect signing method:%v", token.Header["alg"])
		}
		return []byte(Secret), nil
	})
	if err != nil || t == nil {
		return fmt.Errorf("[ParseToken] failed")
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return fmt.Errorf("[ParseToken] invalid token")
	}
	d, err := crypto.Base64DecodeString(claims[TokenData].(string))
	if err != nil {
		return fmt.Errorf("[ParseToken] decode failed")
	}
	decrypt, err := goaes.AesDecrypt(d, []byte(AesKey))
	if err != nil {
		return fmt.Errorf("[ParseToken] aes decrypt failed")
	}
	if err = json.Unmarshal(decrypt, s); err != nil {
		return fmt.Errorf("[ParseToken] json unmarshal failed")
	}
	return nil
}
