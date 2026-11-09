package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"task-management-system/internal/auth"
	"task-management-system/internal/models"
)

var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_.-]{3,32}$`)

const minPasswordLength = 8

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterUser creates a regular user. The role is always "user": admins are
// promoted out of band with `go run ./cmd/promote -username <name>`.
func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input credentials
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		input.Username = strings.TrimSpace(input.Username)

		if !usernamePattern.MatchString(input.Username) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be 3-32 characters: letters, numbers, '.', '_' or '-'"})
			return
		}
		if len(input.Password) < minPasswordLength || len(input.Password) > 72 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be 8-72 characters"})
			return
		}

		var count int64
		db.Model(&models.User{}).Where("LOWER(username) = LOWER(?)", input.Username).Count(&count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
			return
		}

		user := models.User{Username: input.Username, Role: models.RoleUser}
		if err := user.SetPassword(input.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		if err := db.Create(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(strings.ToLower(err.Error()), "unique") {
				c.JSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user.ToResponse()})
	}
}

// LoginUser checks the credentials and returns a JWT.
func LoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input credentials
		if err := c.ShouldBindJSON(&input); err != nil || input.Username == "" || input.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
			return
		}

		var user models.User
		err := db.Where("username = ?", strings.TrimSpace(input.Username)).First(&user).Error
		if err != nil || !user.CheckPassword(input.Password) {
			// Same message for unknown user and wrong password.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token, "role": user.Role, "user": user.ToResponse()})
	}
}

// GetUsers returns all users with their roles (admin only).
func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := db.Select("id", "username", "role", "created_at").Order("username").Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		out := make([]models.UserResponse, 0, len(users))
		for i := range users {
			out = append(out, users[i].ToResponse())
		}
		c.JSON(http.StatusOK, out)
	}
}

// GetAssignees returns the id and username of every user so any signed-in
// user can pick an assignee when creating a task.
func GetAssignees(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type assignee struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
		}
		var out []assignee
		if err := db.Model(&models.User{}).Select("id", "username").Order("username").Scan(&out).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		if out == nil {
			out = []assignee{}
		}
		c.JSON(http.StatusOK, out)
	}
}

// GetProfile returns the current user's profile.
func GetProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.First(&user, c.GetUint("userID")).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile"})
			return
		}
		c.JSON(http.StatusOK, user.ToResponse())
	}
}
