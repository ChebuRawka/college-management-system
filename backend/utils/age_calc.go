package utils

import "time"

// CalculateAge вычисляет возраст на основе даты рождения
func CalculateAge(dateOfBirth time.Time) int {
    now := time.Now()
    age := now.Year() - dateOfBirth.Year()
    if now.Month() < dateOfBirth.Month() || (now.Month() == dateOfBirth.Month() && now.Day() < dateOfBirth.Day()) {
        age--
    }
    return age
}