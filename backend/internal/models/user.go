package models

import (
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
)

// Define user roles
const (
    RoleUser  = "user"
    RoleAdmin = "admin"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique" json:"username"`
    Password string `json:"-"` // "-" means this field won't be included in JSON
    Role     string `gorm:"default:user" json:"role"`
}

// Hash user password before saving
func (u *User) HashPassword() error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    return nil
}

// Check if entered password matches stored hash
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}

// Add a method to convert to profile response
func (u *User) ToProfileResponse() gin.H {
    return gin.H{
        "id":        u.ID,
        "username":  u.Username,
        "role":     u.Role,
        "createdAt": u.CreatedAt,
    }
}
