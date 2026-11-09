package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"task-management-system/internal/models"
	"task-management-system/internal/websocket"
	"task-management-system/pkg/ai"
)

const (
	maxTitleLength       = 200
	maxDescriptionLength = 5000
)

// TaskResponse is the JSON shape returned for tasks.
type TaskResponse struct {
	ID             uint              `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Status         models.TaskStatus `json:"status"`
	Priority       string            `json:"priority"`
	DueDate        *time.Time        `json:"dueDate"`
	AssignedTo     uint              `json:"assignedTo"`
	AssignedToName string            `json:"assignedToName"`
	AssignedBy     uint              `json:"assignedBy"`
	AssignedByName string            `json:"assignedByName"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

func usernames(db *gorm.DB, ids ...uint) map[uint]string {
	names := map[uint]string{}
	if len(ids) == 0 {
		return names
	}
	var users []models.User
	db.Select("id", "username").Where("id IN ?", ids).Find(&users)
	for _, u := range users {
		names[u.ID] = u.Username
	}
	return names
}

func toResponses(db *gorm.DB, tasks []models.Task) []TaskResponse {
	ids := make([]uint, 0, len(tasks)*2)
	for _, t := range tasks {
		ids = append(ids, t.AssignedTo, t.AssignedBy)
	}
	names := usernames(db, ids...)
	out := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		var due *time.Time
		if !t.DueDate.IsZero() {
			d := t.DueDate
			due = &d
		}
		out = append(out, TaskResponse{
			ID:             t.ID,
			Title:          t.Title,
			Description:    t.Description,
			Status:         t.Status,
			Priority:       t.Priority,
			DueDate:        due,
			AssignedTo:     t.AssignedTo,
			AssignedToName: names[t.AssignedTo],
			AssignedBy:     t.AssignedBy,
			AssignedByName: names[t.AssignedBy],
			CreatedAt:      t.CreatedAt,
			UpdatedAt:      t.UpdatedAt,
		})
	}
	return out
}

func isAdmin(c *gin.Context) bool { return c.GetString("role") == models.RoleAdmin }

// canView: admins, the assignee and the creator can see a task.
func canView(c *gin.Context, t *models.Task) bool {
	uid := c.GetUint("userID")
	return isAdmin(c) || t.AssignedTo == uid || t.AssignedBy == uid
}

// canEdit: admins and the creator can edit or reassign a task.
func canEdit(c *gin.Context, t *models.Task) bool {
	return isAdmin(c) || t.AssignedBy == c.GetUint("userID")
}

// loadTask fetches :id and writes 404 if it is missing or not visible to the caller.
func loadTask(c *gin.Context, db *gorm.DB) (*models.Task, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return nil, false
	}
	var task models.Task
	if err := db.First(&task, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		}
		return nil, false
	}
	if !canView(c, &task) {
		// Do not reveal that the task exists.
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return nil, false
	}
	return &task, true
}

func parseDueDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, nil
	}
	if t, err := time.Parse("2006-01-02", value); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, value)
}

func userExists(db *gorm.DB, id uint) bool {
	var count int64
	db.Model(&models.User{}).Where("id = ?", id).Count(&count)
	return count > 0
}

type taskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	AssignedTo  *uint   `json:"assignedTo"`
	DueDate     *string `json:"dueDate"`
	Priority    *string `json:"priority"`
	Status      *string `json:"status"`
}

// apply validates the input and copies it onto the task. It returns a
// user-facing error message, or "" when the input is valid.
func (in taskInput) apply(db *gorm.DB, task *models.Task) string {
	if in.Title != nil {
		title := strings.TrimSpace(*in.Title)
		if title == "" || len(title) > maxTitleLength {
			return "Title is required (max 200 characters)"
		}
		task.Title = title
	}
	if in.Description != nil {
		desc := strings.TrimSpace(*in.Description)
		if len(desc) > maxDescriptionLength {
			return "Description is too long (max 5000 characters)"
		}
		task.Description = desc
	}
	if in.AssignedTo != nil {
		if *in.AssignedTo == 0 || !userExists(db, *in.AssignedTo) {
			return "Assignee does not exist"
		}
		task.AssignedTo = *in.AssignedTo
	}
	if in.DueDate != nil {
		due, err := parseDueDate(*in.DueDate)
		if err != nil {
			return "Due date must be YYYY-MM-DD or RFC 3339"
		}
		task.DueDate = due
	}
	if in.Priority != nil {
		p := strings.ToUpper(strings.TrimSpace(*in.Priority))
		if !models.ValidPriority(p) {
			return "Priority must be LOW, MEDIUM or HIGH"
		}
		task.Priority = p
	}
	if in.Status != nil {
		s := models.TaskStatus(strings.ToUpper(strings.TrimSpace(*in.Status)))
		if !s.Valid() {
			return "Status must be PENDING, IN_PROGRESS or COMPLETED"
		}
		task.Status = s
	}
	return ""
}

// GetTasks lists the caller's tasks (assigned to or created by them); admins see all.
// Optional ?status= filter.
func GetTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := db.Model(&models.Task{})
		if !isAdmin(c) {
			uid := c.GetUint("userID")
			query = query.Where("assigned_to = ? OR assigned_by = ?", uid, uid)
		}
		if status := strings.ToUpper(c.Query("status")); status != "" {
			if !models.TaskStatus(status).Valid() {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status filter"})
				return
			}
			query = query.Where("status = ?", status)
		}

		var tasks []models.Task
		if err := query.Order("created_at DESC").Find(&tasks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
			return
		}
		c.JSON(http.StatusOK, toResponses(db, tasks))
	}
}

// GetTask returns a single task the caller can see.
func GetTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		task, ok := loadTask(c, db)
		if !ok {
			return
		}
		c.JSON(http.StatusOK, toResponses(db, []models.Task{*task})[0])
	}
}

// CreateTask creates a task assigned to an existing user.
func CreateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input taskInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		if input.Title == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
			return
		}
		if input.AssignedTo == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Assignee is required"})
			return
		}
		input.Status = nil // new tasks always start as PENDING

		task := models.Task{
			AssignedBy: c.GetUint("userID"),
			Status:     models.TaskStatusPending,
			Priority:   models.PriorityMedium,
		}
		if msg := input.apply(db, &task); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		if err := db.Create(&task).Error; err != nil {
			log.Printf("create task: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
			return
		}

		websocket.BroadcastTaskEvent("task_created", task.ID)
		c.JSON(http.StatusCreated, toResponses(db, []models.Task{task})[0])
	}
}

// UpdateTask edits a task (creator or admin only).
func UpdateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		task, ok := loadTask(c, db)
		if !ok {
			return
		}
		if !canEdit(c, task) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only the task creator or an admin can edit this task"})
			return
		}

		var input taskInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		if msg := input.apply(db, task); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		if err := db.Save(task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
			return
		}

		websocket.BroadcastTaskEvent("task_updated", task.ID)
		c.JSON(http.StatusOK, toResponses(db, []models.Task{*task})[0])
	}
}

// UpdateTaskStatus changes the status (assignee, creator or admin).
func UpdateTaskStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		task, ok := loadTask(c, db)
		if !ok {
			return
		}

		var input struct {
			Status string `json:"status"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		status := models.TaskStatus(strings.ToUpper(strings.TrimSpace(input.Status)))
		if !status.Valid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be PENDING, IN_PROGRESS or COMPLETED"})
			return
		}

		task.Status = status
		if err := db.Model(task).Update("status", status).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
			return
		}

		websocket.BroadcastTaskEvent("task_updated", task.ID)
		c.JSON(http.StatusOK, toResponses(db, []models.Task{*task})[0])
	}
}

// DeleteTask soft-deletes a task (admin only; enforced by the router).
func DeleteTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		task, ok := loadTask(c, db)
		if !ok {
			return
		}
		if err := db.Delete(task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
			return
		}
		websocket.BroadcastTaskEvent("task_deleted", task.ID)
		c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
	}
}

// GetTaskSuggestions asks Gemini to improve a description and propose subtasks.
func GetTaskSuggestions(client *ai.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if client == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI suggestions are not configured on this server"})
			return
		}

		var input struct {
			Description string `json:"description"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		description := strings.TrimSpace(input.Description)
		if description == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Description cannot be empty"})
			return
		}
		if len(description) > maxDescriptionLength {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Description is too long"})
			return
		}

		suggestions, err := client.TaskSuggestions(c.Request.Context(), description)
		if err != nil {
			log.Printf("gemini suggestions failed: %v", err)
			c.JSON(http.StatusBadGateway, gin.H{"error": "Could not get AI suggestions right now, please try again"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"suggestions":          suggestions,
			"original_description": description,
		})
	}
}
