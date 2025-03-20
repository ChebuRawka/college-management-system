package repositories

import (
    "backend/models"
    "database/sql"
    "errors"
    "strings"
    "fmt"
)

type UserRepository struct {
    DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{DB: db}
}

// CreateUser создает нового пользователя
func (r *UserRepository) CreateUser(user *models.User) error {
    query := `
        INSERT INTO users (username, password_hash, role)
        VALUES ($1, $2, $3)
        RETURNING id
    `
    return r.DB.QueryRow(query, user.Username, user.PasswordHash, user.Role).Scan(&user.ID)
}

// GetUserByUsername находит пользователя по имени
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
    query := `
        SELECT id, username, password_hash, role
        FROM users
        WHERE username = $1
    `
    row := r.DB.QueryRow(query, username)

    var user models.User
    err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) UpdateTeacherProfile(userID int, updates map[string]interface{}) error {
    if len(updates) == 0 {
        fmt.Println("Error: no fields to update")
        return fmt.Errorf("no fields to update")
    }

    allowedFields := map[string]bool{
        "username": true,
        "password": true,
    }

    query := "UPDATE users SET "
    var args []interface{}
    var setClauses []string

    for key, value := range updates {
        if !allowedFields[key] {
            fmt.Printf("Error: invalid field %s\n", key)
            return fmt.Errorf("invalid field: %s", key)
        }
        setClauses = append(setClauses, fmt.Sprintf("%s = ?", key))
        args = append(args, value)
    }

    // Логируем поля для обновления
    fmt.Println("Fields to update:", setClauses)

    // Добавляем все поля через запятую
    query += strings.Join(setClauses, ", ")

    // Добавляем условие WHERE
    query += " WHERE id = ?"
    args = append(args, userID)

    // Логируем сформированный SQL-запрос и аргументы
    fmt.Println("Generated query:", query)
    fmt.Println("Query arguments:", args)

    _, err := r.DB.Exec(query, args...)
    if err != nil {
        fmt.Println("Database error:", err) // Логируем ошибку базы данных
        return err
    }

    return nil
}