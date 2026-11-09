package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "PENDING"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
)

// Valid reports whether s is one of the known task statuses.
func (s TaskStatus) Valid() bool {
	switch s {
	case TaskStatusPending, TaskStatusInProgress, TaskStatusCompleted:
		return true
	}
	return false
}

const (
	PriorityLow    = "LOW"
	PriorityMedium = "MEDIUM"
	PriorityHigh   = "HIGH"
)

// ValidPriority reports whether p is one of the known priorities.
func ValidPriority(p string) bool {
	return p == PriorityLow || p == PriorityMedium || p == PriorityHigh
}

type Task struct {
	gorm.Model         // ID, CreatedAt, UpdatedAt, DeletedAt
	Title       string `gorm:"not null"`
	Description string `gorm:"type:text"`
	AssignedTo  uint   `gorm:"index"` // References User ID
	AssignedBy  uint   `gorm:"index"` // References User ID
	DueDate     time.Time
	Status      TaskStatus `gorm:"type:varchar(20);default:'PENDING'"`
	Priority    string     `gorm:"type:varchar(20);default:'MEDIUM'"`
	Tags        []TaskTag  `gorm:"many2many:task_task_tags;"`
}

type TaskTag struct {
	gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `gorm:"type:varchar(50);unique;not null"`
	Tasks      []Task `gorm:"many2many:task_task_tags;"`
}

// TaskTaskTag is the join table between tasks and tags.
type TaskTaskTag struct {
	TaskID    uint `gorm:"primaryKey"`
	TaskTagID uint `gorm:"primaryKey"`
}

func (TaskTaskTag) TableName() string { return "task_task_tags" }

func (TaskTag) TableName() string { return "task_tags" }

func (Task) TableName() string { return "tasks" }
