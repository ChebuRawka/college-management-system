package handlers

import (
    "backend/services"
    "net/http"

    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    Service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
    return &AuthHandler{Service: service}
}

// Register регистрирует нового пользователя
func (h *AuthHandler) Register(c *gin.Context) {
    var input struct {
        Username string `json:"username"`
        Password string `json:"password"`
        Role     string `json:"role"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    if err := h.Service.Register(input.Username, input.Password, input.Role); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// Login авторизует пользователя
func (h *AuthHandler) Login(c *gin.Context) {
    var input struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    token, err := h.Service.Login(input.Username, input.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}