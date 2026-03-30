package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return setupRouter()
}

func TestCreateOrder(t *testing.T) {
	// Очищаем orders перед тестом
	orders = make(map[string]Order)
	
	router := setupTestRouter()

	order := Order{
		OrderID:       "ORD-001",
		CustomerName:  "John Doe",
		CustomerEmail: "john@example.com",
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
			ZipCode: "10001",
		},
		Products: []Product{
			{ID: "PROD-1", Name: "Laptop", Price: 999.99, Quantity: 1},
			{ID: "PROD-2", Name: "Mouse", Price: 29.99, Quantity: 2},
		},
	}

	jsonData, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/api/v1/order", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check if order was stored
	createdOrder, exists := orders[order.OrderID]
	assert.True(t, exists)
	assert.Equal(t, "created", createdOrder.Status)
	assert.Equal(t, 999.99+2*29.99, createdOrder.TotalAmount)
}

func TestCreateOrderInvalidData(t *testing.T) {
	// Очищаем orders перед тестом
	orders = make(map[string]Order)
	
	router := setupTestRouter()

	invalidOrder := map[string]interface{}{
		"order_id": "ORD-002",
		// Missing required fields
	}

	jsonData, _ := json.Marshal(invalidOrder)
	req, _ := http.NewRequest("POST", "/api/v1/order", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
}

func TestGetOrder(t *testing.T) {
	// Очищаем и заполняем orders перед тестом
	orders = make(map[string]Order)
	
	testOrder := Order{
		OrderID:       "ORD-003",
		CustomerName:  "Jane Doe",
		CustomerEmail: "jane@example.com",
		Address: Address{
			Street:  "456 Oak Ave",
			City:    "Los Angeles",
			Country: "USA",
			ZipCode: "90001",
		},
		Products: []Product{
			{ID: "PROD-3", Name: "Phone", Price: 699.99, Quantity: 1},
		},
	}
	orders[testOrder.OrderID] = testOrder

	router := setupTestRouter()
	req, _ := http.NewRequest("GET", "/api/v1/order/ORD-003", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestGetOrderNotFound(t *testing.T) {
	// Очищаем orders перед тестом
	orders = make(map[string]Order)
	
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/order/NONEXISTENT", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
}

func TestListOrders(t *testing.T) {
	// Очищаем и заполняем orders перед тестом
	orders = make(map[string]Order)
	orders["ORD-004"] = Order{OrderID: "ORD-004", CustomerName: "Test1"}
	orders["ORD-005"] = Order{OrderID: "ORD-005", CustomerName: "Test2"}

	router := setupTestRouter()
	req, _ := http.NewRequest("GET", "/api/v1/orders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 2)
}