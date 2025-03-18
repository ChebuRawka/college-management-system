package services

import (
    "backend/models"
    "backend/repository"
    "errors"
    "time"
    "github.com/dgrijalva/jwt-go"
)

type AuthService struct {
    Repo *repositories.UserRepository
    SecretKey string
}

func NewAuthService(repo *repositories.UserRepository, secretKey string) *AuthService {
    return &AuthService{Repo: repo, SecretKey: secretKey}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(username, password, role string) error {
    // Проверяем роль
    if role != "admin" && role != "teacher" {
        return errors.New("invalid role")
    }

    // Хэшируем пароль
    user := &models.User{
        Username: username,
        Role:     role,
    }
    if err := user.HashPassword(password); err != nil {
        return err
    }

    // Создаем пользователя
    return s.Repo.CreateUser(user)
}

// Login авторизует пользователя и возвращает JWT-токен
func (s *AuthService) Login(username, password string) (string, error) {
    user, err := s.Repo.GetUserByUsername(username)
    if err != nil {
        return "", err
    }
    if user == nil || !user.CheckPassword(password) {
        return "", errors.New("invalid credentials")
    }

    // Генерируем JWT-токен
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
    })

    tokenString, err := token.SignedString([]byte(s.SecretKey))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}