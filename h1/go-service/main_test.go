package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
    "go-jwt-service/handlers"
    "go-jwt-service/middleware"
)

func TestLoginEndpoint(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.POST("/login", handlers.Login)

    tests := []struct {
        name       string
        username   string
        password   string
        wantStatus int
    }{
        {"Valid credentials", "admin", "password123", http.StatusOK},
        {"Invalid username", "wrong", "password123", http.StatusUnauthorized},
        {"Invalid password", "admin", "wrong", http.StatusUnauthorized},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            reqBody := handlers.LoginRequest{
                Username: tt.username,
                Password: tt.password,
            }
            jsonData, _ := json.Marshal(reqBody)

            req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            if w.Code != tt.wantStatus {
                t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
            }

            if tt.wantStatus == http.StatusOK {
                var resp handlers.LoginResponse
                if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
                    t.Errorf("Failed to parse response: %v", err)
                }
                if resp.Token == "" {
                    t.Error("Expected token, got empty")
                }
            }
        })
    }
}

func TestProtectedEndpoint(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    
    protected := router.Group("/protected")
    protected.Use(middleware.AuthMiddleware())
    protected.GET("/data", handlers.ProtectedData)

    // Сначала получим токен
    router.POST("/login", handlers.Login)
    
    reqBody := handlers.LoginRequest{
        Username: "admin",
        Password: "password123",
    }
    jsonData, _ := json.Marshal(reqBody)
    
    loginReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
    loginReq.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, loginReq)
    
    var loginResp handlers.LoginResponse
    json.Unmarshal(w.Body.Bytes(), &loginResp)

    tests := []struct {
        name       string
        token      string
        wantStatus int
    }{
        {"Valid token", loginResp.Token, http.StatusOK},
        {"No token", "", http.StatusUnauthorized},
        {"Invalid token", "invalid.token.here", http.StatusUnauthorized},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, _ := http.NewRequest("GET", "/protected/data", nil)
            if tt.token != "" {
                req.Header.Set("Authorization", "Bearer "+tt.token)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            if w.Code != tt.wantStatus {
                t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
            }
        })
    }
}