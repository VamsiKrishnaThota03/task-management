package middleware

import (
    "net/http"
    "task-management-system/internal/auth"
    "github.com/gin-gonic/gin"
)

// Middleware to allow only admin users
func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
            c.Abort()
            return
        }

        _, role, err := auth.VerifyToken(tokenString)
        if err != nil || role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
            c.Abort()
            return
        }

        c.Next()
    }
}