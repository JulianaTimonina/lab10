// m1/main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"taskapi/handlers"
	"taskapi/middleware"
	"taskapi/models"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())

	api := r.Group("/api")
	{
		tasks := api.Group("/tasks")
		{
			tasks.POST("", handlers.CreateTask)
			tasks.GET("", handlers.GetTasks)
			tasks.GET("/:id", handlers.GetTask)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.DELETE("/:id", handlers.DeleteTask)
		}
	}
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}

func TestHealthCheck(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}

func TestCreateTask(t *testing.T) {
	router := setupRouter()

	task := models.CreateTaskRequest{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
	}

	jsonData, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response models.SuccessResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if response.Message != "Task created successfully" {
		t.Errorf("Expected success message, got '%s'", response.Message)
	}
}

func TestCreateTaskValidationError(t *testing.T) {
	router := setupRouter()

	// Task with invalid status
	task := models.CreateTaskRequest{
		Title:  "Test",
		Status: "invalid",
	}

	jsonData, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetTasks(t *testing.T) {
	router := setupRouter()

	// Create a task first
	task := models.CreateTaskRequest{
		Title:  "Task for list",
		Status: "pending",
	}
	jsonData, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Get all tasks
	req, _ = http.NewRequest("GET", "/api/tasks", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/tasks/non-existent-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateTask(t *testing.T) {
	router := setupRouter()

	// Create a task
	createTask := models.CreateTaskRequest{
		Title:  "Original Title",
		Status: "pending",
	}
	jsonData, _ := json.Marshal(createTask)
	req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var createResponse models.SuccessResponse
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	taskData := createResponse.Data.(map[string]interface{})
	taskID := taskData["id"].(string)

	// Update the task
	updateTask := models.UpdateTaskRequest{
		Title:  "Updated Title",
		Status: "in-progress",
	}
	jsonData, _ = json.Marshal(updateTask)
	req, _ = http.NewRequest("PUT", "/api/tasks/"+taskID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteTask(t *testing.T) {
	router := setupRouter()

	// Create a task
	createTask := models.CreateTaskRequest{
		Title:  "Task to delete",
		Status: "pending",
	}
	jsonData, _ := json.Marshal(createTask)
	req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var createResponse models.SuccessResponse
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	taskData := createResponse.Data.(map[string]interface{})
	taskID := taskData["id"].(string)

	// Delete the task
	req, _ = http.NewRequest("DELETE", "/api/tasks/"+taskID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify task is deleted
	req, _ = http.NewRequest("GET", "/api/tasks/"+taskID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 after deletion, got %d", w.Code)
	}
}