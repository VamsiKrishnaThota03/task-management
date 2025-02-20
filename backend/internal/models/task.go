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

type Task struct {
    gorm.Model         // This adds ID, CreatedAt, UpdatedAt, DeletedAt fields
    Title       string `gorm:"not null"`
    Description string `gorm:"type:text"`
    AssignedTo  uint   `gorm:"index"` // References User ID
    AssignedBy  uint   `gorm:"index"` // References User ID
    DueDate     time.Time
    Status      TaskStatus `gorm:"type:varchar(20);default:'PENDING'"`
    Priority    string    `gorm:"type:varchar(20);default:'MEDIUM'"`
    Tags        []TaskTag `gorm:"many2many:task_task_tags;"`
}

type TaskTag struct {
    gorm.Model         // This adds ID, CreatedAt, UpdatedAt, DeletedAt fields
    Name  string `gorm:"type:varchar(50);unique;not null"`
    Tasks []Task `gorm:"many2many:task_task_tags;"`
}

// Define the join table structure
type TaskTaskTag struct {
    TaskID    uint `gorm:"primaryKey"`
    TaskTagID uint `gorm:"primaryKey"`
}

// TableName specifies the table name for TaskTaskTag
func (TaskTaskTag) TableName() string {
    return "task_task_tags"
}

// TableName specifies the table name for TaskTag
func (TaskTag) TableName() string {
    return "task_tags"
}

// TableName specifies the table name for Task
func (Task) TableName() string {
    return "tasks"
}