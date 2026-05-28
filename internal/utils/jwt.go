package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret string
}

type Claims struct {
	UserID string
	jwt.RegisteredClaims
}

func NewJwtHelper(secret string) *JWT {
	return &JWT{
		secret: secret,
	}
}

func (j *JWT) Generate(
	userID string,
) (string, error) {
	claims := Claims{

		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),
			NotBefore: jwt.NewNumericDate(
				time.Now(),
			),
		},
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	return token.SignedString(
		[]byte(j.secret),
	)
}

func (j *JWT) Verify(

	tokenString string,

) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			// Ensure signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(
					"invalid signing method",
				)
			}
			return []byte(j.secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil

}
