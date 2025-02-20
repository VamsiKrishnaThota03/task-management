package auth

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = []byte("your_secret_key") // Store in .env for security

// Generate JWT token for a user
func GenerateToken(userID uint, role string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "role":    role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
    })

    return token.SignedString(JWT_SECRET)
}

// Verify JWT token and extract user ID & Role
func VerifyToken(tokenString string) (uint, string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return JWT_SECRET, nil
    })

    if err != nil || !token.Valid {
        return 0, "", errors.New("invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return 0, "", errors.New("invalid claims")
    }

    userID := uint(claims["user_id"].(float64))
    role := claims["role"].(string)
    return userID, role, nil
}
