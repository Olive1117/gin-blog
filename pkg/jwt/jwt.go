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
	Type  string   `json:"type"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func (j *JWTHandler) GenerateAccessToken(userID string, roles []string) (string, time.Time, error) {
	issuedAt := time.Now()
	expirationTime := issuedAt.Add(15 * time.Minute)
	claims := Claims{
		Type:  "Access",
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
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
func (j *JWTHandler) GenerateRefreshToken(userID string, roles []string, tokenID string) (string, time.Time, error) {
	issuedAt := time.Now()
	expirationTime := issuedAt.Add(7 * 24 * time.Hour)
	claims := Claims{
		Type:  "Refresh",
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    j.issuer,
			ID:        tokenID,
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
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
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
