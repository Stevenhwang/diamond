package utils

import (
	"diamond/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT 声明签名信息
type JWT struct {
	signKey []byte
}

// J 全局JWT
var J *JWT

func init() {
	secret := config.Config.Get("jwt.secret").(string)
	J = NewJWT(secret)
}

// NewJWT 初始化jwt对象
func NewJWT(jwtSecret string) *JWT {
	return &JWT{
		[]byte(jwtSecret),
	}
}

// EncodeToken token加密
func (j *JWT) EncodeToken(uid uint, username string, isSuperuser bool) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":          time.Now(),
		"uid":          uid,
		"username":     username,
		"is_superuser": isSuperuser,
	})
	tokenString, _ := token.SignedString(j.signKey)
	return tokenString
}

// DecodeToken token解密
func (j *JWT) DecodeToken(tokenString string) (uid uint, username string, isSuperuser bool) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.signKey, nil
	})
	if token != nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return uint(claims["uid"].(float64)), claims["username"].(string), claims["is_superuser"].(bool)
		}
	}
	return 0, "", false
}
