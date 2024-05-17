package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtService struct {
	PrivateKey string
}

type Claims struct {
	UserId int64 `json:"userId"`
	jwt.StandardClaims
}

func (s *JwtService) GenerateToken(userId int64) (string, error) {
	expirationTime := time.Now().Add(9 * time.Hour)
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.PrivateKey))
}

func (s *JwtService) VerifyToken(tokenString string) (int64, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.PrivateKey), nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	return claims.UserId, nil
}
