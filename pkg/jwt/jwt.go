package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrTokenInvalid       = errors.New("token 签名无效或已被篡改")
	ErrTokenExpired       = errors.New("token 已过期")
	ErrTokenClaimsInvalid = errors.New("token 格式错误，载荷无法解析")
	// ErrTokenMalformed     = errors.New("token 格式非法")
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
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func (j *JWTHandler) GenerateToken(userID string, username string) (string, error) {
	issuedAt := time.Now()
	expirationTime := issuedAt.Add(3 * time.Hour)
	claims := Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issuedAt.Unix(),
			ExpiresAt: expirationTime.Unix(),
			Issuer:    j.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(j.secret)
	return tokenString, err
}

func (j *JWTHandler) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) { return j.secret, nil })
	if token == nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrTokenInvalid
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrTokenClaimsInvalid
	}
	if claims.ExpiresAt > 0 && claims.ExpiresAt < time.Now().Unix() {
		return nil, ErrTokenExpired
	}
	return claims, nil
}
