package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired = errors.New("token 已过期")
	ErrTokenInvalid = errors.New("token 无效")
)

type JWTHandler struct {
	secret []byte
	issuer string
}

func NewJWT(secret string, issuer string) *JWTHandler {
	return &JWTHandler{
		secret: []byte(secret),
		issuer: issuer,
	}
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (j *JWTHandler) GenerateToken(userID uint, username string) (string, time.Time, error) {
	issuedAt := time.Now()
	expirationTime := issuedAt.Add(3 * time.Hour)
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    j.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expirationTime, nil
}

func (j *JWTHandler) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
