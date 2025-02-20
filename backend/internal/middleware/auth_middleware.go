package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "task-management-system/internal/auth"
)

func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
            c.Abort()
            return
        }

        // ✅ Ensure "Bearer " prefix is removed before verification
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        userID, role, err := auth.VerifyToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("userID", userID)
        c.Set("role", role)
        c.Next()
    }
}
