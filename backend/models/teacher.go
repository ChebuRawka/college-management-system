package models

import "github.com/go-playground/validator/v10"

type Teacher struct {
    ID      int      `json:"id"`
    Name    string   `json:"name" validate:"required"`
    Subject string   `json:"subject" validate:"required"`
    Courses []string `json:"courses"` 
}

var Validate *validator.Validate

func init() {
    Validate = validator.New()
}