package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/michaelodeh/go-kit/dto"
	"github.com/michaelodeh/go-kit/utils"
)

var JWTSecretKey = utils.GetEnv("JWT_SECRET", "alapa-ai-auth")
var JWTRefreshSecretKey = utils.GetEnv("JWT_REFRESH_SECRET", "alapa-ai-auth-refresh")

type JWTAuth struct {
	secretKey []byte
}

func NewJWTAuth(secretKey string) *JWTAuth {
	return &JWTAuth{
		secretKey: []byte(secretKey),
	}
}

func (j *JWTAuth) Create(userID string, duration time.Duration) (string, error) {
	claims := dto.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    utils.GetEnv("APP_NAME", "alapa-ai-auth"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (j *JWTAuth) Verify(signedToken string) (*dto.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &dto.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*dto.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (j *JWTAuth) Decode(signedToken string) (*dto.JWTClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(signedToken, &dto.JWTClaims{})
	if err != nil {
		return nil, err
	}

	return token.Claims.(*dto.JWTClaims), nil
}
