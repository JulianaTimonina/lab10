// m1/models/task.go
package models

import "time"

// Task represents a task entity
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" binding:"required,min=3,max=100"`
	Description string    `json:"description" binding:"max=500"`
	Status      string    `json:"status" binding:"required,oneof=pending in-progress done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTaskRequest represents request body for creating a task
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"max=500"`
	Status      string `json:"status" binding:"required,oneof=pending in-progress done"`
}

// UpdateTaskRequest represents request body for updating a task
type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"omitempty,min=3,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Status      string `json:"status" binding:"omitempty,oneof=pending in-progress done"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}