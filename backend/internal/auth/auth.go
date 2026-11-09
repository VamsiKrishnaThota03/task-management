package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenTTL is how long issued tokens stay valid.
const TokenTTL = 24 * time.Hour

var secret []byte

// Init sets the HMAC key used to sign and verify tokens. It must be called
// once at startup with the JWT_SECRET from the environment.
func Init(jwtSecret string) {
	secret = []byte(jwtSecret)
}

// GenerateToken issues an HS256 JWT for the user.
func GenerateToken(userID uint, role string) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("auth: signing secret not configured")
	}
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"iat":     now.Unix(),
		"exp":     now.Add(TokenTTL).Unix(),
	})
	return token.SignedString(secret)
}

// VerifyToken validates the token signature, algorithm and expiry and returns
// the user ID and role it carries.
func VerifyToken(tokenString string) (uint, string, error) {
	if len(secret) == 0 {
		return 0, "", errors.New("auth: signing secret not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil || !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid claims")
	}
	rawID, ok := claims["user_id"].(float64)
	if !ok || rawID <= 0 {
		return 0, "", errors.New("invalid claims")
	}
	role, _ := claims["role"].(string)
	return uint(rawID), role, nil
}
