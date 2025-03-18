package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// RoleMiddleware проверяет роль пользователя
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Получаем роль из контекста
        role := c.GetString("role")
        if role == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
            return
        }

        // Проверяем, разрешена ли роль
        for _, allowedRole := range allowedRoles {
            if role == allowedRole {
                c.Next()
                return
            }
        }

        // Если роль не разрешена, возвращаем ошибку
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
    }
}