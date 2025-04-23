package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

/**
 * @File: jwt.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/23 下午2:59
 * @Software: GoLand
 * @Version:  1.0
 */

var JwtKey = []byte("A9v$34Fjfl1.pMv@z6k1!qwe93Km!")

func ParseToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return JwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
