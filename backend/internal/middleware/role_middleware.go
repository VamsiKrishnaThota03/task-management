package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"task-management-system/internal/models"
)

// AdminMiddleware allows only admins. It must run after JWTAuthMiddleware,
// which puts the user's current role in the context.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != models.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}
		c.Next()
	}
}
