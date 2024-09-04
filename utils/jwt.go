package utils

import (
	"github.com/dotdancer/gogofly/global"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWT struct {
	SigningKey []byte
}

type JwtCustomClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(global.Config.Jwt.SigningKey),
	}
}

func (j *JWT) GenerateToken(id int, name string) (token string, err error) {

	iJwtCustomClaims := JwtCustomClaims{
		ID:       id,
		Username: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.Config.Jwt.TokenExpire) * time.Minute)), //过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                                 //签发时间
			Subject:   "token",                                                                                        //主题
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, iJwtCustomClaims)
	return tokenClaims.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return tokenClaims, nil
	} else {
		return nil, err
	}
}

func (j *JWT) IsTokenValid(tokenString string) bool {
	_, err := j.ParseToken(tokenString)
	if err != nil {
		return false
	}
	return true
}
