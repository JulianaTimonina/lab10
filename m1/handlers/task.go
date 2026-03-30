// m1/handlers/task.go
package handlers

import (
	"net/http"
	"sync"
	"time"

	"taskapi/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TaskStore is an in-memory store for tasks
type TaskStore struct {
	mu    sync.RWMutex
	tasks map[string]models.Task
}

// NewTaskStore creates a new task store
func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]models.Task),
	}
}

var store = NewTaskStore()

// CreateTask handles POST /api/tasks
func CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	now := time.Now()
	task := models.Task{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	store.mu.Lock()
	store.tasks[task.ID] = task
	store.mu.Unlock()

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Task created successfully",
		Data:    task,
	})
}

// GetTasks handles GET /api/tasks
func GetTasks(c *gin.Context) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	tasks := make([]models.Task, 0, len(store.tasks))
	for _, task := range store.tasks {
		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Tasks retrieved successfully",
		Data:    tasks,
	})
}

// GetTask handles GET /api/tasks/:id
func GetTask(c *gin.Context) {
	id := c.Param("id")

	store.mu.RLock()
	task, exists := store.tasks[id]
	store.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Task not found"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Task retrieved successfully",
		Data:    task,
	})
}

// UpdateTask handles PUT /api/tasks/:id
func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	task, exists := store.tasks[id]
	if !exists {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Task not found"})
		return
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	task.UpdatedAt = time.Now()

	store.tasks[id] = task

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Task updated successfully",
		Data:    task,
	})
}

// DeleteTask handles DELETE /api/tasks/:id
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.tasks[id]; !exists {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Task not found"})
		return
	}

	delete(store.tasks, id)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Task deleted successfully",
	})
}