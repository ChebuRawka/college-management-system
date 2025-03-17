package models

import "time"

type Schedule struct {
    ID            int       `json:"id"`
    TeacherName   string    `json:"teacher_name"`   //  (подтягивается через JOIN)
    ClassroomName string    `json:"classroom_name"` //  (подтягивается через JOIN)
    GroupName     string    `json:"group_name"`    
    StartTime     time.Time `json:"start_time"`    
    EndTime       time.Time `json:"end_time"`      
    DayOfWeek     string    `json:"day_of_week"`   //(например, "Monday")
}
