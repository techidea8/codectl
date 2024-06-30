package token

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GenToken 生成token
func GenerateTokenWithSecret(values map[string]interface{}, secretKey string) (string, error) {
	claim := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"iss": "turing-microapp-server",
		"aud": "turing-microapp-server",
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
		"sub": "user",
	}
	for key, value := range values {
		claim[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// VerfyToken 验证token
func ParseTokenWithSecret(token, secretKey string) (map[string]interface{}, error) {
	tokenObj, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tokenObj.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot convert claim to mapclaim")
	}

	if !tokenObj.Valid {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}
