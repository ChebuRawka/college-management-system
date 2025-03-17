package models

import "time"

type ScheduleResponse struct {
    ID            int       `json:"id"`
    TeacherName   string    `json:"teacher_name"`   // Имя преподавателя
    ClassroomName string    `json:"classroom_name"` // Название аудитории
    GroupName     string    `json:"group_name"`     // Название группы (может быть пустым)
    StartTime     time.Time `json:"start_time"`     // Время начала занятия
    EndTime       time.Time `json:"end_time"`       // Время окончания занятия
    DayOfWeek     string    `json:"day_of_week"`    // День недели (например, "Monday")
}