package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "go-jwt-service/auth"
)

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string `json:"token"`
}

func Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Простая проверка - в реальном приложении используйте БД
    if req.Username != "admin" || req.Password != "password123" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := auth.GenerateToken(req.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, LoginResponse{Token: token})
}

func ProtectedData(c *gin.Context) {
    claims, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":  "This is protected data",
        "username": claims.(*auth.Claims).Username,
    })
}

func PublicData(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "This is public data",
    })
}