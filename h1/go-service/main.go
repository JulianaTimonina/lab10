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
    "go-jwt-service/handlers"
    "go-jwt-service/middleware"
)

func main() {
    router := gin.Default()

    // Публичные маршруты
    router.POST("/login", handlers.Login)
    router.GET("/public", handlers.PublicData)

    // Защищенные маршруты
    protected := router.Group("/protected")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/data", handlers.ProtectedData)
    }

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    log.Println("Go service running on http://localhost:8080")

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