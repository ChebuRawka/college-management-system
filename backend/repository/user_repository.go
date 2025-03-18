package repositories

import (
    "backend/models"
    "database/sql"
    "errors"
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