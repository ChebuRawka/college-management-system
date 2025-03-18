package middleware

import (

    "net/http"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
	"fmt"
	"strings"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Получаем заголовок Authorization
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            fmt.Println("Error: Missing token in Authorization header")
            return
        }

        // Удаляем префикс "Bearer " из токена
        if len(tokenString) > 7 && strings.ToUpper(tokenString[:7]) == "BEARER " {
            tokenString = tokenString[7:]
        } else {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
            fmt.Println("Error: Token format is invalid")
            return
        }

        fmt.Println("Received token:", tokenString) // Отладочное сообщение

        // Парсим токен
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secretKey), nil
        })
        if err != nil {
            fmt.Println("Error parsing token:", err) // Отладочное сообщение
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "details": err.Error()})
            return
        }

        if !token.Valid {
            fmt.Println("Token is not valid") // Отладочное сообщение
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        // Проверяем claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            fmt.Println("Invalid token claims") // Отладочное сообщение
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
            return
        }

        fmt.Println("Token is valid. Claims:", claims) // Отладочное сообщение

        // Устанавливаем user_id и role в контексте запроса
        c.Set("user_id", int(claims["user_id"].(float64)))
        c.Set("role", claims["role"].(string))
        c.Next()
    }
}