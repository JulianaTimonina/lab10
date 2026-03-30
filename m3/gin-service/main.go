package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Models
type Address struct {
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	Country string `json:"country" binding:"required"`
	ZipCode string `json:"zip_code" binding:"required"`
}

type Product struct {
	ID       string  `json:"id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Quantity int     `json:"quantity" binding:"required,gte=0"`
}

type Order struct {
	OrderID      string    `json:"order_id" binding:"required"`
	CustomerName string    `json:"customer_name" binding:"required"`
	CustomerEmail string   `json:"customer_email" binding:"required,email"`
	Address      Address   `json:"address" binding:"required"`
	Products     []Product `json:"products" binding:"required,min=1"`
	TotalAmount  float64   `json:"total_amount"`
	Status       string    `json:"status"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var orders = make(map[string]Order)

func main() {
	router := setupRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Gin service started on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	// Middleware for logging
	router.Use(LoggerMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{Success: true, Message: "OK"})
	})

	// Order endpoints
	api := router.Group("/api/v1")
	{
		api.POST("/order", createOrder)
		api.GET("/order/:id", getOrder)
		api.GET("/orders", listOrders)
	}

	return router
}

// Middleware for logging
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		log.Printf("[%s] %s %s - %d (%v)",
			startTime.Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(startTime),
		)
	}
}

func createOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Calculate total amount
	total := 0.0
	for _, product := range order.Products {
		total += product.Price * float64(product.Quantity)
	}
	order.TotalAmount = total
	order.Status = "created"

	orders[order.OrderID] = order

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Order created successfully",
		Data:    order,
	})
}

func getOrder(c *gin.Context) {
	id := c.Param("id")
	order, exists := orders[id]
	if !exists {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Order not found",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    order,
	})
}

func listOrders(c *gin.Context) {
	orderList := make([]Order, 0, len(orders))
	for _, order := range orders {
		orderList = append(orderList, order)
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    orderList,
	})
}