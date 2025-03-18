package models

import "golang.org/x/crypto/bcrypt"

type User struct {
    ID           int    `json:"id"`
    Username     string `json:"username" validate:"required"`
    PasswordHash string `json:"-"`
    Role         string `json:"role" validate:"oneof=admin teacher"`
}

// HashPassword хэширует пароль
func (u *User) HashPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.PasswordHash = string(hashedPassword)
    return nil
}

// CheckPassword проверяет пароль
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
    return err == nil
}