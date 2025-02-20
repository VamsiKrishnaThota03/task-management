package handlers

import (
    "net/http"
    "task-management-system/internal/models"
    "task-management-system/internal/auth"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
)

// Register a new user
func RegisterUser(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var user models.User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Default role to "user" if not provided
        if user.Role == "" {
            user.Role = models.RoleUser
        }

        if err := user.HashPassword(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
            return
        }

        if err := db.Create(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "role": user.Role})
    }
}

func LoginUser(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input models.User
        var user models.User

        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
            return
        }

        if !user.CheckPassword(input.Password) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
            return
        }

        token, err := auth.GenerateToken(user.ID, user.Role) // ✅ Include role in token
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"token": token, "role": user.Role}) // ✅ Return role to frontend
    }
}

// GetUsers returns all users (admin only)
func GetUsers(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var users []models.User
        
        // Fetch all users but exclude password field
        if err := db.Select("id, username, role, created_at, updated_at").Find(&users).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
            return
        }

        c.JSON(http.StatusOK, users)
    }
}

// GetProfile returns the user's profile
func GetProfile(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get userID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
            return
        }

        var user models.User
        if err := db.First(&user, userID).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
                return
            }
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile"})
            return
        }

        c.JSON(http.StatusOK, user.ToProfileResponse())
    }
}
