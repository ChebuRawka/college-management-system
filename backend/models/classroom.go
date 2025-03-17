package models

type Classroom struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`        // Название аудитории (например, "Аудитория 101")
    Capacity    int     `json:"capacity"`    // Вместимость аудитории (количество мест)
    Description string  `json:"description"` // Описание аудитории (необязательное поле)
}