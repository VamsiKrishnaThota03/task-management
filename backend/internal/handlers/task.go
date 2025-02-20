package handlers

import (
    "net/http"
    "task-management-system/internal/models"
    "task-management-system/internal/config"
    "task-management-system/internal/websocket"
    "task-management-system/pkg/ai"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "fmt"
    "log"
)

// Get all tasks
func GetTasks(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var tasks []models.Task
        userID, _ := c.Get("userID")

        // If user is admin, show all tasks, otherwise show only assigned tasks
        query := db
        if role, _ := c.Get("role"); role != "admin" {
            query = query.Where("assigned_to = ?", userID)
        }

        if err := query.Find(&tasks).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
            return
        }

        c.JSON(http.StatusOK, tasks)
    }
}

// Update task status
func UpdateTaskStatus(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        taskID := c.Param("id")
        var task models.Task
        
        if err := db.First(&task, taskID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
            return
        }

        var input struct {
            Status string `json:"status"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        task.Status = models.TaskStatus(input.Status)
        if err := db.Save(&task).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
            return
        }

        // Broadcast task update
        websocket.BroadcastMessage("Task updated: " + task.Title)
        c.JSON(http.StatusOK, task)
    }
}

// Get AI suggestions for task
func GetTaskSuggestions(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("Received suggestion request")

        var input struct {
            Description string `json:"description" binding:"required"`
        }
        
        if err := c.ShouldBindJSON(&input); err != nil {
            log.Printf("Invalid input: %v", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %v", err)})
            return
        }

        if input.Description == "" {
            log.Println("Empty description received")
            c.JSON(http.StatusBadRequest, gin.H{"error": "Description cannot be empty"})
            return
        }

        log.Printf("Processing suggestion for description: %s", input.Description)

        if cfg.GeminiKey == "" {
            log.Println("Gemini API key not configured")
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Gemini API key not configured"})
            return
        }

        prompt := fmt.Sprintf("Please suggest improvements and subtasks for this task description: %s", input.Description)
        
        log.Println("Calling Gemini API")
        suggestions, err := ai.GetTaskSuggestions(prompt, cfg.GeminiKey)
        if err != nil {
            log.Printf("Error getting suggestions: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": fmt.Sprintf("Failed to get AI suggestions: %v", err),
                "details": err.Error(),
            })
            return
        }

        if suggestions == "" {
            log.Println("Received empty suggestions")
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Received empty response from Gemini"})
            return
        }

        log.Printf("Successfully received suggestions from Gemini")
        c.JSON(http.StatusOK, gin.H{
            "suggestions": suggestions,
            "original_description": input.Description,
        })
    }
}

func CreateTask(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get the current user ID from the context
        userID, exists := c.Get("userID")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
            return
        }

        var input struct {
            Title       string    `json:"title" binding:"required"`
            Description string    `json:"description" binding:"required"`
            AssignedTo  uint      `json:"assignedTo" binding:"required"`
            DueDate     string    `json:"dueDate" binding:"required"`
            Priority    string    `json:"priority"`
        }

        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        task := models.Task{
            Title:       input.Title,
            Description: input.Description,
            AssignedTo:  input.AssignedTo,
            AssignedBy:  userID.(uint),
            Status:      models.TaskStatusPending,
            Priority:    input.Priority,
        }

        if err := db.Create(&task).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task: " + err.Error()})
            return
        }

        c.JSON(http.StatusOK, task)
    }
}

func DeleteTask(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        taskID := c.Param("id")

        if err := db.Delete(&models.Task{}, taskID).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
    }
}
