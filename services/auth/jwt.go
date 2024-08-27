package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lunatictiol/that-pet-place-backend-go/config"
)

var secretKey = config.Envs.JWTSecret

func GenerateToken(userId int64) (string, error) {
	exp := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{

		"userId": userId,
		"exp":    exp,
	})
	return token.SignedString([]byte(secretKey))
}
func VerifyToken(token string) (int64, error) {
	parse, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, errors.New("could not parse token")
	}
	if !parse.Valid {
		return 0, errors.New("not a valid token")
	}
	claims, ok := parse.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}
	userId := int64(claims["userId"].(float64))

	return userId, nil

}
