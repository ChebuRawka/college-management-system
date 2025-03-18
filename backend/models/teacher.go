package models

import (
    "github.com/go-playground/validator/v10"
)

type Teacher struct {
    ID           int       `json:"id"`
    Name         string    `json:"name" validate:"required"`       // Имя преподавателя
    Subject      string    `json:"subject" validate:"required"`    // Предмет, который преподает
    Courses      []string  `json:"courses" validate:"omitempty,dive,gt=0"` // Список имен курсов
    WorkingHours float64   `json:"working_hours" validate:"gte=0"` // Количество рабочих часов
}

var Validate *validator.Validate

func init() {
    Validate = validator.New()
}